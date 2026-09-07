package grafana

import (
	"fmt"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

type resolvedItem struct {
	index      int
	refID      string
	itemType   string
	datasource Datasource
	expr       string
	instant    bool
	baseURL    string
	from       string
	to         string
}

func resolveItems(items []telemetry.QueryItem, datasources map[string]Datasource) ([]resolvedItem, error) {
	resolved := make([]resolvedItem, len(items))
	for i, item := range items {
		ds, ok := datasources[item.Datasource]
		if !ok {
			return nil, fmt.Errorf("items[%d]: unknown datasource %q", i, item.Datasource)
		}
		expr, ok := item.Params["expr"].(string)
		if !ok || expr == "" {
			return nil, fmt.Errorf("items[%d]: expr must be a non-empty string", i)
		}
		from, ok := item.Params["from"].(string)
		if !ok || from == "" {
			return nil, fmt.Errorf("items[%d]: from must be a non-empty string", i)
		}
		to, ok := item.Params["to"].(string)
		if !ok || to == "" {
			return nil, fmt.Errorf("items[%d]: to must be a non-empty string", i)
		}
		instant, _ := item.Params["instant"].(bool)
		resolved[i] = resolvedItem{
			index:      i,
			refID:      refID(i),
			itemType:   item.Type,
			datasource: ds,
			expr:       expr,
			instant:    instant,
			baseURL:    ds.BaseURL,
			from:       from,
			to:         to,
		}
	}
	return resolved, nil
}

func refID(i int) string {
	return string(rune('A' + i))
}

type requestKey struct {
	baseURL string
	from    string
	to      string
}

type requestGroup struct {
	baseURL string
	from    string
	to      string
	items   []resolvedItem
}

func groupByRequest(items []resolvedItem) []requestGroup {
	byKey := map[requestKey]*requestGroup{}
	var order []requestKey
	for _, item := range items {
		key := requestKey{baseURL: item.baseURL, from: item.from, to: item.to}
		g, ok := byKey[key]
		if !ok {
			g = &requestGroup{baseURL: key.baseURL, from: key.from, to: key.to}
			byKey[key] = g
			order = append(order, key)
		}
		g.items = append(g.items, item)
	}
	groups := make([]requestGroup, len(order))
	for i, key := range order {
		groups[i] = *byKey[key]
	}
	return groups
}
