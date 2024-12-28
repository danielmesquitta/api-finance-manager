package root

import (
	"embed"
	_ "embed"
)

//go:embed .env
var Env []byte

//go:embed docs/openapi.yaml docs/openapi.json
var StaticFiles embed.FS

//go:embed mock
var MockData embed.FS
