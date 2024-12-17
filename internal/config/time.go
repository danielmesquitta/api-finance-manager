package config

import "time"

func init() {
	// loc, _ := time.LoadLocation("America/Sao_Paulo")
	time.Local = time.UTC
}
