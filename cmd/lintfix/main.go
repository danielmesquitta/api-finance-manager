package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	dirPaths := []string{
		"./internal/domain/usecase",
		"./internal/domain/entity",
	}

	for _, dirPath := range dirPaths {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(dirPath, file.Name())
			cmd := exec.
				Command(
					"gomodifytags",
					"-w",
					"-all",
					"-quiet",
					"-add-tags",
					"json",
					"-skip-unexported",
					"-sort",
					"-file",
					filePath,
				)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			cmd = exec.Command("formattag", "-file", filePath)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
