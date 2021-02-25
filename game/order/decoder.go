package order

import (
	"fmt"
	"strings"

	"github.com/burrbd/dip/game/order/board"
)

func Decode(order, country string) (interface{}, error) {
	tokens := strings.Split(strings.ToLower(order), " ")
	n := len(tokens)
	switch n {
	case 2:
		return decodeMove(tokens, country)
	case 3:
		return decodeHold(tokens, country)
	case 5:
		if tokens[2] == "c" {
			return decodeConvoy(tokens, country)
		} else if tokens[2] == "s" {
			return decodeSupport(tokens, country)
		}
		fallthrough
	default:
		return nil, fmt.Errorf("invalid order: %s", order)
	}
}

func decodeMove(tokens []string, country string) (interface{}, error) {
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
		Country:  country,
		UnitType: unit,
		From:     board.LookupTerritory(from),
		To:       board.LookupTerritory(to)}, nil
}

func decodeHold(tokens []string, country string) (interface{}, error) {
	unit, err := unitType(tokens[0])
	if err != nil {
		return nil, err
	}
	at := tokens[1]
	if tokens[2] != "h" {
		return nil, fmt.Errorf("invalid order: %s", strings.Join(tokens, " "))
	}
	return Hold{
		Country:  country,
		UnitType: unit,
		At:       board.LookupTerritory(at),
	}, nil
}

func decodeSupport(tokens []string, country string) (interface{}, error) {
	unit, err := unitType(tokens[0])
	if err != nil {
		return nil, err
	}
	fromTo := strings.Split(tokens[4], "-")
	if len(fromTo) == 2 {
		move, err := decodeMove(tokens[3:], country)
		if err != nil {
			return nil, err
		}
		return MoveSupport{
			Country:  country,
			UnitType: unit,
			By:       board.LookupTerritory(tokens[1]),
			Move:     move.(Move),
		}, nil
	}
	supportedUnit, err := unitType(tokens[3])
	if err != nil {
		return nil, err
	}
	return HoldSupport{
		Country:  country,
		UnitType: unit,
		By:       board.LookupTerritory(tokens[1]),
		Hold:     Hold{UnitType: supportedUnit, At: board.LookupTerritory(tokens[4])},
	}, nil
}

func decodeConvoy(tokens []string, country string) (interface{}, error) {
	if tokens[0] != "f" {
		return nil, fmt.Errorf("invalid order; only fleet can convoy: %s", strings.Join(tokens, " "))
	}
	fromTo := strings.Split(tokens[4], "-")
	if len(fromTo) != 2 {
		return nil, fmt.Errorf("invalid order: %s", strings.Join(tokens, " "))
	}
	move, err := decodeMove(tokens[3:], country)
	if err != nil {
		return nil, err
	}
	return MoveConvoy{
		Country: country,
		By:      board.LookupTerritory(tokens[1]),
		Move:    move.(Move),
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
		return board.UnitType(""), fmt.Errorf("invalid unit type: %s", unitToken)
	}
}
