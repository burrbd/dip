package order

import (
	"fmt"
	"strings"

	"github.com/burrbd/dip/game/order/board"
)

func Decode(order string) (interface{}, error) {
	tokens := strings.Split(strings.ToLower(order), " ")
	n := len(tokens)
	switch n {
	case 2:
		return decodeMove(tokens)
	case 3:
		return decodeHold(tokens)
	case 5:
		if tokens[2] == "c" {
			return decodeConvoy(tokens)
		} else if tokens[2] == "s" {
			return decodeSupport(tokens)
		}
		fallthrough
	default:
		return nil, fmt.Errorf("invalid order: %s", order)
	}
}

func decodeMove(tokens []string) (interface{}, error) {
	unit, err := unitType(tokens[0])
	if err != nil {
		return nil, err
	}
	fromTo := strings.Split(tokens[1], "-")
	if len(fromTo) != 2 {
		return nil, fmt.Errorf("invalid order: %s", strings.Join(tokens, " "))
	}
	from, to := fromTo[0], fromTo[1]
	return Move{
		UnitType: unit,
		From:     board.Territory{Abbr: from},
		To:       board.Territory{Abbr: to}}, nil
}

func decodeHold(tokens []string) (interface{}, error) {
	unit, err := unitType(tokens[0])
	if err != nil {
		return nil, err
	}
	at := tokens[1]
	if tokens[2] != "h" {
		return nil, fmt.Errorf("invalid order: %s", strings.Join(tokens, " "))
	}
	return Hold{
		UnitType: unit,
		At:       board.Territory{Abbr: at},
	}, nil
}

func decodeSupport(tokens []string) (interface{}, error) {
	unit, err := unitType(tokens[0])
	if err != nil {
		return nil, err
	}
	fromTo := strings.Split(tokens[4], "-")
	if len(fromTo) == 2 {
		move, err := decodeMove(tokens[3:])
		if err != nil {
			return nil, err
		}
		return MoveSupport{
			UnitType: unit,
			By:       board.Territory{Abbr: tokens[1]},
			Move:     move.(Move),
		}, nil
	}
	supportedUnit, err := unitType(tokens[3])
	if err != nil {
		return nil, err
	}
	return HoldSupport{
		UnitType: unit,
		By:       board.Territory{Abbr: tokens[1]},
		Hold:     Hold{UnitType: supportedUnit, At: board.Territory{Abbr: tokens[4]}},
	}, nil
}

func decodeConvoy(tokens []string) (interface{}, error) {
	if tokens[0] != "f" {
		return nil, fmt.Errorf("invalid order; only fleet can convoy: %s", strings.Join(tokens, " "))
	}
	fromTo := strings.Split(tokens[4], "-")
	if len(fromTo) != 2 {
		return nil, fmt.Errorf("invalid order: %s", strings.Join(tokens, " "))
	}
	move, err := decodeMove(tokens[3:])
	if err != nil {
		return nil, err
	}
	return MoveConvoy{
		By:   board.Territory{Abbr: tokens[1]},
		Move: move.(Move),
	}, nil
}

func unitType(unitToken string) (board.UnitType, error) {
	n := len(unitToken)
	switch {
	case n == 1 && unitToken[0] == 'f':
		return board.Fleet, nil
	case n == 1 && unitToken[0] == 'a':
		return board.Army, nil
	default:
		return board.UnitType('?'), fmt.Errorf("invalid unit type: %s", string(unitToken))
	}
}
