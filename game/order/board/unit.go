package board

type Graph interface {
	IsNeighbour(t1, t2 string) (bool, error)
}

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
	Country string
	Type    UnitType
}
