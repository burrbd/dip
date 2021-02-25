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
		if strings.ToUpper(abbr) == terr.Abbr {
			return terr
		}
	}
	return Territory{id: -1} // TODO: why this way?
}

var armyTerritories = []Territory{
	Territory{id: 0, Name: "Albania", Abbr: "ALB", edges: []string{"GRE", "SER", "TRI"}},                               //	Balkans
	Territory{id: 1, Name: "Ankara", Abbr: "ANK", edges: []string{"ARM", "CON", "SMY"}},                                //	TurkeyTRUE
	Territory{id: 2, Name: "Apulia", Abbr: "APU", edges: []string{"NAP", "ROM", "VEN"}},                                //	Italy
	Territory{id: 3, Name: "Armenia", Abbr: "ARM", edges: []string{"SEV", "SMY", "SYR"}},                               //	Turkey
	Territory{id: 4, Name: "Belgium", Abbr: "BEL", edges: []string{"BUR", "HOL", "PIC", "RUR"}},                        //	Low Countries	TRUE
	Territory{id: 5, Name: "Berlin", Abbr: "BER", edges: []string{"KIE", "MUN", "PRU", "SIL"}},                         //	GermanyTRUE
	Territory{id: 6, Name: "Bohemia", Abbr: "BOH", edges: []string{"GAL", "MUN", "SIL", "TYR"}},                        //	Austria
	Territory{id: 7, Name: "Brest", Abbr: "BRE", edges: []string{"GAS", "PAR", "PIC"}},                                 //	FranceTRUE
	Territory{id: 8, Name: "Budapest", Abbr: "BUD", edges: []string{"GAL", "RUM", "SER", "TRI", "VIE"}},                //	AustriaTRUE
	Territory{id: 9, Name: "Bulgaria", Abbr: "BUL", edges: []string{"CON", "GRE", "RUM", "SER"}},                       //	BalkansTRUE
	Territory{id: 10, Name: "Burgundy", Abbr: "BUR", edges: []string{"BEL", "GAS", "MAR", "MUN", "PAR", "PIC", "RUH"}}, //	France
	Territory{id: 11, Name: "Clyde", Abbr: "CLY", edges: []string{"EDI", "LVP"}},                                       //	England
	Territory{id: 12, Name: "Constantinople", Abbr: "CON", edges: []string{"SMY", "ANK", "BUL"}},                       //	TurkeyTRUE
	Territory{id: 13, Name: "Denmark", Abbr: "DEN", edges: []string{"KIE", "SWE"}},                                     //	ScandinaviaTRUE
	Territory{id: 14, Name: "Edinburgh", Abbr: "EDI", edges: []string{"LVP", "YOR", "CLY"}},                            //	EnglandTRUE
	Territory{id: 15, Name: "Finland", Abbr: "FIN", edges: []string{"NWY", "STP", "SWE"}},                              //	Russia
	Territory{id: 16, Name: "Galacia", Abbr: "GAL", edges: []string{"RUM", "SIL", "UKR", "VIE", "WAR", "BOH", "BUD"}},  //	Austria
	Territory{id: 17, Name: "Gascony", Abbr: "GAS", edges: []string{"MAR", "PAR", "SPA", "BRE", "BUR"}},                //	France
	Territory{id: 18, Name: "Greece", Abbr: "GRE", edges: []string{"SER", "ALB", "BUL"}},                               //	BalkansTRUE
	Territory{id: 19, Name: "Holland", Abbr: "HOL", edges: []string{"KIE", "RUH", "BEL"}},                              //	Low Countries	TRUE
	Territory{id: 20, Name: "Kiel", Abbr: "KIE", edges: []string{"MUN", "RUH", "BER", "DEN", "HOL"}},                   //	GermanyTRUE
	Territory{id: 21, Name: "Liverpool", Abbr: "LVP", edges: []string{"WAL", "YOR", "CLY", "EDI"}},                     // or LPL	EnglandTRUE
	Territory{id: 22, Name: "Livonia", Abbr: "LVN", edges: []string{"MOS", "PRU", "STP", "WAR"}},                       // or LVA	Russia
	Territory{id: 23, Name: "London", Abbr: "LON", edges: []string{"WAL", "YOR"}},                                      //	EnglandTRUE
	Territory{id: 24, Name: "Marseilles", Abbr: "MAR", edges: []string{"PIE", "SPA", "BUR", "GAS"}},                    //	FranceTRUE
	Territory{id: 25, Name: "Moscow", Abbr: "MOS", edges: []string{"SEV", "STP", "URK", "WAR", "LVN"}},                 //	RussiaTRUE
	Territory{id: 26, Name: "Munich", Abbr: "MUN", edges: []string{"RUH", "SIL", "TYR", "BER", "BOH", "BUR", "KIE"}},   //	GermanyTRUE
	Territory{id: 27, Name: "Naples", Abbr: "NAP", edges: []string{"ROM", "APU"}},                                      //	ItalyTRUE
	Territory{id: 28, Name: "North Africa", Abbr: "NAF", edges: []string{"TUN"}},                                       //	Africa
	Territory{id: 29, Name: "Norway", Abbr: "NWY", edges: []string{"STP", "SWE", "FIN"}},                               //	ScandinaviaTRUE
	Territory{id: 30, Name: "Paris", Abbr: "PAR", edges: []string{"PIC", "BRE", "BUR", "GAS"}},                         //	FranceTRUE
	Territory{id: 31, Name: "Picardy", Abbr: "PIC", edges: []string{"BEL", "BRE", "BUR", "PAR"}},                       //	France
	Territory{id: 32, Name: "Piedmont", Abbr: "PIE", edges: []string{"TUS", "TYR", "VEN", "MAR"}},                      //	Italy
	Territory{id: 33, Name: "Portugal", Abbr: "POR", edges: []string{"SPA"}},                                           //	IberiaTRUE
	Territory{id: 34, Name: "Prussia", Abbr: "PRU", edges: []string{"SIL", "WAR", "BER", "LVN"}},                       //	Germany
	Territory{id: 35, Name: "Rome", Abbr: "ROM", edges: []string{"TUS", "VEN", "APU", "NAP"}},                          // or RME	ItalyTRUE
	Territory{id: 36, Name: "Ruhr", Abbr: "RUH", edges: []string{"BEL", "BUR", "HOL", "KIE", "MUN"}},                   // or RHR	Germany
	Territory{id: 37, Name: "Rumania", Abbr: "RUM", edges: []string{"SER", "SEV", "UKR", "BUD", "BUL", "GAL"}},         // or RMA	BalkansTRUE
	Territory{id: 38, Name: "Saint Petersburg", Abbr: "STP", edges: []string{"FIN", "LVN", "NWY", "MOS"}},              //	RussiaTRUE
	Territory{id: 39, Name: "Serbia", Abbr: "SER", edges: []string{"TRI", "ALB", "BUD", "BUL", "GRE", "RUM"}},          //	BalkansTRUE
	Territory{id: 40, Name: "Sevastopol", Abbr: "SEV", edges: []string{"UKR", "ARM", "MOS", "RUM"}},                    //	RussiaTRUE
	Territory{id: 41, Name: "Silesia", Abbr: "SIL", edges: []string{"WAR", "BER", "BOH", "GAL", "MUN", "PRU"}},         //	Germany
	Territory{id: 42, Name: "Smyrna", Abbr: "SMY", edges: []string{"SYR", "ANK", "ARM", "CON"}},                        //	TurkeyTRUE
	Territory{id: 43, Name: "Spain", Abbr: "SPA", edges: []string{"GAS", "MAR", "POR"}},                                //	IberiaTRUE
	Territory{id: 44, Name: "Sweden", Abbr: "SWE", edges: []string{"DEN", "FIN", "NWY"}},                               //	ScandinaviaTRUE
	Territory{id: 45, Name: "Syria", Abbr: "SYR", edges: []string{"ARM", "SMY"}},                                       //	Turkey
	Territory{id: 46, Name: "Trieste", Abbr: "TRI", edges: []string{"TYR", "VEN", "VIE", "ALB", "BUD", "SER"}},         //	AustriaTRUE
	Territory{id: 47, Name: "Tunis", Abbr: "TUN", edges: []string{"NAF"}},                                              //	AfricaTRUE
	Territory{id: 48, Name: "Tuscany", Abbr: "TUS", edges: []string{"VEN", "PIE", "ROM"}},                              //	Italy
	Territory{id: 49, Name: "Tyrolia", Abbr: "TYR", edges: []string{"VEN", "VIE", "BOH", "MUN", "PIE", "TRI"}},         //	Austria
	Territory{id: 50, Name: "Ukraine", Abbr: "UKR", edges: []string{"WAR", "GAL", "MOS", "RUM", "SEV"}},                //	Russia
	Territory{id: 51, Name: "Venice", Abbr: "VEN", edges: []string{"APU", "PIE", "ROM", "TRI", "TUS", "TYR"}},          //	ItalyTRUE
	Territory{id: 52, Name: "Vienna", Abbr: "VIE", edges: []string{"BOH", "BUD", "GAL", "TRI", "TYR"}},                 //	AustriaTRUE
	Territory{id: 53, Name: "Wales", Abbr: "WAL", edges: []string{"YOR", "LON", "LVP"}},                                //	England
	Territory{id: 54, Name: "Warsaw", Abbr: "WAR", edges: []string{"GAL", "LVN", "MOS", "PRU", "SIL", "UKR"}},          //	RussiaTRUE
	Territory{id: 55, Name: "Yorkshire", Abbr: "YOR", edges: []string{"EDI", "LON", "LVN", "WAL"}},                     //	England
}
