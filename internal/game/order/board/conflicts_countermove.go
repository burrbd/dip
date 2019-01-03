package board

import (
	"sort"
	"strings"
)

type counterMoveConflicts map[string][]*Unit

func (conflicts counterMoveConflicts) add(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	key := movePairKey(u.Position().Territory, u.PrevPosition().Territory)
	if _, ok := conflicts[key]; !ok {
		conflicts[key] = make([]*Unit, 0, 2)
	}
	conflicts[key] = append(conflicts[key], u)
}

func (conflicts counterMoveConflicts) del(u *Unit) {
	if u.PrevPosition() == nil {
		return
	}
	key := movePairKey(u.PrevPosition().Territory, u.Position().Territory)
	pair := conflicts[key]
	for i, cu := range pair {
		if u == cu {
			conflicts[key] = removeIndex(i, pair)
		}
	}
}

func movePairKey(t1, t2 Territory) string {
	s := []string{t1.Abbr, t2.Abbr}
	sort.Strings(s)
	return strings.Join(s, "")
}
