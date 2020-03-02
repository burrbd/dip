package order_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
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

func TestDecode_Move_Errors(t *testing.T) {
	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		move := "Z Abc-Def"
		_, err := order.Decode(move)
		is.Err(err)
	})

	t.Run("hyphen split error", func(t *testing.T) {
		is := is.New(t)
		move := "A Abc-Def-Ghi"
		_, err := order.Decode(move)
		is.Err(err)
	})
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

func TestDecode_Support_Errors(t *testing.T) {
	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		support := "& Ghi S A Abc"
		_, err := order.Decode(support)
		is.Err(err)
	})

	t.Run("invalid move unit type", func(t *testing.T) {
		is := is.New(t)
		support := "F Ghi S & Abc-Def"
		_, err := order.Decode(support)
		is.Err(err)
	})

	t.Run("invalid hold unit type", func(t *testing.T) {
		is := is.New(t)
		support := "F Ghi S & Abc"
		_, err := order.Decode(support)
		is.Err(err)
	})
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

func TestDecode_ConvoyMove_Errors(t *testing.T) {
	t.Run("convoy with invalid unit type", func(t *testing.T) {
		is := is.New(t)
		convoy := "& Ghi C A Abc-Def"
		_, err := order.Decode(convoy)
		is.Err(err)
	})

	t.Run("convoy invalid move", func(t *testing.T) {
		is := is.New(t)
		convoy := "F Ghi C & Abc-Def"
		_, err := order.Decode(convoy)
		is.Err(err)
	})

}

func TestDecode_Hold(t *testing.T) {
	is := is.New(t)
	hold := "A Abc H"
	decHold, err := order.Decode(hold)
	is.NoErr(err)
	expHold := order.Hold{
		UnitType: board.Army,
		At:       board.Territory{Abbr: "abc"}}
	is.Equal(expHold, decHold)
}

func TestDecode_Hold_Errors(t *testing.T) {
	t.Run("with wrong hold token", func(t *testing.T) {
		is := is.New(t)
		hold := "A Abc s"
		_, err := order.Decode(hold)
		is.Err(err)
	})

	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		move := "Z Abc H"
		_, err := order.Decode(move)
		is.Err(err)
	})
}

func TestDecode_SupportUnitHold(t *testing.T) {
	is := is.New(t)
	support := "A Ghi S A Abc"
	decSupport, err := order.Decode(support)
	is.NoErr(err)
	expSupport := order.HoldSupport{
		UnitType: board.Army,
		By:       ghi,
		Hold:     order.Hold{UnitType: board.Army, At: abc}}
	is.Equal(expSupport, decSupport)

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
