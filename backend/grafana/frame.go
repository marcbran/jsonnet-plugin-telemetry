package grafana

import (
	"fmt"
	"strings"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

type frameField struct {
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Labels map[string]string `json:"labels"`
}

type frameSchema struct {
	Fields []frameField `json:"fields"`
}

type frameData struct {
	Values [][]any `json:"values"`
}

type frame struct {
	Schema frameSchema `json:"schema"`
	Data   frameData   `json:"data"`
}

func framesToResult(itemType string, frames []frame) telemetry.QueryResult {
	if itemType == "logql" {
		records := make([]telemetry.LogRecord, 0, len(frames))
		for _, f := range frames {
			records = append(records, frameToRecords(f)...)
		}
		return telemetry.QueryResult{Type: itemType, Records: records}
	}
	series := make([]telemetry.Series, 0, len(frames))
	for _, f := range frames {
		if s, ok := frameToSeries(f); ok {
			series = append(series, s)
		}
	}
	return telemetry.QueryResult{Type: itemType, Series: series}
}

func frameToSeries(f frame) (telemetry.Series, bool) {
	if len(f.Schema.Fields) < 2 || len(f.Data.Values) < 2 {
		return telemetry.Series{}, false
	}
	times := f.Data.Values[0]
	values := f.Data.Values[1]
	points := make([][2]any, len(times))
	for i := range times {
		points[i] = [2]any{times[i], values[i]}
	}
	return telemetry.Series{Labels: f.Schema.Fields[1].Labels, Points: points}, true
}

func frameToRecords(f frame) []telemetry.LogRecord {
	cols := map[string][]any{}
	for i, field := range f.Schema.Fields {
		if i < len(f.Data.Values) {
			cols[field.Name] = f.Data.Values[i]
		}
	}
	times := cols["Time"]
	lines := cols["Line"]
	if times == nil || lines == nil {
		return nil
	}
	labels := cols["labels"]
	ids := cols["id"]
	records := make([]telemetry.LogRecord, len(times))
	for i := range times {
		ts, _ := times[i].(float64)
		fields := map[string]any{}
		if labels != nil {
			if m, ok := labels[i].(map[string]any); ok {
				fields = m
			}
		}
		id := ""
		if ids != nil {
			id = fmt.Sprintf("%v", ids[i])
		}
		records[i] = telemetry.LogRecord{
			Timestamp: ts,
			Body:      strings.TrimRight(fmt.Sprintf("%v", lines[i]), "\r\n"),
			Severity:  severityFrom(fields),
			Fields:    fields,
			ID:        id,
		}
	}
	return records
}

func severityFrom(fields map[string]any) string {
	for _, k := range []string{"detected_level", "level", "severity"} {
		if v, ok := fields[k].(string); ok && v != "" {
			return strings.ToLower(v)
		}
	}
	return ""
}
