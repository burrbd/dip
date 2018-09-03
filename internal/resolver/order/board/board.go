package board

import (
	"github.com/burrbd/kit/graph"
	"github.com/satori/go.uuid"
)

type UnitType string

const (
	Army  UnitType = "army"
	Fleet UnitType = "fleet"
)

type Territory struct {
	Abbr string
	Name string
}

func (t Territory) ID() string {
	return t.Abbr
}

type Unit struct {
	id       string
	Country  string
	Type     UnitType
	Position *Position
}

func NewArmy(country string, t Territory, g *graph.Simple) (Unit, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Unit{}, err
	}
	return Unit{
		id:      id.String(),
		Country: country,
		Type:    Army,
		Position: &Position{
			Territory: t,
			Previous:  make([]Territory, 0),
			graph:     g,
		},
	}, nil
}

type Position struct {
	Territory Territory
	Previous  []Territory
	graph     *graph.Simple
}
