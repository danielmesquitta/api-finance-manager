package root

import (
	"embed"
)

//go:embed .env
var Env []byte

//go:embed docs/openapi.yaml docs/openapi.json
var StaticFiles embed.FS

//go:embed testdata
var TestData embed.FS
