package order_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/diplomacy/internal/game/order"
	"github.com/burrbd/diplomacy/internal/game/order/board"
)

var abc = board.Territory{Abbr: "abc"}
var def = board.Territory{Abbr: "def"}
var ghi = board.Territory{Abbr: "ghi"}

func TestDecode_ArmyMove(t *testing.T) {
	is := is.New(t)
	move := "A Abc-Def"
	decMove, err := order.Decode(move)
	is.NoErr(err)
	is.Equal(order.Move{UnitType: board.Army, From: abc, To: def}, decMove)
}

func TestDecode_FleetMove(t *testing.T) {
	is := is.New(t)
	move := "F Abc-Def"
	decMove, err := order.Decode(move)
	is.NoErr(err)
	is.Equal(order.Move{UnitType: board.Fleet, From: abc, To: def}, decMove)
}

func TestDecode_SupportMove(t *testing.T) {
	is := is.New(t)
	support := "A Ghi S A Abc-Def"
	decSupport, err := order.Decode(support)
	is.NoErr(err)
	expSupport := order.MoveSupport{
		UnitType: board.Army,
		By:       ghi,
		Move:     order.Move{UnitType: board.Army, From: abc, To: def}}
	is.Equal(expSupport, decSupport)
}

func TestDecode_SupportMoveWithFleet(t *testing.T) {
	is := is.New(t)
	support := "F Ghi S A Abc-Def"
	decSupport, err := order.Decode(support)
	is.NoErr(err)
	expSupport := order.MoveSupport{
		UnitType: board.Fleet,
		By:       ghi,
		Move:     order.Move{UnitType: board.Army, From: abc, To: def}}
	is.Equal(expSupport, decSupport)
}

func TestDecode_ConvoyMove(t *testing.T) {
	is := is.New(t)
	convoy := "F Ghi C A Abc-Def"
	decSupport, err := order.Decode(convoy)
	is.NoErr(err)
	expSupport := order.MoveConvoy{
		By:   ghi,
		Move: order.Move{UnitType: board.Army, From: abc, To: def}}
	is.Equal(expSupport, decSupport)
}

func TestDecode_ConvoyMoveWithArmyReturnsError(t *testing.T) {
	is := is.New(t)
	convoy := "A Ghi C A Abc-Def"
	_, err := order.Decode(convoy)
	is.Err(err)
}

func TestDecode_UnitHolds(t *testing.T) {
	is := is.New(t)
	hold := "A Abc H"
	decHold, err := order.Decode(hold)
	is.NoErr(err)
	expHold := order.Hold{
		UnitType: board.Army,
		Pos:      board.Territory{Abbr: "abc"}}
	is.Equal(expHold, decHold)
}

func TestDecode_WhenInvalidReturnError(t *testing.T) {
	is := is.New(t)
	move := "this is not a move"
	_, err := order.Decode(move)
	is.Err(err)
}

func TestDecode_WhenPrefixIsInvalidAndMoveValidReturnsError(t *testing.T) {
	is := is.New(t)
	convoy := "foobar C A Abc-Def"
	_, err := order.Decode(convoy)
	is.Err(err)
}

func TestDecode_WhenPrefixIsValidAndMoveInvalidReturnsError(t *testing.T) {
	is := is.New(t)
	convoy := "F Ghi C A AbcDef"
	_, err := order.Decode(convoy)
	is.Err(err)
}

func TestDecode_NoSpaceBetweenSuffixAndMove(t *testing.T) {
	is := is.New(t)
	convoy := "F Ghi CA Abc-Def"
	_, err := order.Decode(convoy)
	is.Err(err)
}
