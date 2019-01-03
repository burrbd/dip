package board

import (
	"sort"
	"strings"
)

type counterMoveConflicts map[string][2]*Unit

func (conflicts counterMoveConflicts) add(u *Unit) {
	key := movePairKey(u.Position().Territory, u.PrevPosition().Territory)
	pair := conflicts[key]
	if pair[0] == nil {
		conflicts[key] = [2]*Unit{u, pair[1]}
	} else {
		conflicts[key] = [2]*Unit{pair[0], u}
	}
}

func (conflicts counterMoveConflicts) del(u *Unit) {
	key := movePairKey(u.PrevPosition().Territory, u.Position().Territory)
	pair := conflicts[key]
	if pair[0] == u {
		conflicts[key] = [2]*Unit{nil, pair[1]}
	} else if pair[1] == u {
		conflicts[key] = [2]*Unit{pair[0], nil}
	}
}

func movePairKey(t1, t2 Territory) string {
	s := []string{t1.Abbr, t2.Abbr}
	sort.Strings(s)
	return strings.Join(s, "")
}
