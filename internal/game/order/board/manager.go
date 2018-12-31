package board

type Manager interface {
	Units() []*Unit
	Move(unit *Unit, territory Territory, strength int)
	Bounce(unit *Unit)
	SetDefeated(unit *Unit)
	Conflict() []*Unit
}

type PositionManager struct {
	positions positionRecorder
}

func NewPositionManager(units []*Unit) PositionManager {
	return PositionManager{newPositionRecorder(units)}
}

func (m PositionManager) Units() []*Unit {
	all := make([]*Unit, 0)
	for _, units := range m.positions.unitsByTerr {
		all = append(all, units...)
	}
	return all
}

func (m PositionManager) Conflict() []*Unit {
	for _, units := range m.positions.unitsByMovePair {
		if len(units) == 2 && !units[0].Defeated() && !units[1].Defeated() {
			unitsCopy := make([]*Unit, 2)
			copy(unitsCopy, units)
			return unitsCopy
		}
	}
	for _, units := range m.positions.unitsByTerr {
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
	m.positions.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: next, Strength: strength, Cause: Moved})
	m.positions.add(u)
	m.positions.addMovePair(u)
}

func (m PositionManager) Bounce(u *Unit) {
	prev := u.PrevPosition()
	if prev == nil {
		return
	}
	next := prev.Territory
	m.positions.delMovePair(u)
	m.positions.del(u)
	u.PhaseHistory = append(u.PhaseHistory, Position{
		Territory: next, Strength: 0, Cause: Bounced})
	m.positions.add(u)
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
