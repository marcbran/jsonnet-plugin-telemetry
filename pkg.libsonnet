local p = import 'pkg/main.libsonnet';

p.pkg({
  source: 'https://github.com/marcbran/jsonnet-plugin-telemetry',
  repo: 'https://github.com/marcbran/jsonnet.git',
  branch: 'plugin/telemetry',
  path: 'plugin/telemetry',
  target: 'telemetry',
}, |||
  Batched telemetry queries: send a mix of query mechanisms (`promql`, `logql`, ...) in one call, distinguished by a `type` field, and let the configured backend batch as many of them as it can into as few round trips as possible.
|||, {
  query: p.desc(|||
    Runs a batch of query items and returns one result per item, in order. Each item is `{type, datasource, ...params}`, where `type` names the query mechanism (e.g. `promql`, `logql`) and decides how `params` is interpreted and which backend serves it. `datasource` is an opaque identifier meaningful only to that backend.
  |||),
})
