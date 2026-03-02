package dipmap

// Graph is the adjacency interface used for BFS neighbourhood queries.
// Province names are the string keys returned by Edges.
type Graph interface {
	// Edges returns the names of all provinces directly adjacent to territory.
	Edges(territory string) []string
}

// EmptyGraph is a Graph with no edges; every province is isolated.
// It is used as a default when no real board graph is available.
type EmptyGraph struct{}

// Edges returns nil for every territory.
func (EmptyGraph) Edges(_ string) []string { return nil }

// Neighborhood returns all province names reachable within n hops from
// territory in g. n=0 returns only the territory itself; n=1 returns the
// territory plus all directly adjacent provinces; and so on. The result
// slice has no guaranteed order, but territory is always included.
func Neighborhood(g Graph, territory string, n int) []string {
	visited := map[string]bool{territory: true}
	frontier := []string{territory}
	for i := 0; i < n; i++ {
		var next []string
		for _, t := range frontier {
			for _, e := range g.Edges(t) {
				if !visited[e] {
					visited[e] = true
					next = append(next, e)
				}
			}
		}
		frontier = next
	}
	result := make([]string, 0, len(visited))
	for t := range visited {
		result = append(result, t)
	}
	return result
}
