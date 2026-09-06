package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name       string
		items      []QueryItem
		keyForType map[string]string
		want       []group[string]
		wantErr    string
	}{
		{
			name: "items sharing a key are grouped together, in first-seen order",
			items: []QueryItem{
				{Type: "promql"},
				{Type: "logzioql"},
				{Type: "logql"},
			},
			keyForType: map[string]string{
				"promql":   "grafana",
				"logql":    "grafana",
				"logzioql": "logzio",
			},
			want: []group[string]{
				{key: "grafana", indexes: []int{0, 2}},
				{key: "logzio", indexes: []int{1}},
			},
		},
		{
			name:       "empty input yields no groups",
			items:      nil,
			keyForType: map[string]string{},
			want:       []group[string]{},
		},
		{
			name:       "unknown type is an error",
			items:      []QueryItem{{Type: "promql"}},
			keyForType: map[string]string{},
			wantErr:    `no backend registered for type "promql"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := groupBy(tt.items, tt.keyForType)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeGroupResults(t *testing.T) {
	tests := []struct {
		name           string
		total          int
		indexesByGroup [][]int
		resultsByGroup [][]QueryResult
		want           []QueryResult
		wantErr        string
	}{
		{
			name:           "scatters each group's results back to original positions",
			total:          3,
			indexesByGroup: [][]int{{0, 2}, {1}},
			resultsByGroup: [][]QueryResult{
				{{Type: "promql-a"}, {Type: "promql-c"}},
				{{Type: "logzioql-b"}},
			},
			want: []QueryResult{
				{Type: "promql-a"},
				{Type: "logzioql-b"},
				{Type: "promql-c"},
			},
		},
		{
			name:           "mismatched result count is an error",
			total:          1,
			indexesByGroup: [][]int{{0}},
			resultsByGroup: [][]QueryResult{{}},
			wantErr:        "backend returned 0 results for 1 queries",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeGroupResults(tt.total, tt.indexesByGroup, tt.resultsByGroup)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
