package telemetry

import (
	"fmt"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

func Query(backend Backend) jsonnet.NativeFunction {
	return jsonnet.NativeFunction{
		Name:   "query",
		Params: ast.Identifiers{"items"},
		Func: func(input []any) (any, error) {
			if len(input) != 1 {
				return nil, fmt.Errorf("expected a single items array argument")
			}
			items, err := parseQueryItems(input[0])
			if err != nil {
				return nil, err
			}
			results, err := backend.Query(items)
			if err != nil {
				return nil, err
			}
			if len(results) != len(items) {
				return nil, fmt.Errorf("backend returned %d results for %d queries", len(results), len(items))
			}
			return encodeResults(results), nil
		},
	}
}

func Fetch(backend Backend) jsonnet.NativeFunction {
	return jsonnet.NativeFunction{
		Name:   "fetch",
		Params: ast.Identifiers{"type", "datasource", "id"},
		Func: func(input []any) (any, error) {
			if len(input) != 3 {
				return nil, fmt.Errorf("expected type, datasource, and id arguments")
			}
			typ, ok := input[0].(string)
			if !ok || typ == "" {
				return nil, fmt.Errorf("type must be a non-empty string")
			}
			datasource, ok := input[1].(string)
			if !ok {
				return nil, fmt.Errorf("datasource must be a string")
			}
			id, ok := input[2].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("id must be a non-empty string")
			}
			result, err := backend.Fetch(typ, datasource, id)
			if err != nil {
				return nil, err
			}
			return encodeFetchResult(result), nil
		},
	}
}

func parseQueryItems(raw any) ([]QueryItem, error) {
	list, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("items must be an array")
	}
	items := make([]QueryItem, len(list))
	for i, v := range list {
		item, err := parseQueryItem(v)
		if err != nil {
			return nil, fmt.Errorf("items[%d]: %w", i, err)
		}
		items[i] = item
	}
	return items, nil
}

func parseQueryItem(raw any) (QueryItem, error) {
	m, ok := raw.(map[string]any)
	if !ok {
		return QueryItem{}, fmt.Errorf("must be an object")
	}
	queryType, ok := m["type"].(string)
	if !ok || queryType == "" {
		return QueryItem{}, fmt.Errorf("type must be a non-empty string")
	}
	datasource, _ := m["datasource"].(string)
	params := make(map[string]any, len(m))
	for k, v := range m {
		if k == "type" || k == "datasource" {
			continue
		}
		params[k] = v
	}
	return QueryItem{
		Type:       queryType,
		Datasource: datasource,
		Params:     params,
	}, nil
}

func encodeResults(results []QueryResult) any {
	out := make([]any, len(results))
	for i, r := range results {
		out[i] = encodeResult(r)
	}
	return map[string]any{"results": out}
}

func encodeResult(r QueryResult) map[string]any {
	m := map[string]any{"type": r.Type}
	if r.Series != nil {
		m["series"] = encodeSeries(r.Series)
	}
	if r.Records != nil {
		m["records"] = encodeRecords(r.Records)
	}
	return m
}

func encodeSeries(series []Series) []any {
	out := make([]any, len(series))
	for i, s := range series {
		labels := make(map[string]any, len(s.Labels))
		for k, v := range s.Labels {
			labels[k] = v
		}
		points := make([]any, len(s.Points))
		for j, p := range s.Points {
			points[j] = []any{p[0], p[1]}
		}
		out[i] = map[string]any{"labels": labels, "points": points}
	}
	return out
}

func encodeRecords(records []LogRecord) []any {
	out := make([]any, len(records))
	for i, r := range records {
		out[i] = encodeRecord(r)
	}
	return out
}

func encodeRecord(r LogRecord) map[string]any {
	fields := r.Fields
	if fields == nil {
		fields = map[string]any{}
	}
	return map[string]any{
		"timestamp": r.Timestamp,
		"body":      r.Body,
		"fields":    fields,
		"severity":  r.Severity,
		"id":        r.ID,
	}
}

func encodeFetchResult(r FetchResult) map[string]any {
	m := map[string]any{"type": r.Type}
	if r.Record != nil {
		m["record"] = encodeRecord(*r.Record)
	}
	return m
}
