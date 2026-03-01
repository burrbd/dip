// Package variants defines the Diplomacy variant type used to start different
// game configurations (Classical, etc.).
package variants

import "github.com/zond/godip"

// Variant describes a Diplomacy map variant and how to start a new game on it.
type Variant struct {
	Name  string
	Start func() (godip.Adjudicator, error)
}
