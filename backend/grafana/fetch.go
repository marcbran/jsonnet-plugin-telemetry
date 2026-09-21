package grafana

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

func (b *Backend) Fetch(typ string, datasource string, id string) (telemetry.FetchResult, error) {
	if typ != "logql" {
		return telemetry.FetchResult{}, fmt.Errorf("fetch not supported for type %q", typ)
	}
	ds, ok := b.datasources[datasource]
	if !ok {
		return telemetry.FetchResult{}, fmt.Errorf("unknown datasource %q", datasource)
	}
	nanos := id
	if i := strings.IndexByte(id, '_'); i >= 0 {
		nanos = id[:i]
	}
	ns, err := strconv.ParseInt(nanos, 10, 64)
	if err != nil {
		return telemetry.FetchResult{}, fmt.Errorf("invalid record id %q", id)
	}
	from := strconv.FormatInt(ns/1_000_000, 10)
	to := strconv.FormatInt(ns/1_000_000+1, 10)
	g := requestGroup{
		baseURL: ds.BaseURL,
		from:    from,
		to:      to,
		items: []resolvedItem{{
			index:      0,
			refID:      refID(0),
			itemType:   "logql",
			datasource: ds,
			expr:       `{service_name=~".+"}`,
			baseURL:    ds.BaseURL,
			from:       from,
			to:         to,
		}},
	}
	framesByRefID, err := b.runQuery(g)
	if err != nil {
		return telemetry.FetchResult{}, err
	}
	result := framesToResult("logql", framesByRefID[refID(0)])
	for i := range result.Records {
		if result.Records[i].ID == id {
			record := result.Records[i]
			return telemetry.FetchResult{Type: "logql", Record: &record}, nil
		}
	}
	return telemetry.FetchResult{}, fmt.Errorf("record %q not found", id)
}
