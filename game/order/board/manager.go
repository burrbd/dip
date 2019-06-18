package board

import (
	"sort"
	"strings"
)

// PositionEvent describes the causes of a unit's position
type PositionEvent int

const (
	// Added unit added to the board at the beginning of the phase
	Added PositionEvent = iota
	// Held unit has been held by player
	Held
	// Moved unit has moved territories
	Moved
	// Bounced unit has bounced from territory
	Bounced
	// Defeated unit has been defeated
	Defeated
)

// Position describes the unit's board position
type Position struct {
	Territory Territory
	Strength  int
	Cause     PositionEvent
}

// Manager is a board manager that records positions of units on a board
type Manager interface {
	Position(*Unit) *Position
	Positions() map[*Unit]Position
	Move(*Unit, Territory, int)
	Hold(*Unit, int)
	Bounce(*Unit)
	SetDefeated(*Unit)
	Conflict() []*Unit
	AtOrigin(*Unit) bool
}

// PositionManager implements Manager
type PositionManager struct {
	history map[*Unit][]Position
}

// NewPositionManager constructor for PositionManager
func NewPositionManager() PositionManager {
	return PositionManager{
		history: make(map[*Unit][]Position),
	}
}

// AddUnit places a unit on the board
func (m PositionManager) AddUnit(unit *Unit, territory Territory) {
	m.history[unit] = []Position{{
		Territory: territory,
		Cause:     Added,
	}}
}

// Positions returns all positions managed
func (m PositionManager) Positions() map[*Unit]Position {
	positions := make(map[*Unit]Position)
	for u := range m.history {
		positions[u] = *m.Position(u)
	}
	return positions
}

// Conflict returns the first conflict found on the board
func (m PositionManager) Conflict() []*Unit {
	conflicts := make(map[string][]*Unit)
	for u, position := range m.Positions() {
		if position.Cause == Defeated {
			continue
		}
		if position.Cause == Moved {
			conflicts = m.appendCounterAttackConflict(conflicts, u)
		}
		conflicts = m.appendTerritoryConflict(conflicts, u)
	}
	for _, units := range conflicts {
		if len(units) > 1 {
			return units
		}
	}
	return nil
}

// Move moves a unit from one territory to another
func (m PositionManager) Move(u *Unit, to Territory, strength int) {
	m.history[u] = append(m.history[u], Position{Territory: to, Strength: strength, Cause: Moved})
}

// Hold holds the unit in the current position
func (m PositionManager) Hold(u *Unit, strength int) {
	m.history[u] = append(m.history[u], Position{Territory: m.Position(u).Territory, Cause: Held, Strength: strength})
}

// Bounce bounces two or more units when there is no winner
func (m PositionManager) Bounce(u *Unit) {
	m.history[u] = append(m.history[u], Position{
		Territory: m.prevPosition(u).Territory, Strength: 0, Cause: Bounced})
}

// SetDefeated sets a unit's position as defeated
func (m PositionManager) SetDefeated(u *Unit) {
	m.history[u] = append(m.history[u], Position{Territory: m.Position(u).Territory, Cause: Defeated})
}

// Position returns the current board position of a unit
func (m PositionManager) Position(u *Unit) *Position {
	hist := m.positionHistory(u)
	return &hist[len(hist)-1]
}

// AtOrigin determines if the unit is located at the phase starting position
func (m PositionManager) AtOrigin(u *Unit) bool {
	pos, hist := m.Position(u), m.positionHistory(u)
	if pos == nil || hist == nil {
		return false
	}
	return pos.Territory.Abbr == hist[0].Territory.Abbr
}

// Defeated determines if the unit has been defeated
func (m PositionManager) Defeated(u *Unit) bool {
	return m.Position(u).Cause == Defeated
}

func (m PositionManager) positionHistory(u *Unit) []Position {
	hist, ok := m.history[u]
	if !ok || len(hist) == 0 {
		return nil
	}
	return hist
}

func (m PositionManager) prevPosition(u *Unit) *Position {
	hist := m.positionHistory(u)
	n := len(hist)
	if n < 2 {
		return nil
	}
	return &hist[n-2]
}

func (m PositionManager) appendCounterAttackConflict(conflicts map[string][]*Unit, u *Unit) map[string][]*Unit {
	s := []string{m.Position(u).Territory.Abbr, m.prevPosition(u).Territory.Abbr}
	sort.Strings(s)
	return appendConflict(conflicts, strings.Join(s, "."), u)
}

func (m PositionManager) appendTerritoryConflict(conflicts map[string][]*Unit, u *Unit) map[string][]*Unit {
	return appendConflict(conflicts, m.Position(u).Territory.Abbr, u)
}

func appendConflict(conflicts map[string][]*Unit, key string, u *Unit) map[string][]*Unit {
	if conflicts[key] == nil {
		conflicts[key] = make([]*Unit, 0)
	}
	conflicts[key] = append(conflicts[key], u)
	return conflicts
}
