package telemetry

type QueryItem struct {
	Type       string
	Datasource string
	Params     map[string]any
}

type Series struct {
	Labels map[string]string
	Points [][2]any
}

type LogRecord struct {
	Timestamp float64
	Body      string
	Severity  string
	Fields    map[string]any
	ID        string
}

type QueryResult struct {
	Type    string
	Series  []Series
	Records []LogRecord
}

type FetchResult struct {
	Type   string
	Record *LogRecord
}

type Backend interface {
	Query(items []QueryItem) ([]QueryResult, error)
	Fetch(typ string, datasource string, id string) (FetchResult, error)
}
