package order

import (
	"fmt"

	"github.com/burrbd/dip/game/order/board"
)

type simpleGraph interface {
	HasEdgeBetween(xid, yid int64) bool
}

type Validator struct {
	armyGraph simpleGraph
}

func NewValidator(g simpleGraph) Validator {
	return Validator{armyGraph: g}
}

func (v Validator) ValidateMove(unit board.Unit, move Move) error {
	if move.Country != unit.Country {
		return fmt.Errorf("invalid country: %s", move.Country)
	}
	if !v.armyGraph.HasEdgeBetween(move.From.ID(), move.To.ID()) {
		return fmt.Errorf("cannot move from %s to %s", move.From.Abbr, move.To.Abbr)
	}
	return nil
}

func (v Validator) ValidateMoveSupport(unit board.Unit, sup MoveSupport) error {
	if sup.Country != unit.Country {
		return fmt.Errorf("invalid country: %s", sup.Country)
	}

	if !v.armyGraph.HasEdgeBetween(sup.By.ID(), sup.Move.To.ID()) {
		return fmt.Errorf("cannot support to %s by %s", sup.Move.To.Abbr, sup.By.Abbr)
	}

	return nil
}
