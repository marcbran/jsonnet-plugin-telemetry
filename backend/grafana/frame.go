package grafana

import (
	"fmt"

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
		streams := make([]telemetry.Stream, 0, len(frames))
		for _, f := range frames {
			if s, ok := frameToStream(f); ok {
				streams = append(streams, s)
			}
		}
		return telemetry.QueryResult{Type: itemType, Streams: streams}
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

func frameToStream(f frame) (telemetry.Stream, bool) {
	if len(f.Schema.Fields) < 2 || len(f.Data.Values) < 2 {
		return telemetry.Stream{}, false
	}
	times := f.Data.Values[0]
	lines := f.Data.Values[1]
	out := make([][2]string, len(times))
	for i := range times {
		out[i] = [2]string{fmt.Sprintf("%v", times[i]), fmt.Sprintf("%v", lines[i])}
	}
	return telemetry.Stream{Labels: f.Schema.Fields[1].Labels, Lines: out}, true
}
