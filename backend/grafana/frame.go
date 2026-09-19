package grafana

import (
	"fmt"
	"strings"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

type frameField struct {
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
	if len(f.Schema.Fields) < 2 || len(f.Data.Values) < 2 {
		return nil
	}
	times := f.Data.Values[0]
	lines := f.Data.Values[1]
	labels := f.Schema.Fields[1].Labels
	fields := make(map[string]any, len(labels))
	for k, v := range labels {
		fields[k] = v
	}
	severity := severityFrom(labels)
	records := make([]telemetry.LogRecord, len(times))
	for i := range times {
		ts, _ := times[i].(float64)
		records[i] = telemetry.LogRecord{
			Timestamp: ts,
			Body:      fmt.Sprintf("%v", lines[i]),
			Severity:  severity,
			Fields:    fields,
		}
	}
	return records
}

func severityFrom(labels map[string]string) string {
	for _, k := range []string{"detected_level", "level", "severity"} {
		if v, ok := labels[k]; ok && v != "" {
			return strings.ToLower(v)
		}
	}
	return ""
}
