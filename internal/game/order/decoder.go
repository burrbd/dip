package order

import (
	"errors"
	"regexp"
	"strings"

	"github.com/burrbd/diplomacy/internal/game/order/board"
)

var (
	movePrefix = `([A|F])\s([A-Za-z]{3})\s([S|C|H])\s?`
	moveSuffix = `([A|F])\s([A-Za-z]{3})\-([A-Za-z]{3})`
	orderRE    = regexp.MustCompile(`(` + movePrefix + `)?(` + moveSuffix + `)?`)
)

func Decode(order string) (interface{}, error) {
	result := orderRE.FindAllStringSubmatch(order, -1)
	matches := result[0]
	if matches[0] == "" {
		return nil, errors.New("invalid order")
	}
	prefixMatch := matches[1] != ""
	moveMatch := matches[5] != ""
	if prefixMatch && !moveMatch {
		return decodeHold(matches)
	}
	unit := unitType(rune(matches[6][0]))
	from := strings.ToLower(matches[7])
	to := strings.ToLower(matches[8])
	move := Move{UnitType: unit, From: board.Territory{Abbr: from}, To: board.Territory{Abbr: to}}
	if prefixMatch {
		return decodePrefix(matches, move)
	}
	return move, nil
}

func decodePrefix(matches []string, move Move) (interface{}, error) {
	by := board.Territory{Abbr: strings.ToLower(matches[3])}
	unit := unitType(rune(matches[2][0]))
	if matches[4] == "C" {
		if unit == board.Army {
			return nil, errors.New("cannot convoy with army")
		}
		return MoveConvoy{Move: move, By: by}, nil
	}
	return MoveSupport{UnitType: unit, Move: move, By: by}, nil
}

func decodeHold(matches []string) (interface{}, error) {
	if matches[4] != "H" {
		return nil, errors.New("invalid order")
	}
	pos := board.Territory{Abbr: strings.ToLower(matches[3])}
	unit := unitType(rune(matches[2][0]))
	return Hold{UnitType: unit, Pos: pos}, nil
}

func unitType(letter rune) board.UnitType {
	if letter == 'F' {
		return board.Fleet
	}
	return board.Army
}
