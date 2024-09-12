package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"text/template"
)

const ServiceEnvVarName = "ENV"

func GetEnvironmentName() string {
	return os.Getenv(ServiceEnvVarName)
}

func Parse(settings any) error {
	env := GetEnvironmentName()

	configReader, err := readConfig(env)
	if err != nil {
		return err
	}

	return json.NewDecoder(configReader).Decode(settings)
}

// TemplateParams содержит доступные для применения в шаблонах переменные
type TemplateParams struct {
	Env    string
	EnvLow string
}

func readConfig(env string) (io.Reader, error) {
	if len(env) == 0 {
		env = "dev"
	}

	var envTitle string
	if len(env) > 0 {
		envTitle = strings.ToUpper(string(env[0])) + env[1:]
	}

	fileName := configProfile(env) + ".json"
	filePath := "./.config/" + fileName

	tmpl, err := template.New(fileName).ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(buffer, TemplateParams{
		Env:    envTitle,
		EnvLow: strings.ToLower(env),
	})
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func configProfile(env string) string {
	switch env {
	case "", "dev":
		return "dev"
	case "prod":
		return "prod"
	default:
		return "dev"
	}
}
