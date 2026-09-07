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
			name:     "logql frame becomes stream",
			itemType: "logql",
			frames: []frame{{
				Schema: frameSchema{Fields: []frameField{{}, {Labels: map[string]string{"app": "api"}}}},
				Data:   frameData{Values: [][]any{{float64(0)}, {"boom"}}},
			}},
			want: telemetry.QueryResult{
				Type: "logql",
				Streams: []telemetry.Stream{{
					Labels: map[string]string{"app": "api"},
					Lines:  [][2]string{{"0", "boom"}},
				}},
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
