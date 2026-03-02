package dipmap

import (
	"sort"
	"testing"

	"github.com/cheekybits/is"
)

// testGraph is a simple in-memory graph for unit tests.
type testGraph struct {
	edges map[string][]string
}

func (g *testGraph) Edges(t string) []string { return g.edges[t] }

// sortedNeighborhood calls Neighborhood and returns the result sorted for
// deterministic comparison in tests.
func sortedNeighborhood(g Graph, territory string, n int) []string {
	result := Neighborhood(g, territory, n)
	sort.Strings(result)
	return result
}

func TestNeighborhood_ZeroRadius_ReturnsTerritoryOnly(t *testing.T) {
	is := is.New(t)
	g := &testGraph{edges: map[string][]string{
		"Vienna":   {"Budapest", "Trieste"},
		"Budapest": {"Vienna"},
	}}
	result := sortedNeighborhood(g, "Vienna", 0)
	is.Equal(result, []string{"Vienna"})
}

func TestNeighborhood_RadiusOne_ReturnsTerritoryAndAdjacent(t *testing.T) {
	is := is.New(t)
	g := &testGraph{edges: map[string][]string{
		"Vienna":   {"Budapest", "Trieste"},
		"Budapest": {"Vienna", "Galicia"},
		"Trieste":  {"Vienna"},
	}}
	result := sortedNeighborhood(g, "Vienna", 1)
	is.Equal(result, []string{"Budapest", "Trieste", "Vienna"})
}

func TestNeighborhood_RadiusTwo_ExtendsOneHopFurther(t *testing.T) {
	is := is.New(t)
	g := &testGraph{edges: map[string][]string{
		"Vienna":   {"Budapest"},
		"Budapest": {"Vienna", "Galicia"},
		"Galicia":  {"Budapest", "Warsaw"},
		"Warsaw":   {"Galicia"},
	}}
	result := sortedNeighborhood(g, "Vienna", 2)
	is.Equal(result, []string{"Budapest", "Galicia", "Vienna"})
}

func TestNeighborhood_DisconnectedGraph_ReturnsOnlyReachable(t *testing.T) {
	is := is.New(t)
	g := &testGraph{edges: map[string][]string{
		"Vienna": {"Budapest"},
		"London": {"Edinburgh"}, // disconnected island
	}}
	result := sortedNeighborhood(g, "Vienna", 5)
	is.Equal(result, []string{"Budapest", "Vienna"})
}

func TestNeighborhood_EmptyGraph_ReturnsTerritoryOnly(t *testing.T) {
	is := is.New(t)
	result := sortedNeighborhood(EmptyGraph{}, "Vienna", 3)
	is.Equal(result, []string{"Vienna"})
}

func TestNeighborhood_NoCycles_VisitsEachOnce(t *testing.T) {
	is := is.New(t)
	// A triangle: A-B-C-A
	g := &testGraph{edges: map[string][]string{
		"A": {"B", "C"},
		"B": {"A", "C"},
		"C": {"A", "B"},
	}}
	result := sortedNeighborhood(g, "A", 10)
	is.Equal(result, []string{"A", "B", "C"})
}
