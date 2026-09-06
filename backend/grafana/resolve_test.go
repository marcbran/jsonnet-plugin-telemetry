package grafana

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

func TestResolveItems(t *testing.T) {
	prod := Datasource{BaseURL: "https://grafana.example.com", UID: "abc", Type: "prometheus"}

	tests := []struct {
		name        string
		items       []telemetry.QueryItem
		datasources map[string]Datasource
		want        []resolvedItem
		wantErr     string
	}{
		{
			name: "resolves datasource and params",
			items: []telemetry.QueryItem{
				{Type: "promql", Datasource: "prod", Params: map[string]any{"expr": "up", "from": "0", "to": "100", "instant": true}},
			},
			datasources: map[string]Datasource{"prod": prod},
			want: []resolvedItem{
				{index: 0, refID: "A", itemType: "promql", datasource: prod, expr: "up", instant: true, baseURL: prod.BaseURL, from: "0", to: "100"},
			},
		},
		{
			name: "instant defaults to false",
			items: []telemetry.QueryItem{
				{Type: "promql", Datasource: "prod", Params: map[string]any{"expr": "up", "from": "0", "to": "100"}},
			},
			datasources: map[string]Datasource{"prod": prod},
			want: []resolvedItem{
				{index: 0, refID: "A", itemType: "promql", datasource: prod, expr: "up", instant: false, baseURL: prod.BaseURL, from: "0", to: "100"},
			},
		},
		{
			name:        "unknown datasource",
			items:       []telemetry.QueryItem{{Type: "promql", Datasource: "missing", Params: map[string]any{}}},
			datasources: map[string]Datasource{"prod": prod},
			wantErr:     `unknown datasource "missing"`,
		},
		{
			name:        "missing expr",
			items:       []telemetry.QueryItem{{Type: "promql", Datasource: "prod", Params: map[string]any{"from": "0", "to": "100"}}},
			datasources: map[string]Datasource{"prod": prod},
			wantErr:     "expr must be a non-empty string",
		},
		{
			name:        "missing from",
			items:       []telemetry.QueryItem{{Type: "promql", Datasource: "prod", Params: map[string]any{"expr": "up", "to": "100"}}},
			datasources: map[string]Datasource{"prod": prod},
			wantErr:     "from must be a non-empty string",
		},
		{
			name:        "missing to",
			items:       []telemetry.QueryItem{{Type: "promql", Datasource: "prod", Params: map[string]any{"expr": "up", "from": "0"}}},
			datasources: map[string]Datasource{"prod": prod},
			wantErr:     "to must be a non-empty string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveItems(tt.items, tt.datasources)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRefID(t *testing.T) {
	tests := []struct {
		index int
		want  string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, refID(tt.index))
	}
}

func TestGroupByRequest(t *testing.T) {
	tests := []struct {
		name  string
		items []resolvedItem
		want  []requestGroup
	}{
		{
			name: "groups by baseURL, from and to together",
			items: []resolvedItem{
				{index: 0, baseURL: "https://a", from: "0", to: "100"},
				{index: 1, baseURL: "https://b", from: "0", to: "100"},
				{index: 2, baseURL: "https://a", from: "0", to: "100"},
				{index: 3, baseURL: "https://a", from: "0", to: "200"},
			},
			want: []requestGroup{
				{
					baseURL: "https://a", from: "0", to: "100",
					items: []resolvedItem{
						{index: 0, baseURL: "https://a", from: "0", to: "100"},
						{index: 2, baseURL: "https://a", from: "0", to: "100"},
					},
				},
				{
					baseURL: "https://b", from: "0", to: "100",
					items: []resolvedItem{
						{index: 1, baseURL: "https://b", from: "0", to: "100"},
					},
				},
				{
					baseURL: "https://a", from: "0", to: "200",
					items: []resolvedItem{
						{index: 3, baseURL: "https://a", from: "0", to: "200"},
					},
				},
			},
		},
		{
			name:  "empty input yields no groups",
			items: nil,
			want:  []requestGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, groupByRequest(tt.items))
		})
	}
}
