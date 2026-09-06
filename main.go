package main

import (
	"github.com/marcbran/jsonnet-plugin-telemetry/telemetry"
)

func main() {
	telemetry.Plugin("telemetry", telemetry.NewRouter(map[string]telemetry.Backend{})).Serve()
}
