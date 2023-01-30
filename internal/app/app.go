package app

import (
	"fmt"
	"path/filepath"

	"github.com/cqroot/sawmill/internal/script"
	"github.com/cqroot/sawmill/internal/templater"
	"github.com/cqroot/sawmill/internal/toml"
)

func Run(tomlPath string, outputDir string) error {
	rootDir := filepath.Dir(tomlPath)

	p, err := toml.New(tomlPath)
    if err != nil {
        return err
    }

	co, err := p.Parse()
	if err != nil {
		return err
	}

	ret, err := p.Run(co)
	if err != nil {
		return err
	}

	templateDir := filepath.Join(rootDir, "template")
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(rootDir, outputDir)
	}

	tmpl := templater.New(
		templateDir, outputDir, ret, co.IncludePathRules, co.ExcludePathRules)

	fmt.Println()
	fmt.Println("Template path :", templateDir)
	fmt.Println("Output path   :", outputDir)
	fmt.Printf("Variables     : %+v\n", ret)
	fmt.Println()

	err = tmpl.Execute()
	if err != nil {
		return err
	}

	fmt.Println("")

	for _, scriptPath := range co.Scripts.AfterScripts {
        err = script.Run(scriptPath, outputDir)
        if err != nil {
            return err
        }
	}

	return nil
}
