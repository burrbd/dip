package session

import (
	"encoding/json"
	"fmt"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/events"
)

// EngineLoader restores an engine.Engine from a JSON snapshot produced by Dump.
type EngineLoader func(snapshot []byte) (engine.Engine, error)

// Load rebuilds a Session for channelID from the channel's event log.
// loader is called to restore the engine from the most recent snapshot;
// pass engine.Load for production use.
func Load(ch events.Channel, channelID string, notifier Notifier, loader EngineLoader) (*Session, error) {
	envs, err := events.Scan(ch, channelID)
	if err != nil {
		return nil, fmt.Errorf("session: scan: %w", err)
	}

	s := &Session{
		ChannelID:    channelID,
		StagedOrders: make(map[string][]string),
		Players:      make(map[string]string),
		Submitted:    make(map[string]bool),
		ch:           ch,
		notifier:     notifier,
	}

	// snapshotIdx tracks the position of the last snapshot event so that only
	// OrderSubmitted events after it are included in StagedOrders.
	snapshotIdx := -1

	for i, env := range envs {
		switch env.Type {
		case events.TypeGameCreated:
			var gc events.GameCreated
			if err := json.Unmarshal(env.Payload, &gc); err != nil {
				continue
			}
			s.GMID = gc.GMUserID
			s.DeadlineHours = gc.DeadlineHours

		case events.TypePlayerJoined:
			var pj events.PlayerJoined
			if err := json.Unmarshal(env.Payload, &pj); err != nil {
				continue
			}
			s.Players[pj.UserID] = pj.Nation

		case events.TypeGameStarted:
			snapshotIdx = i
			s.StagedOrders = make(map[string][]string)
			s.Submitted = make(map[string]bool)

		case events.TypePhaseResolved:
			var pr events.PhaseResolved
			if err := json.Unmarshal(env.Payload, &pr); err != nil {
				continue
			}
			s.Phase = pr.Phase
			snapshotIdx = i
			s.StagedOrders = make(map[string][]string)
			s.Submitted = make(map[string]bool)

		case events.TypeOrderSubmitted:
			if snapshotIdx >= 0 && i > snapshotIdx {
				var os events.OrderSubmitted
				if err := json.Unmarshal(env.Payload, &os); err != nil {
					continue
				}
				s.StagedOrders[os.Nation] = append(s.StagedOrders[os.Nation], os.Orders...)
			}
		}
	}

	// Rebuild the engine from the last snapshot, replaying any orders after it.
	var eng engine.Engine
	if _, err := events.Rebuild(ch, channelID, func(snap []byte) (events.EngineState, error) {
		e, loadErr := loader(snap)
		if loadErr != nil {
			return nil, loadErr
		}
		eng = e
		return e, nil
	}); err != nil {
		return nil, fmt.Errorf("session: rebuild engine: %w", err)
	}

	s.Eng = eng
	return s, nil
}
