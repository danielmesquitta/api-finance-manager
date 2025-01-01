package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	matches, err := filepath.Glob(
		filepath.Join("sql", "migrations", "**", "*.sql"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(matches)
}
