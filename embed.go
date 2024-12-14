package root

import (
	_ "embed"
)

//go:embed .env
var Env []byte
