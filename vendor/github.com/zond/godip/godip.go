// Package godip provides the core types for the Diplomacy adjudication engine.
package godip

// Nation is a string identifying a Great Power (e.g. "England", "France").
type Nation string

// Province is a string identifying a territory on the board (e.g. "Vie", "Bud").
type Province string

// UnitType identifies the kind of military unit.
type UnitType string

// PhaseType identifies the kind of game phase.
type PhaseType string

// Season identifies the in-game season.
type Season string

// OrderType identifies the kind of order.
type OrderType string

// Flag is a bitmask flag used by orders.
type Flag int

const (
	Army  UnitType = "Army"
	Fleet UnitType = "Fleet"
)

const (
	Movement   PhaseType = "Movement"
	Retreat    PhaseType = "Retreat"
	Adjustment PhaseType = "Adjustment"
)

const (
	Spring Season = "Spring"
	Fall   Season = "Fall"
)

// Unit is a military unit on the board.
type Unit struct {
	Type   UnitType
	Nation Nation
}

// Order is an instruction for a unit.
type Order interface {
	Type() OrderType
	Flags() map[Flag]bool
}

// Phase represents the current game phase.
type Phase interface {
	Type() PhaseType
	Year() int
	Season() Season
	DefaultOrder(Province) Order
}

// Adjudicator is the stateful godip game state that can stage orders and be adjudicated.
type Adjudicator interface {
	Phase() Phase
	Orders() map[Province]Order
	Units() map[Province]Unit
	Dislodgeds() map[Province]Unit
	SetOrder(Province, Order)
	Resolve(Province) error
	Next() (Adjudicator, error)
	SoloWinner() Nation
	Dump() ([]byte, error)
}
