local p = import 'pkg/main.libsonnet';

p.ex({
  query: p.ex([{
    name: 'range query',
    inputs: [[{ type: 'promql', datasource: 'prod', expr: 'up', instant: false }]],
  }, {
    name: 'instant query',
    inputs: [[{ type: 'promql', datasource: 'prod', expr: 'up', instant: true }]],
  }, {
    name: 'mixed batch of metric and log queries',
    inputs: [[
      { type: 'promql', datasource: 'prod', expr: 'up', instant: false },
      { type: 'logql', datasource: 'prod', expr: '{app="api"} |= "error"', limit: 100 },
    ]],
  }]),
})
