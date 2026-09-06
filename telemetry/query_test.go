package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQueryItems(t *testing.T) {
	tests := []struct {
		name    string
		raw     any
		want    []QueryItem
		wantErr string
	}{
		{
			name: "splits type and datasource out of params",
			raw: []any{
				map[string]any{"type": "promql", "datasource": "prod", "expr": "up", "instant": false},
			},
			want: []QueryItem{
				{Type: "promql", Datasource: "prod", Params: map[string]any{"expr": "up", "instant": false}},
			},
		},
		{
			name:    "must be an array",
			raw:     map[string]any{},
			wantErr: "items must be an array",
		},
		{
			name: "missing type",
			raw: []any{
				map[string]any{"expr": "up"},
			},
			wantErr: "type must be a non-empty string",
		},
		{
			name: "error is indexed",
			raw: []any{
				map[string]any{"type": "promql"},
				map[string]any{"expr": "up"},
			},
			wantErr: "items[1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseQueryItems(tt.raw)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeResults(t *testing.T) {
	tests := []struct {
		name    string
		results []QueryResult
		want    any
	}{
		{
			name: "series-shaped result",
			results: []QueryResult{{
				Type:   "promql",
				Series: []Series{{Labels: map[string]string{"instance": "a"}, Points: [][2]any{{float64(0), 1.5}}}},
			}},
			want: map[string]any{
				"results": []any{
					map[string]any{
						"type": "promql",
						"series": []any{
							map[string]any{
								"labels": map[string]any{"instance": "a"},
								"points": []any{[]any{float64(0), 1.5}},
							},
						},
					},
				},
			},
		},
		{
			name: "stream-shaped result",
			results: []QueryResult{{
				Type:    "logql",
				Streams: []Stream{{Labels: map[string]string{"app": "api"}, Lines: [][2]string{{"0", "boom"}}}},
			}},
			want: map[string]any{
				"results": []any{
					map[string]any{
						"type": "logql",
						"streams": []any{
							map[string]any{
								"labels": map[string]any{"app": "api"},
								"lines":  []any{[]any{"0", "boom"}},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, encodeResults(tt.results))
		})
	}
}
