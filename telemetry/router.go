package telemetry

import "fmt"

type Router struct {
	backends map[string]Backend
}

func NewRouter(backends map[string]Backend) *Router {
	return &Router{backends: backends}
}

func (r *Router) Query(items []QueryItem) ([]QueryResult, error) {
	groups, err := groupBy(items, r.backends)
	if err != nil {
		return nil, err
	}
	indexesByGroup := make([][]int, len(groups))
	resultsByGroup := make([][]QueryResult, len(groups))
	for i, g := range groups {
		sub := make([]QueryItem, len(g.indexes))
		for j, idx := range g.indexes {
			sub[j] = items[idx]
		}
		subResults, err := g.key.Query(sub)
		if err != nil {
			return nil, err
		}
		indexesByGroup[i] = g.indexes
		resultsByGroup[i] = subResults
	}
	return mergeGroupResults(len(items), indexesByGroup, resultsByGroup)
}

type group[K comparable] struct {
	key     K
	indexes []int
}

func groupBy[K comparable](items []QueryItem, keyForType map[string]K) ([]group[K], error) {
	byKey := map[K]*group[K]{}
	var order []K
	for i, item := range items {
		key, ok := keyForType[item.Type]
		if !ok {
			return nil, fmt.Errorf("no backend registered for type %q", item.Type)
		}
		g, ok := byKey[key]
		if !ok {
			g = &group[K]{key: key}
			byKey[key] = g
			order = append(order, key)
		}
		g.indexes = append(g.indexes, i)
	}
	groups := make([]group[K], len(order))
	for i, key := range order {
		groups[i] = *byKey[key]
	}
	return groups, nil
}

func mergeGroupResults(total int, indexesByGroup [][]int, resultsByGroup [][]QueryResult) ([]QueryResult, error) {
	results := make([]QueryResult, total)
	for gi, indexes := range indexesByGroup {
		sub := resultsByGroup[gi]
		if len(sub) != len(indexes) {
			return nil, fmt.Errorf("backend returned %d results for %d queries", len(sub), len(indexes))
		}
		for j, idx := range indexes {
			results[idx] = sub[j]
		}
	}
	return results, nil
}
