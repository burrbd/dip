package board_test

import (
	"testing"

	"github.com/burrbd/dip/game/order/board"
	"github.com/cheekybits/is"
)

func TestPositionManager_Position(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	terr := board.Territory{Abbr: "t1"}
	u := &board.Unit{}

	m.AddUnit(u, terr)

	is.Equal(terr, m.Position(u).Territory)
}

func TestPositionManager_Positions(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u1 := &board.Unit{}
	u2 := &board.Unit{}
	m.AddUnit(u1, t1)
	m.AddUnit(u2, t2)

	positions := m.Positions()

	is.Equal(2, len(positions))
	is.Equal("t1", positions[u1].Territory.Abbr)
	is.Equal("t2", positions[u2].Territory.Abbr)
}

func TestPositionManager_Conflict(t *testing.T) {
	is := is.New(t)

	specs := []struct {
		desc                string
		positions           []unitPositionInstruction
		conflictedUnits     int
		conflictedTerritory string
	}{
		{
			desc:            "1 unit",
			positions:       units(add(with(board.UnitPlaced, "a_terr"))),
			conflictedUnits: 0,
		},
		{
			desc: "2 units on different territory",
			positions: units(
				add(with(board.UnitPlaced, "a_terr")),
				add(with(board.UnitPlaced, "b_terr"))),
			conflictedUnits: 0,
		},
		{
			desc: "2 units on same territory",
			positions: units(
				add(with(board.UnitPlaced, "a_terr")),
				add(with(board.UnitPlaced, "a_terr"))),
			conflictedUnits:     2,
			conflictedTerritory: "a_terr",
		},
		{
			desc: "3 units on same territory",
			positions: units(
				add(with(board.UnitPlaced, "a_terr")),
				add(with(board.UnitPlaced, "a_terr")),
				add(with(board.UnitPlaced, "a_terr"))),
			conflictedUnits:     3,
			conflictedTerritory: "a_terr",
		},
		{
			desc: "2 units on same territory, 1 unit different",
			positions: units(
				add(with(board.UnitPlaced, "b_terr")),
				add(with(board.UnitPlaced, "a_terr")),
				add(with(board.UnitPlaced, "b_terr"))),
			conflictedUnits:     2,
			conflictedTerritory: "b_terr",
		},
		{
			desc: "2 units move to same territory",
			positions: units(
				add(with(board.UnitPlaced, "b_terr"), with(board.Moved, "c_terr", 0)),
				add(with(board.UnitPlaced, "a_terr"), with(board.Moved, "c_terr", 0))),
			conflictedUnits:     2,
			conflictedTerritory: "c_terr",
		},
		{
			desc: "2 units move into each other's territory",
			positions: units(
				add(with(board.UnitPlaced, "a_terr"), with(board.Moved, "b_terr", 0)),
				add(with(board.UnitPlaced, "b_terr"), with(board.Moved, "a_terr", 0))),
			conflictedUnits: 2,
		},
		{
			desc: "2 units in same territory, 1 unit defeated",
			positions: units(
				add(with(board.UnitPlaced, "a_terr"), with(board.Moved, "b_terr", 0)),
				add(with(board.UnitPlaced, "b_terr"), with(board.Defeated, "b_terr", 0))),
			conflictedUnits: 0,
		},
	}

	for _, spec := range specs {
		m := board.NewPositionManager()
		t.Run(spec.desc, func(t *testing.T) {
			for _, unitPositions := range spec.positions {
				u := &board.Unit{}
				for _, position := range unitPositions {
					switch position[0] {
					case board.UnitPlaced:
						m.AddUnit(u, board.Territory{Abbr: position[1].(string)})
					case board.Moved:
						m.Move(u, board.Territory{Abbr: position[1].(string)}, 0)
					case board.Defeated:
						m.SetDefeated(u)
					}
				}
			}
			conflicts := m.Conflict()
			is.Equal(spec.conflictedUnits, len(conflicts))
			if spec.conflictedTerritory != "" {
				for _, unit := range conflicts {
					is.Equal(spec.conflictedTerritory, m.Position(unit).Territory.Abbr)
				}
			}
		})
	}
}

func TestPositionManager_AddUnit(t *testing.T) {
	is := is.New(t)
	terr := board.Territory{Abbr: "terr"}
	u := &board.Unit{}
	m := board.NewPositionManager()
	m.AddUnit(u, terr)
	is.Equal(board.UnitPlaced, m.Position(u).Cause)
	is.Equal(terr, m.Position(u).Territory)
	is.Equal(0, m.Position(u).Strength)
}

func TestPositionManager_Move(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u := &board.Unit{}
	m.AddUnit(u, t1)
	m.Move(u, t2, 1)

	is.Equal(t2, m.Position(u).Territory)
	is.Equal(board.Moved, m.Position(u).Cause)
	is.Equal(1, m.Position(u).Strength)
}

func TestPositionManager_Hold(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	ter := board.Territory{Abbr: "t"}
	u := &board.Unit{}
	m.AddUnit(u, ter)

	m.Hold(u, 2)

	is.Equal(ter, m.Position(u).Territory)
	is.Equal(board.Held, m.Position(u).Cause)
	is.Equal(2, m.Position(u).Strength)
}

func TestManager_Bounce(t *testing.T) {
	is := is.New(t)

	m := board.NewPositionManager()

	t1 := board.Territory{Abbr: "t1"}
	t2 := board.Territory{Abbr: "t2"}
	u := &board.Unit{}
	m.AddUnit(u, t1)

	is.Equal(t1, m.Position(u).Territory)
	is.Equal(board.UnitPlaced, m.Position(u).Cause)

	m.Move(u, t2, 1)
	m.Bounce(u)

	is.Equal(t1, m.Position(u).Territory)
	is.Equal(board.Bounced, m.Position(u).Cause)
	is.Equal(0, m.Position(u).Strength)
}

func TestManager_SetDefeated(t *testing.T) {
	is := is.New(t)
	m := board.NewPositionManager()
	terr := board.Territory{Abbr: "terr"}
	u := &board.Unit{}
	m.AddUnit(u, terr)

	m.SetDefeated(u)

	is.Equal(terr, m.Position(u).Territory)
	is.Equal(board.Defeated, m.Position(u).Cause)
	is.Equal(0, m.Position(u).Strength)
}

func with(position ...interface{}) positionInstruction {
	return position
}

type positionInstruction []interface{}

func add(positions ...positionInstruction) []positionInstruction {
	return positions
}

type unitPositionInstruction []positionInstruction

func units(positions ...unitPositionInstruction) []unitPositionInstruction {
	return positions
}
