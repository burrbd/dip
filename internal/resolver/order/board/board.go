package board

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

type Position struct {
	Unit Unit
	//Retreats  []Unit
	Territory Territory
}
