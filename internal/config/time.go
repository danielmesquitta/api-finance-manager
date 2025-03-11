package config

import "time"

func setServerTimeZone() {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	time.Local = loc
}
