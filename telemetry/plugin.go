package telemetry

import (
	"github.com/google/go-jsonnet"
	"github.com/marcbran/jpoet/pkg/jpoet"
)

func Plugin(name string, backend Backend, opts ...jpoet.PluginOption) *jpoet.Plugin {
	return jpoet.NewPlugin(name, []jsonnet.NativeFunction{
		Query(backend),
	}, opts...)
}
