package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cqroot/prompt"
	"github.com/cqroot/sawmill/internal/script"
	"github.com/cqroot/sawmill/internal/templater"
	"github.com/cqroot/sawmill/internal/templates"
	"github.com/cqroot/sawmill/internal/toml"
)

func Run() error {
	projName, tomlPath, outputDir, err := getPaths()
	if err != nil {
		return err
	}

	rootDir := filepath.Dir(tomlPath)

	templateDir := filepath.Join(rootDir, "template")
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(rootDir, outputDir)
	}

	co, vars, err := getTomlData(tomlPath)
	if err != nil {
		return err
	}

	vars["project_name"] = projName

	tmpl := templater.New(
		templateDir, outputDir, vars, co.IncludePathRules, co.ExcludePathRules)

	fmt.Println()
	fmt.Println("Template path :", templateDir)
	fmt.Println("Output path   :", outputDir)
	fmt.Printf("Variables     : %+v\n", vars)
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

func getTomlPath() (string, error) {
	choice, err := templates.ChooseTemplate()
	if err != nil {
		return "", err
	}

	dataDir, err := templates.DataDir()
	if err != nil {
		return "", err
	}

	tomlPath := filepath.Join(dataDir, choice, "sawmill.toml")
	return tomlPath, nil
}

func getOutputDir() (string, string, error) {
	projName, err := prompt.New().Ask("Your project name:").Input("project")
	if err != nil {
		return "", "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	outputDir := filepath.Join(cwd, projName)
	return projName, outputDir, nil
}

func getPaths() (string, string, string, error) {
	tomlPath, err := getTomlPath()
	if err != nil {
		return "", "", "", err
	}

	projName, outputDir, err := getOutputDir()
	if err != nil {
		return "", "", "", err
	}

	return projName, tomlPath, outputDir, nil
}

func getTomlData(tomlPath string) (*toml.ConfigObject, map[string]string, error) {
	p, err := toml.New(tomlPath)
	if err != nil {
		return nil, nil, err
	}

	co, err := p.Parse()
	if err != nil {
		return nil, nil, err
	}

	ret, err := p.Run(co)
	if err != nil {
		return nil, nil, err
	}

	return co, ret, nil
}
