package grafana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

func TestFramesToResult(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		frames   []frame
		want     telemetry.QueryResult
	}{
		{
			name:     "promql frame becomes series",
			itemType: "promql",
			frames: []frame{{
				Schema: frameSchema{Fields: []frameField{{}, {Labels: map[string]string{"instance": "a"}}}},
				Data:   frameData{Values: [][]any{{float64(0), float64(60000)}, {float64(1.5), float64(2.5)}}},
			}},
			want: telemetry.QueryResult{
				Type: "promql",
				Series: []telemetry.Series{{
					Labels: map[string]string{"instance": "a"},
					Points: [][2]any{{float64(0), float64(1.5)}, {float64(60000), float64(2.5)}},
				}},
			},
		},
		{
			name:     "logql frame becomes records",
			itemType: "logql",
			frames: []frame{{
				Schema: frameSchema{Fields: []frameField{
					{Name: "labels", Type: "other"},
					{Name: "Time", Type: "time"},
					{Name: "Line", Type: "string"},
					{Name: "id", Type: "string"},
				}},
				Data: frameData{Values: [][]any{
					{map[string]any{"app": "api", "detected_level": "ERROR"}, map[string]any{"app": "api"}},
					{float64(0), float64(60000)},
					{"boom\n", "bang"},
					{"id-0", "id-1"},
				}},
			}},
			want: telemetry.QueryResult{
				Type: "logql",
				Records: []telemetry.LogRecord{
					{Timestamp: float64(0), Body: "boom", Severity: "error", Fields: map[string]any{"app": "api", "detected_level": "ERROR"}, ID: "id-0"},
					{Timestamp: float64(60000), Body: "bang", Severity: "", Fields: map[string]any{"app": "api"}, ID: "id-1"},
				},
			},
		},
		{
			name:     "frame with no data field is skipped",
			itemType: "promql",
			frames: []frame{{
				Schema: frameSchema{Fields: []frameField{{}}},
				Data:   frameData{Values: [][]any{{float64(0)}}},
			}},
			want: telemetry.QueryResult{Type: "promql", Series: []telemetry.Series{}},
		},
		{
			name:     "no frames at all",
			itemType: "promql",
			frames:   nil,
			want:     telemetry.QueryResult{Type: "promql", Series: []telemetry.Series{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, framesToResult(tt.itemType, tt.frames))
		})
	}
}
