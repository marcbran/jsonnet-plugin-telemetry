package telemetry

import (
	"github.com/google/go-jsonnet"
	"github.com/marcbran/jpoet/pkg/jpoet"
)

func Plugin(backend Backend, opts ...jpoet.PluginOption) *jpoet.Plugin {
	return jpoet.NewPlugin("telemetry", []jsonnet.NativeFunction{
		Query(backend),
		Fetch(backend),
	}, opts...)
}
