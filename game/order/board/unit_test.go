package board_test

//import (
//	"sort"
//	"testing"
//
//	"github.com/burrbd/dip/game/order/board"
//	"github.com/cheekybits/is"
//)
//
//func TestUnitsByStrength(t *testing.T) {
//	// A Bud-Gal
//	// A Vie S Bud-Gal
//	// A Par-Lon
//	is := is.New(t)
//	u1 := &board.Unit{PhaseHistory: []board.Position{{Strength: 3}}}
//	u2 := &board.Unit{PhaseHistory: []board.Position{{Strength: 2}}}
//	u3 := &board.Unit{PhaseHistory: []board.Position{{Strength: 1}}}
//
//	units := []*board.Unit{u2, u1, u3}
//	sort.Sort(board.UnitPositionsByStrength(units))
//
//	is.Equal(u1, units[0])
//	is.Equal(u2, units[1])
//	is.Equal(u3, units[2])
//}
//
//func TestUnitsByStrength_NoPosition(t *testing.T) {
//	is := is.New(t)
//	u1 := &board.Unit{PhaseHistory: []board.Position{{Strength: 1}}}
//	u2 := &board.Unit{}
//
//	units := []*board.Unit{u2, u1}
//	sort.Sort(board.UnitPositionsByStrength(units))
//
//	is.Equal(u1, units[0])
//	is.Equal(u2, units[1])
//}
//
//func TestUnit_AtOrigin(t *testing.T) {
//	is := is.New(t)
//
//	u := board.Unit{
//		PhaseHistory: []board.Position{
//			{Territory: board.Territory{Abbr: "t1"}, Cause: board.Added}}}
//	is.True(u.AtOrigin())
//}
//
//func TestUnit_AtOrigin_Held(t *testing.T) {
//	is := is.New(t)
//
//	u := board.Unit{
//		PhaseHistory: []board.Position{
//			{Territory: board.Territory{Abbr: "t1"}, Cause: board.Held}}}
//	is.True(u.AtOrigin())
//}
//
//func TestUnit_AtOrigin_WhenNotAtOrigin(t *testing.T) {
//	is := is.New(t)
//
//	u := board.Unit{
//		PhaseHistory: []board.Position{
//			{Territory: board.Territory{Abbr: "t1"}, Cause: board.Added},
//			{Territory: board.Territory{Abbr: "t2"}}}}
//	is.False(u.AtOrigin())
//}
//
//func TestUnit_Defeated(t *testing.T) {
//	is := is.New(t)
//
//	u := board.Unit{
//		PhaseHistory: []board.Position{
//			{Territory: board.Territory{Abbr: "t1"}, Cause: board.Defeated}}}
//	is.True(u.Defeated())
//}
