package board

type UnitType string

const (
	Army  UnitType = "army"
	Fleet UnitType = "fleet"
)

type Unit struct {
	Country string
	Type    UnitType
}
