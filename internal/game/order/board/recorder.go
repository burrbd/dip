package board

import (
	"sort"
	"strings"
)

type positionRecorder struct {
	unitsByTerr     map[string][]*Unit
	unitsByMovePair map[string][]*Unit
}

func newPositionRecorder(units []*Unit) positionRecorder {
	terrUnits := make(map[string][]*Unit)
	for _, u := range units {
		if _, ok := terrUnits[u.Position().Territory.Abbr]; !ok {
			terrUnits[u.Position().Territory.Abbr] = make([]*Unit, 0)
		}
		terrUnits[u.Position().Territory.Abbr] = append(terrUnits[u.Position().Territory.Abbr], u)
	}
	return positionRecorder{
		unitsByTerr:     terrUnits,
		unitsByMovePair: make(map[string][]*Unit),
	}
}

func (r positionRecorder) add(u *Unit) {
	terr := u.Position().Territory.Abbr
	if _, ok := r.unitsByTerr[terr]; !ok {
		r.unitsByTerr[terr] = make([]*Unit, 0)
	}
	r.unitsByTerr[terr] = append(r.unitsByTerr[terr], u)
}

func (r positionRecorder) del(u *Unit) {
	terr := u.Position().Territory.Abbr
	units, ok := r.unitsByTerr[terr]
	if !ok {
		return
	}
	for i, unit := range units {
		if u == unit {
			units = removeIndex(i, units)
		}
	}
	r.unitsByTerr[terr] = units
}

func (r positionRecorder) addMovePair(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	key := movePairKey(u.Position().Territory, u.PrevPosition().Territory)
	if _, ok := r.unitsByMovePair[key]; !ok {
		r.unitsByMovePair[key] = make([]*Unit, 0, 2)
	}
	r.unitsByMovePair[key] = append(r.unitsByMovePair[key], u)
}

func (r positionRecorder) delMovePair(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	key := movePairKey(u.PrevPosition().Territory, u.Position().Territory)
	units := r.unitsByMovePair[key]
	for i, cu := range units {
		if u == cu {
			r.unitsByMovePair[key] = removeIndex(i, units)
		}
	}
}

func movePairKey(t1, t2 Territory) string {
	s := []string{t1.Abbr, t2.Abbr}
	sort.Strings(s)
	return strings.Join(s, "")
}

func removeIndex(i int, units []*Unit) []*Unit {
	copy(units[i:], units[i+1:])
	units[len(units)-1] = nil
	return units[:len(units)-1]
}
