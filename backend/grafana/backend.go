package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

type Datasource struct {
	BaseURL string
	UID     string
	Type    string
}

type Auth func(*http.Request) error

type Option func(*Backend)

func WithAuth(auth Auth) Option {
	return func(b *Backend) { b.auth = auth }
}

func WithHTTPClient(client *http.Client) Option {
	return func(b *Backend) { b.client = client }
}

type Backend struct {
	auth        Auth
	datasources map[string]Datasource
	client      *http.Client
}

func New(datasources map[string]Datasource, opts ...Option) *Backend {
	b := &Backend{
		datasources: datasources,
		client:      &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *Backend) Query(items []telemetry.QueryItem) ([]telemetry.QueryResult, error) {
	resolved, err := resolveItems(items, b.datasources)
	if err != nil {
		return nil, err
	}
	groups := groupByRequest(resolved)

	results := make([]telemetry.QueryResult, len(items))
	for _, g := range groups {
		framesByRefID, err := b.runQuery(g)
		if err != nil {
			return nil, err
		}
		for _, ri := range g.items {
			results[ri.index] = framesToResult(ri.itemType, framesByRefID[ri.refID])
		}
	}
	return results, nil
}

type dsQueryTarget struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

type dsQuery struct {
	RefID      string        `json:"refId"`
	Datasource dsQueryTarget `json:"datasource"`
	Expr       string        `json:"expr"`
	Instant    bool          `json:"instant"`
	Range      bool          `json:"range"`
}

type dsQueryRequest struct {
	Queries []dsQuery `json:"queries"`
	From    string    `json:"from"`
	To      string    `json:"to"`
}

type dsQueryResult struct {
	Frames []frame `json:"frames"`
}

type dsQueryResponse struct {
	Results map[string]dsQueryResult `json:"results"`
}

func (b *Backend) runQuery(g requestGroup) (map[string][]frame, error) {
	queries := make([]dsQuery, len(g.items))
	for i, ri := range g.items {
		queries[i] = dsQuery{
			RefID:      ri.refID,
			Datasource: dsQueryTarget{Type: ri.datasource.Type, UID: ri.datasource.UID},
			Expr:       ri.expr,
			Instant:    ri.instant,
			Range:      !ri.instant,
		}
	}
	reqBody, err := json.Marshal(dsQueryRequest{Queries: queries, From: g.from, To: g.to})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(g.baseURL, "/")+"/api/ds/query", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if b.auth != nil {
		if err := b.auth(req); err != nil {
			return nil, err
		}
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("grafana request failed: %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	var parsed dsQueryResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	out := make(map[string][]frame, len(parsed.Results))
	for refID, r := range parsed.Results {
		out[refID] = r.Frames
	}
	return out, nil
}
