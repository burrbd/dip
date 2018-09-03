package board

import (
	"fmt"
)

type Map struct {
	positions map[string][]Unit
}

func (m Map) AddPosition(u Unit) error {
	pos := u.Position
	if pos == nil {
		return fmt.Errorf("unit does not have a position defined")
	}
	if pos.graph == nil {
		return fmt.Errorf("position does not belong to a graph")
	}
	terr := pos.Territory.ID()
	if _, ok := m.positions[terr]; !ok {
		m.positions[terr] = make([]Unit, 0)
	}
	m.positions[terr] = append(m.positions[terr], u)
	return nil
}

func (m Map) Contests(terr string) []Unit {
	if _, ok := m.positions[terr]; !ok {
		return nil
	}
	if len(m.positions[terr]) == 0 {
		return nil
	}
	return m.positions[terr]
}

// Move a unit to a territory.
func (m Map) Move(u Unit, t Territory) error {
	pos := u.Position
	if pos == nil {
		return fmt.Errorf("unit does not have position")
	}
	if pos.graph == nil {
		return fmt.Errorf("position does not have a graph")
	}
	ok, err := pos.graph.IsNeighbor(pos.Territory, t)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s not a neighbor of %s", t.ID(), pos.Territory.ID())
	}
	i, err := m.findPosition(u)
	if err != nil {
		return err
	}
	m.deletePosition(i, pos.Territory)
	if _, ok := m.positions[t.ID()]; !ok {
		m.positions[t.ID()] = make([]Unit, 0)
	}
	pos.Previous = append([]Territory{pos.Territory}, pos.Previous...)
	pos.Territory = t
	m.positions[t.ID()] = append(m.positions[t.ID()], u)
	return nil
}

func (m Map) Bounce(u Unit, t Territory) error {
	return nil
}

func (m Map) findPosition(u Unit) (int, error) {
	if u.Position == nil {
		return -1, fmt.Errorf("unit does not have position")
	}
	units := m.positions[u.Position.Territory.ID()]
	for i, uu := range units {
		if uu.id == u.id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("map does not have unit")
}

func (m *Map) deletePosition(index int, t Territory) {
	if _, ok := m.positions[t.ID()]; !ok {
		return
	}
	if len(m.positions[t.ID()]) <= index {
		return
	}
	units := m.positions[t.ID()]
	units[index] = units[0]
	m.positions[t.ID()] = units[1:]
}
