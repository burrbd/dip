package board

type Manager interface {
	Units() []*Unit
	Move(unit *Unit, territory Territory, strength int)
	Bounce(unit *Unit)
	SetDefeated(unit *Unit)
	Conflict() []*Unit
}

type PositionManager struct {
	territoryConflicts   territoryConflicts
	counterMoveConflicts counterMoveConflicts
}

func NewPositionManager() PositionManager {
	return PositionManager{
		territoryConflicts:   make(territoryConflicts, 0),
		counterMoveConflicts: make(counterMoveConflicts)}
}

func (m PositionManager) AddUnit(unit *Unit, territory Territory) {
	unit.PhaseHistory = make([]Position, 0)
	unit.PhaseHistory = append(unit.PhaseHistory, Position{
		Territory: territory,
		Cause:     Originated,
	})
	m.territoryConflicts.add(unit)
}

func (m PositionManager) Units() []*Unit {
	all := make([]*Unit, 0)
	for _, units := range m.territoryConflicts {
		all = append(all, units...)
	}
	return all
}

func (m PositionManager) Conflict() []*Unit {
	for _, moveConflict := range m.counterMoveConflicts {
		if moveConflict[0] != nil &&
			moveConflict[1] != nil &&
			!moveConflict[0].Defeated() &&
			!moveConflict[1].Defeated() {
			return []*Unit{moveConflict[0], moveConflict[1]}
		}
	}
	for _, units := range m.territoryConflicts {
		nonRetreatingUnits := unitFilter(units, func(u *Unit) bool { return !u.Defeated() })
		conflicts := len(nonRetreatingUnits)
		if conflicts > 1 {
			unitsCopy := make([]*Unit, conflicts)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	return nil
}

func (m PositionManager) Move(u *Unit, next Territory, strength int) {
	m.territoryConflicts.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: next, Strength: strength, Cause: Moved})
	m.territoryConflicts.add(u)
	if u.PrevPosition() != nil {
		m.counterMoveConflicts.add(u)
	}
}

func (m PositionManager) Bounce(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	m.counterMoveConflicts.del(u)
	m.territoryConflicts.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: u.PrevPosition().Territory, Strength: 0, Cause: Bounced})
	m.territoryConflicts.add(u)
}

func (m PositionManager) SetDefeated(u *Unit) {
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: u.Position().Territory, Cause: Defeated})
}

func unitFilter(units []*Unit, f func(*Unit) bool) []*Unit {
	filtered := make([]*Unit, 0)
	for _, u := range units {
		if f(u) {
			filtered = append(filtered, u)
		}
	}
	return filtered
}
