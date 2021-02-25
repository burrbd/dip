package order_test

import (
	"testing"

	"github.com/cheekybits/is"

	"github.com/burrbd/dip/game/order"
	"github.com/burrbd/dip/game/order/board"
)

var (
	vie     = board.LookupTerritory("vie")
	bud     = board.LookupTerritory("bud")
	gal     = board.LookupTerritory("gal")
	country = "a_country"
)

func TestDecode_ArmyMove(t *testing.T) {
	is := is.New(t)
	move := "A Vie-Bud"
	decMove, err := order.Decode(move, country)
	is.NoErr(err)
	is.Equal(order.Move{UnitType: board.Army, From: vie, To: bud, Country: country}, decMove)
}

func TestDecode_FleetMove(t *testing.T) {
	// this is an invalid move
	is := is.New(t)
	move := "F Vie-Bud"
	decMove, err := order.Decode(move, country)
	is.NoErr(err)
	is.Equal(order.Move{UnitType: board.Fleet, From: vie, To: bud, Country: country}, decMove)
}

func TestDecode_Move_Errors(t *testing.T) {
	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		move := "Z Vie-Bud"
		_, err := order.Decode(move, country)
		is.Err(err)
	})

	t.Run("hyphen split error", func(t *testing.T) {
		is := is.New(t)
		move := "A Vie-Bud-Gal"
		_, err := order.Decode(move, country)
		is.Err(err)
	})
}

func TestDecode_SupportMove(t *testing.T) {
	is := is.New(t)
	support := "A Gal S A Vie-Bud"
	decSupport, err := order.Decode(support, country)
	is.NoErr(err)
	expSupport := order.MoveSupport{
		UnitType: board.Army,
		By:       gal,
		Move:     order.Move{UnitType: board.Army, From: vie, To: bud, Country: country},
		Country:  country,
	}
	is.Equal(expSupport, decSupport)
}

func TestDecode_SupportMoveWithFleet(t *testing.T) {
	is := is.New(t)
	support := "F Gal S A Vie-Bud"
	decSupport, err := order.Decode(support, country)
	is.NoErr(err)
	expSupport := order.MoveSupport{
		UnitType: board.Fleet,
		By:       gal,
		Move:     order.Move{UnitType: board.Army, From: vie, To: bud, Country: country},
		Country:  country,
	}
	is.Equal(expSupport, decSupport)
}

func TestDecode_Support_Errors(t *testing.T) {
	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		support := "& Gal S A Vie"
		_, err := order.Decode(support, country)
		is.Err(err)
	})

	t.Run("invalid move unit type", func(t *testing.T) {
		is := is.New(t)
		support := "F Gal S & Vie-Bud"
		_, err := order.Decode(support, country)
		is.Err(err)
	})

	t.Run("invalid hold unit type", func(t *testing.T) {
		is := is.New(t)
		support := "F Gal S & Vie"
		_, err := order.Decode(support, country)
		is.Err(err)
	})
}

func TestDecode_ConvoyMove(t *testing.T) {
	is := is.New(t)
	convoy := "F Gal C A Vie-Bud"
	decSupport, err := order.Decode(convoy, country)
	is.NoErr(err)
	expSupport := order.MoveConvoy{
		By:      gal,
		Move:    order.Move{UnitType: board.Army, From: vie, To: bud, Country: country},
		Country: country,
	}
	is.Equal(expSupport, decSupport)
}

func TestDecode_ConvoyMove_Errors(t *testing.T) {
	t.Run("convoy with invalid unit type", func(t *testing.T) {
		is := is.New(t)
		convoy := "& Gal C A Vie-Bud"
		_, err := order.Decode(convoy, country)
		is.Err(err)
	})

	t.Run("convoy invalid move", func(t *testing.T) {
		is := is.New(t)
		convoy := "F Gal C & Vie-Bud"
		_, err := order.Decode(convoy, country)
		is.Err(err)
	})

}

func TestDecode_Hold(t *testing.T) {
	is := is.New(t)
	hold := "A Vie H"
	decHold, err := order.Decode(hold, country)
	is.NoErr(err)
	expHold := order.Hold{
		UnitType: board.Army,
		At:       vie,
		Country:  country,
	}
	is.Equal(expHold, decHold)
}

func TestDecode_Hold_Errors(t *testing.T) {
	t.Run("with wrong hold token", func(t *testing.T) {
		is := is.New(t)
		hold := "A Vie s"
		_, err := order.Decode(hold, country)
		is.Err(err)
	})

	t.Run("invalid unit type", func(t *testing.T) {
		is := is.New(t)
		move := "Z Vie H"
		_, err := order.Decode(move, country)
		is.Err(err)
	})
}

func TestDecode_SupportUnitHold(t *testing.T) {
	is := is.New(t)
	support := "A Gal S A Vie"
	decSupport, err := order.Decode(support, country)
	is.NoErr(err)
	expSupport := order.HoldSupport{
		UnitType: board.Army,
		By:       gal,
		Hold:     order.Hold{UnitType: board.Army, At: vie},
		Country:  country,
	}
	is.Equal(expSupport, decSupport)

}

func TestDecode_WhenInvalidReturnError(t *testing.T) {
	is := is.New(t)
	move := "this is not a move"
	_, err := order.Decode(move, country)
	is.Err(err)
}

func TestDecode_WhenPrefixIsInvalidAndMoveValidReturnsError(t *testing.T) {
	is := is.New(t)
	convoy := "foobar C A Vie-Bud"
	_, err := order.Decode(convoy, country)
	is.Err(err)
}

func TestDecode_WhenPrefixIsValidAndMoveInvalidReturnsError(t *testing.T) {
	is := is.New(t)
	convoy := "F Gal C A VieBud"
	_, err := order.Decode(convoy, country)
	is.Err(err)
}

func TestDecode_NoSpaceBetweenSuffixAndMove(t *testing.T) {
	is := is.New(t)
	convoy := "F Gal CA Vie-Bud"
	_, err := order.Decode(convoy, country)
	is.Err(err)
}
