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
	Territory{id: 0, Name: "Albania", Abbr: "alb", edges: []string{"gre", "ser", "tri"}},                               //	Balkans
	Territory{id: 1, Name: "Ankara", Abbr: "ank", edges: []string{"arm", "con", "smy"}},                                //	TurkeyTRUE
	Territory{id: 2, Name: "Apulia", Abbr: "apu", edges: []string{"nap", "rom", "ven"}},                                //	Italy
	Territory{id: 3, Name: "Armenia", Abbr: "arm", edges: []string{"sev", "smy", "syr"}},                               //	Turkey
	Territory{id: 4, Name: "Belgium", Abbr: "bel", edges: []string{"bur", "hol", "pic", "rur"}},                        //	Low Countries	TRUE
	Territory{id: 5, Name: "Berlin", Abbr: "ber", edges: []string{"kie", "mun", "pru", "sil"}},                         //	GermanyTRUE
	Territory{id: 6, Name: "Bohemia", Abbr: "boh", edges: []string{"gal", "mun", "sil", "tyr"}},                        //	Austria
	Territory{id: 7, Name: "Brest", Abbr: "bre", edges: []string{"gas", "par", "pic"}},                                 //	FranceTRUE
	Territory{id: 8, Name: "Budapest", Abbr: "bud", edges: []string{"gal", "rum", "ser", "tri", "vie"}},                //	AustriaTRUE
	Territory{id: 9, Name: "Bulgaria", Abbr: "bul", edges: []string{"con", "gre", "rum", "ser"}},                       //	BalkansTRUE
	Territory{id: 10, Name: "Burgundy", Abbr: "bur", edges: []string{"bel", "gas", "mar", "mun", "par", "pic", "ruh"}}, //	France
	Territory{id: 11, Name: "Clyde", Abbr: "cly", edges: []string{"edi", "lvp"}},                                       //	England
	Territory{id: 12, Name: "Constantinople", Abbr: "con", edges: []string{"smy", "ank", "bul"}},                       //	TurkeyTRUE
	Territory{id: 13, Name: "Denmark", Abbr: "den", edges: []string{"kie", "swe"}},                                     //	ScandinaviaTRUE
	Territory{id: 14, Name: "Edinburgh", Abbr: "edi", edges: []string{"lvp", "yor", "cly"}},                            //	EnglandTRUE
	Territory{id: 15, Name: "Finland", Abbr: "fin", edges: []string{"nwy", "stp", "swe"}},                              //	Russia
	Territory{id: 16, Name: "Galacia", Abbr: "gal", edges: []string{"rum", "sil", "ukr", "vie", "war", "boh", "bud"}},  //	Austria
	Territory{id: 17, Name: "Gascony", Abbr: "gas", edges: []string{"mar", "par", "spa", "bre", "bur"}},                //	France
	Territory{id: 18, Name: "Greece", Abbr: "gre", edges: []string{"ser", "alb", "bul"}},                               //	BalkansTRUE
	Territory{id: 19, Name: "Holland", Abbr: "hol", edges: []string{"kie", "ruh", "bel"}},                              //	Low Countries	TRUE
	Territory{id: 20, Name: "Kiel", Abbr: "kie", edges: []string{"mun", "ruh", "ber", "den", "hol"}},                   //	GermanyTRUE
	Territory{id: 21, Name: "Liverpool", Abbr: "lvp", edges: []string{"wal", "yor", "cly", "edi"}},                     // or lpl	EnglandTRUE
	Territory{id: 22, Name: "Livonia", Abbr: "lvn", edges: []string{"mos", "pru", "stp", "war"}},                       // or lva	Russia
	Territory{id: 23, Name: "London", Abbr: "lon", edges: []string{"wal", "yor"}},                                      //	EnglandTRUE
	Territory{id: 24, Name: "Marseilles", Abbr: "mar", edges: []string{"pie", "spa", "bur", "gas"}},                    //	FranceTRUE
	Territory{id: 25, Name: "Moscow", Abbr: "mos", edges: []string{"sev", "stp", "urk", "war", "lvn"}},                 //	RussiaTRUE
	Territory{id: 26, Name: "Munich", Abbr: "mun", edges: []string{"ruh", "sil", "tyr", "ber", "boh", "bur", "kie"}},   //	GermanyTRUE
	Territory{id: 27, Name: "Naples", Abbr: "nap", edges: []string{"rom", "apu"}},                                      //	ItalyTRUE
	Territory{id: 28, Name: "North Africa", Abbr: "naf", edges: []string{"tun"}},                                       //	Africa
	Territory{id: 29, Name: "Norway", Abbr: "nwy", edges: []string{"stp", "swe", "fin"}},                               //	ScandinaviaTRUE
	Territory{id: 30, Name: "Paris", Abbr: "par", edges: []string{"pic", "bre", "bur", "gas"}},                         //	FranceTRUE
	Territory{id: 31, Name: "Picardy", Abbr: "pic", edges: []string{"bel", "bre", "bur", "par"}},                       //	France
	Territory{id: 32, Name: "Piedmont", Abbr: "pie", edges: []string{"tus", "tyr", "ven", "mar"}},                      //	Italy
	Territory{id: 33, Name: "Portugal", Abbr: "por", edges: []string{"spa"}},                                           //	IberiaTRUE
	Territory{id: 34, Name: "Prussia", Abbr: "pru", edges: []string{"sil", "war", "ber", "lvn"}},                       //	Germany
	Territory{id: 35, Name: "Rome", Abbr: "rom", edges: []string{"tus", "ven", "apu", "nap"}},                          // or rme	ItalyTRUE
	Territory{id: 36, Name: "Ruhr", Abbr: "ruh", edges: []string{"bel", "bur", "hol", "kie", "mun"}},                   // or rhr	Germany
	Territory{id: 37, Name: "Rumania", Abbr: "rum", edges: []string{"ser", "sev", "ukr", "bud", "bul", "gal"}},         // or rma	BalkansTRUE
	Territory{id: 38, Name: "Saint Petersburg", Abbr: "stp", edges: []string{"fin", "lvn", "nwy", "mos"}},              //	RussiaTRUE
	Territory{id: 39, Name: "Serbia", Abbr: "ser", edges: []string{"tri", "alb", "bud", "bul", "gre", "rum"}},          //	BalkansTRUE
	Territory{id: 40, Name: "Sevastopol", Abbr: "sev", edges: []string{"ukr", "arm", "mos", "rum"}},                    //	RussiaTRUE
	Territory{id: 41, Name: "Silesia", Abbr: "sil", edges: []string{"war", "ber", "boh", "gal", "mun", "pru"}},         //	Germany
	Territory{id: 42, Name: "Smyrna", Abbr: "smy", edges: []string{"syr", "ank", "arm", "con"}},                        //	TurkeyTRUE
	Territory{id: 43, Name: "Spain", Abbr: "spa", edges: []string{"gas", "mar", "por"}},                                //	IberiaTRUE
	Territory{id: 44, Name: "Sweden", Abbr: "swe", edges: []string{"den", "fin", "nwy"}},                               //	ScandinaviaTRUE
	Territory{id: 45, Name: "Syria", Abbr: "syr", edges: []string{"arm", "smy"}},                                       //	Turkey
	Territory{id: 46, Name: "Trieste", Abbr: "tri", edges: []string{"tyr", "ven", "vie", "alb", "bud", "ser"}},         //	AustriaTRUE
	Territory{id: 47, Name: "Tunis", Abbr: "tun", edges: []string{"naf"}},                                              //	AfricaTRUE
	Territory{id: 48, Name: "Tuscany", Abbr: "tus", edges: []string{"ven", "pie", "rom"}},                              //	Italy
	Territory{id: 49, Name: "Tyrolia", Abbr: "tyr", edges: []string{"ven", "vie", "boh", "mun", "pie", "tri"}},         //	Austria
	Territory{id: 50, Name: "Ukraine", Abbr: "ukr", edges: []string{"war", "gal", "mos", "rum", "sev"}},                //	Russia
	Territory{id: 51, Name: "Venice", Abbr: "ven", edges: []string{"apu", "pie", "rom", "tri", "tus", "tyr"}},          //	ItalyTRUE
	Territory{id: 52, Name: "Vienna", Abbr: "vie", edges: []string{"boh", "bud", "gal", "tri", "tyr"}},                 //	AustriaTRUE
	Territory{id: 53, Name: "Wales", Abbr: "wal", edges: []string{"yor", "lon", "lvp"}},                                //	England
	Territory{id: 54, Name: "Warsaw", Abbr: "war", edges: []string{"gal", "lvn", "mos", "pru", "sil", "ukr"}},          //	RussiaTRUE
	Territory{id: 55, Name: "Yorkshire", Abbr: "yor", edges: []string{"edi", "lon", "lvn", "wal"}},                     //	England
}
