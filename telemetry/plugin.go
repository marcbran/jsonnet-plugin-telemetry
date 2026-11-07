package telemetry

import (
	"github.com/marcbran/jpoet/pkg/jpoet"
)

func Plugin(backend Backend, opts ...jpoet.PluginOption) *jpoet.Plugin {
	return jpoet.NewPlugin("telemetry", nil, append([]jpoet.PluginOption{
		jpoet.WithFunctions(Query(backend), Fetch(backend)),
	}, opts...)...)
}
