package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type wireConfigData struct {
	PackageName      string
	Imports          []string
	DefaultProviders []string
	Environments     []environmentConfig
}

type environmentConfig struct {
	Name      string
	Providers []string
}

// Use a text/template to generate the final wire.go.
var wireTemplate = `// Code generated by Wire. DO NOT EDIT.

//go:build wireinject
// +build wireinject

package {{.PackageName}}

import (
{{- range .Imports }}
    "{{.}}"
{{- end }}
)

{{- range .Environments }}
// New{{ .Name }} wires up the application in {{ .Name | lower }} mode.
func New{{ .Name }}(
    v *validator.Validator,
    e *config.Env,
    t *testing.T,
) *App {
    wire.Build(
  {{- range .Providers }}
      {{.}}
  {{- end }}
    )
    return &App{}
}
{{- end }}
`

func main() {
	wireConfigFilename := "internal/app/restapi/config.go"          // Source
	wireOutputFilename := "internal/app/restapi/wire_config_gen.go" // Destination

	data, err := parseWireConfig(wireConfigFilename)
	if err != nil {
		panic(err)
	}

	err = generateWireFile(wireOutputFilename, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("wire: wrote %s\n", wireOutputFilename)
}

// parseWireConfig reads wire_config.go line by line
func parseWireConfig(filename string) (*wireConfigData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := &wireConfigData{
		Imports:          []string{},
		DefaultProviders: []string{},
		Environments:     []environmentConfig{},
	}

	scanner := bufio.NewScanner(file)

	inImportBlock := false
	currentSection := ""

	environmentNames := []string{"Dev", "Staging", "Test", "Prod"}

	// Map to track provider sections and their corresponding slices
	providerSections := map[string]*[]string{
		"var providers = []any{": &data.DefaultProviders,
	}

	environments := map[string]*environmentConfig{}
	for _, env := range environmentNames {
		envConfig := environmentConfig{
			Name:      env,
			Providers: []string{},
		}
		environments[env] = &envConfig

		key := fmt.Sprintf("var %sProviders = []any{", strings.ToLower(env))
		providerSections[key] = &envConfig.Providers
	}

	// Keep track of imports in a set to avoid duplicates
	importSet := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Detect package name
		if strings.HasPrefix(line, "package ") && data.PackageName == "" {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				data.PackageName = parts[1]
				continue
			}
		}

		// Detect start/end of import block
		if strings.HasPrefix(line, "import (") {
			inImportBlock = true
			continue
		}
		if inImportBlock && strings.HasPrefix(line, ")") {
			inImportBlock = false
			continue
		}

		// Inside import block
		if inImportBlock {
			impLine := strings.Trim(line, `"`)
			impLine = strings.TrimSpace(impLine)
			if impLine == "" || strings.HasPrefix(impLine, "(") ||
				strings.HasPrefix(impLine, ")") {
				continue
			}
			// De-duplicate imports
			if _, seen := importSet[impLine]; !seen {
				importSet[impLine] = true
				data.Imports = append(data.Imports, impLine)
			}
			continue
		}

		// Check if this line starts a provider section
		for prefix := range providerSections {
			if strings.HasPrefix(line, prefix) {
				currentSection = prefix
				break
			}
		}

		// End of a section
		if line == "}" && currentSection != "" {
			currentSection = ""
			continue
		}

		// Process content within a provider section
		if currentSection != "" && line != currentSection {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && trimmed != "}" {
				// Get the slice pointer for this section and append to it
				if slice, ok := providerSections[currentSection]; ok {
					*slice = append(*slice, trimmed)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Add required imports
	mustHaveImports := []string{
		"github.com/google/wire",
		"github.com/danielmesquitta/api-finance-manager/internal/config",
		"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator",
		"testing",
	}
	for _, imp := range mustHaveImports {
		if _, seen := importSet[imp]; !seen {
			importSet[imp] = true
			data.Imports = append(data.Imports, imp)
		}
	}

	for _, env := range environments {
		providers := append(env.Providers, data.DefaultProviders...)

		data.Environments = append(data.Environments, environmentConfig{
			Name:      env.Name,
			Providers: providers,
		})
	}

	return data, nil
}

func generateWireFile(outputFilename string, data *wireConfigData) error {
	// Add custom function to template for converting to lowercase
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("wire").Funcs(funcMap).Parse(wireTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
