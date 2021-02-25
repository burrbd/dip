package board

import (
	"strings"

	"gonum.org/v1/gonum/graph/simple"
)

type Territory struct {
	id    int64
	Abbr  string
	Name  string
	edges []string
}

func (t Territory) ID() int64 {
	return t.id
}

func (t Territory) Is(cmp Territory) bool {
	return t.id == cmp.id
}

func (t Territory) IsNot(cmp Territory) bool {
	return t.id != cmp.id
}

func CreateArmyGraph() *simple.UndirectedGraph {
	tm := make(map[string]Territory)
	for _, terr := range armyTerritories {
		tm[terr.Abbr] = terr
	}

	g := simple.NewUndirectedGraph()
	for _, terr := range tm {
		for _, neighbor := range terr.edges {
			g.SetEdge(g.NewEdge(terr, tm[neighbor]))

		}
	}
	return g
}

func LookupTerritory(abbr string) Territory {
	for _, terr := range armyTerritories {
		if strings.ToLower(abbr) == terr.Abbr {
			return terr
		}
	}
	return Territory{id: -1} // TODO: why this way?
}

var armyTerritories = []Territory{
	Territory{id: 0, Name: "Albania", Abbr: "alb", edges: []string{"gre", "ser", "tri"}},
	Territory{id: 1, Name: "Ankara", Abbr: "ank", edges: []string{"arm", "con", "smy"}},
	Territory{id: 2, Name: "Apulia", Abbr: "apu", edges: []string{"nap", "rom", "ven"}},
	Territory{id: 3, Name: "Armenia", Abbr: "arm", edges: []string{"sev", "smy", "syr"}},
	Territory{id: 4, Name: "Belgium", Abbr: "bel", edges: []string{"bur", "hol", "pic", "rur"}},
	Territory{id: 5, Name: "Berlin", Abbr: "ber", edges: []string{"kie", "mun", "pru", "sil"}},
	Territory{id: 6, Name: "Bohemia", Abbr: "boh", edges: []string{"gal", "mun", "sil", "tyr"}},
	Territory{id: 7, Name: "Brest", Abbr: "bre", edges: []string{"gas", "par", "pic"}},
	Territory{id: 8, Name: "Budapest", Abbr: "bud", edges: []string{"gal", "rum", "ser", "tri", "vie"}},
	Territory{id: 9, Name: "Bulgaria", Abbr: "bul", edges: []string{"con", "gre", "rum", "ser"}},
	Territory{id: 10, Name: "Burgundy", Abbr: "bur", edges: []string{"bel", "gas", "mar", "mun", "par", "pic", "ruh"}},
	Territory{id: 11, Name: "Clyde", Abbr: "cly", edges: []string{"edi", "lvp"}},
	Territory{id: 12, Name: "Constantinople", Abbr: "con", edges: []string{"smy", "ank", "bul"}},
	Territory{id: 13, Name: "Denmark", Abbr: "den", edges: []string{"kie", "swe"}},
	Territory{id: 14, Name: "Edinburgh", Abbr: "edi", edges: []string{"lvp", "yor", "cly"}},
	Territory{id: 15, Name: "Finland", Abbr: "fin", edges: []string{"nwy", "stp", "swe"}},
	Territory{id: 16, Name: "Galacia", Abbr: "gal", edges: []string{"rum", "sil", "ukr", "vie", "war", "boh", "bud"}},
	Territory{id: 17, Name: "Gascony", Abbr: "gas", edges: []string{"mar", "par", "spa", "bre", "bur"}},
	Territory{id: 18, Name: "Greece", Abbr: "gre", edges: []string{"ser", "alb", "bul"}},
	Territory{id: 19, Name: "Holland", Abbr: "hol", edges: []string{"kie", "ruh", "bel"}},
	Territory{id: 20, Name: "Kiel", Abbr: "kie", edges: []string{"mun", "ruh", "ber", "den", "hol"}},
	Territory{id: 21, Name: "Liverpool", Abbr: "lvp", edges: []string{"wal", "yor", "cly", "edi"}},
	Territory{id: 22, Name: "Livonia", Abbr: "lvn", edges: []string{"mos", "pru", "stp", "war"}},
	Territory{id: 23, Name: "London", Abbr: "lon", edges: []string{"wal", "yor"}},
	Territory{id: 24, Name: "Marseilles", Abbr: "mar", edges: []string{"pie", "spa", "bur", "gas"}},
	Territory{id: 25, Name: "Moscow", Abbr: "mos", edges: []string{"sev", "stp", "urk", "war", "lvn"}},
	Territory{id: 26, Name: "Munich", Abbr: "mun", edges: []string{"ruh", "sil", "tyr", "ber", "boh", "bur", "kie"}},
	Territory{id: 27, Name: "Naples", Abbr: "nap", edges: []string{"rom", "apu"}},
	Territory{id: 28, Name: "North Africa", Abbr: "naf", edges: []string{"tun"}},
	Territory{id: 29, Name: "Norway", Abbr: "nwy", edges: []string{"stp", "swe", "fin"}},
	Territory{id: 30, Name: "Paris", Abbr: "par", edges: []string{"pic", "bre", "bur", "gas"}},
	Territory{id: 31, Name: "Picardy", Abbr: "pic", edges: []string{"bel", "bre", "bur", "par"}},
	Territory{id: 32, Name: "Piedmont", Abbr: "pie", edges: []string{"tus", "tyr", "ven", "mar"}},
	Territory{id: 33, Name: "Portugal", Abbr: "por", edges: []string{"spa"}},
	Territory{id: 34, Name: "Prussia", Abbr: "pru", edges: []string{"sil", "war", "ber", "lvn"}},
	Territory{id: 35, Name: "Rome", Abbr: "rom", edges: []string{"tus", "ven", "apu", "nap"}},
	Territory{id: 36, Name: "Ruhr", Abbr: "ruh", edges: []string{"bel", "bur", "hol", "kie", "mun"}},
	Territory{id: 37, Name: "Rumania", Abbr: "rum", edges: []string{"ser", "sev", "ukr", "bud", "bul", "gal"}},
	Territory{id: 38, Name: "Saint Petersburg", Abbr: "stp", edges: []string{"fin", "lvn", "nwy", "mos"}},
	Territory{id: 39, Name: "Serbia", Abbr: "ser", edges: []string{"tri", "alb", "bud", "bul", "gre", "rum"}},
	Territory{id: 40, Name: "Sevastopol", Abbr: "sev", edges: []string{"ukr", "arm", "mos", "rum"}},
	Territory{id: 41, Name: "Silesia", Abbr: "sil", edges: []string{"war", "ber", "boh", "gal", "mun", "pru"}},
	Territory{id: 42, Name: "Smyrna", Abbr: "smy", edges: []string{"syr", "ank", "arm", "con"}},
	Territory{id: 43, Name: "Spain", Abbr: "spa", edges: []string{"gas", "mar", "por"}},
	Territory{id: 44, Name: "Sweden", Abbr: "swe", edges: []string{"den", "fin", "nwy"}},
	Territory{id: 45, Name: "Syria", Abbr: "syr", edges: []string{"arm", "smy"}},
	Territory{id: 46, Name: "Trieste", Abbr: "tri", edges: []string{"tyr", "ven", "vie", "alb", "bud", "ser"}},
	Territory{id: 47, Name: "Tunis", Abbr: "tun", edges: []string{"naf"}},
	Territory{id: 48, Name: "Tuscany", Abbr: "tus", edges: []string{"ven", "pie", "rom"}},
	Territory{id: 49, Name: "Tyrolia", Abbr: "tyr", edges: []string{"ven", "vie", "boh", "mun", "pie", "tri"}},
	Territory{id: 50, Name: "Ukraine", Abbr: "ukr", edges: []string{"war", "gal", "mos", "rum", "sev"}},
	Territory{id: 51, Name: "Venice", Abbr: "ven", edges: []string{"apu", "pie", "rom", "tri", "tus", "tyr"}},
	Territory{id: 52, Name: "Vienna", Abbr: "vie", edges: []string{"boh", "bud", "gal", "tri", "tyr"}},
	Territory{id: 53, Name: "Wales", Abbr: "wal", edges: []string{"yor", "lon", "lvp"}},
	Territory{id: 54, Name: "Warsaw", Abbr: "war", edges: []string{"gal", "lvn", "mos", "pru", "sil", "ukr"}},
	Territory{id: 55, Name: "Yorkshire", Abbr: "yor", edges: []string{"edi", "lon", "lvn", "wal"}},
}
