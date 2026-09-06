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

type Stream struct {
	Labels map[string]string
	Lines  [][2]string
}

type QueryResult struct {
	Type    string
	Series  []Series
	Streams []Stream
}

type Backend interface {
	Query(items []QueryItem) ([]QueryResult, error)
}
