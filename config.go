package address_formatter

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

const componentFileDelimiter = "---"

type template interface{}
type abbreviation map[string]map[string]string
type ConfigFiles struct {
	CountriesPath     string
	ComponentsPath    string
	StateCodesPath    string
	CountryToLangPath string
	CountyCodesPath   string
	CountryCodesPath  string
	AbbreviationFiles string
}
type Config struct {
	Components       map[string][]string
	ComponentAliases map[string]string
	Templates        map[string]interface{}
	StateCodes       map[string]map[string]interface{}
	CountryToLang    map[string]interface{}
	CountyCodes      map[string]interface{}
	CountryCodes     map[string]string
	Abbreviations    map[string]abbreviation
}

// LoadConfig parses the configuration files into a Config structure
func LoadConfig(configFiles ConfigFiles) *Config {
	var config Config

	// todo: check if passing a reference is cleaner than returning values for non-interface types
	config.Components, config.ComponentAliases = loadComponentsAndAliasesConfig(configFiles.ComponentsPath)
	config.Abbreviations = loadAbbreviationConfig(configFiles.AbbreviationFiles)
	config.CountryCodes = loadCountryCodesConfig(configFiles.CountryCodesPath)
	loadConfig(configFiles.CountriesPath, &config.Templates)
	loadConfig(configFiles.StateCodesPath, &config.StateCodes)
	loadConfig(configFiles.CountryToLangPath, &config.CountryToLang)
	loadConfig(configFiles.CountyCodesPath, &config.CountyCodes)

	return &config
}

func getFileContent(path string) string {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatalf("Could not read %s: %v", path, err)
	}

	return string(content)
}

func loadCountryCodesConfig(countryCodesPath string) map[string]string {
	fileContent := getFileContent(countryCodesPath)
	fileContent = regexp.MustCompile(` #`).ReplaceAllString(fileContent, "")

	var countryCodes map[string]string

	err := yaml.Unmarshal([]byte(fileContent), &countryCodes)

	if err != nil {
		log.Fatalf("Could not load countries config file: %v", err)
	}

	return countryCodes
}

func loadComponentsAndAliasesConfig(componentsPath string) (map[string][]string, map[string]string) {
	componentFileContent := getFileContent(componentsPath)
	componentParts := strings.Split(componentFileContent, componentFileDelimiter)
	componentAliases := make(map[string]string)
	components := make(map[string][]string)

	for _, componentPart := range componentParts {
		var component struct {
			Name    string   `yaml:"Name"`
			Aliases []string `yaml:"aliases"`
		}

		err := yaml.Unmarshal([]byte(componentPart), &component)

		if err != nil {
			log.Fatalf("Could not load components config file: %v", err)
		}

		if len(component.Aliases) > 0 {
			for _, alias := range component.Aliases {
				componentAliases[alias] = component.Name
			}
		}

		components[component.Name] = component.Aliases
	}

	return components, componentAliases
}

func loadAbbreviationConfig(abbreviationPath string) map[string]abbreviation {
	abbreviationFiles, err := filepath.Glob(abbreviationPath)
	abbreviations := make(map[string]abbreviation, len(abbreviationFiles))

	if err != nil {
		log.Fatalf("Could not load abbreviation config file: %v", err)
	}

	for _, filePath := range abbreviationFiles {
		var abbreviation abbreviation
		loadConfig(filePath, &abbreviation)

		fileBase := filepath.Base(filePath)
		language := fileBase[0 : len(fileBase)-len(filepath.Ext(fileBase))]

		abbreviations[language] = abbreviation
	}

	return abbreviations
}

func loadConfig(path string, config interface{}) {
	fileContent := getFileContent(path)
	err := yaml.Unmarshal([]byte(fileContent), config)

	if err != nil {
		log.Fatalf("Could not load %s: %v", path, err)
	}
}
