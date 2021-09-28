package addrFmt

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

// testing relies on testcases by OpenCageData (use testcases from submodule by executing copy-testcases.cmd or create your own)
func TestAddressFormatterTestSuite(t *testing.T) {
	suite.Run(t, new(FormatterTestSuite))
}

type FormatterTestSuite struct {
	suite.Suite
	Config    *Config
	TestCases []*TestCase
}
type TestCase struct {
	Name           string
	Address        *Address
	ExpectedOutput string
}

func (suite *FormatterTestSuite) SetupTest() {
	// Load Config
	suite.Config = LoadConfig(ConfigFiles{
		CountriesPath:     "templates/countries/worldwide.yaml",
		ComponentsPath:    "templates/components.yaml",
		StateCodesPath:    "templates/state_codes.yaml",
		CountryToLangPath: "templates/country2lang.yaml",
		CountyCodesPath:   "templates/county_codes.yaml",
		CountryCodesPath:  "templates/country_codes.yaml",
		AbbreviationFiles: "templates/abbreviations/*.yaml",
	})
	suite.Config.OutputFormat = PostalFormat
	suite.Config.UnknownAsAttention = true

	// Load TestCases
	testFiles, _ := filepath.Glob("testcases/*.yaml")

	for _, filePath := range testFiles {
		fileContent, _ := ioutil.ReadFile(filePath)
		fileContentStr := string(fileContent)

		var testCaseParts []string
		if strings.Count(fileContentStr, "---") > 1 {
			testCaseParts = strings.Split(fileContentStr, "---")[1:]
		} else {
			testCaseParts = append(testCaseParts, fileContentStr)
		}
		for _, testCasePart := range testCaseParts {
			var parsedFileContent struct {
				Components map[string]string `yaml:"components"`
				Expected   string            `yaml:"expected"`
			}

			yaml.Unmarshal([]byte(testCasePart), &parsedFileContent)

			testCase := new(TestCase)

			testCase.ExpectedOutput = parsedFileContent.Expected
			testCase.Name = filepath.Base(filePath)
			var err error
			testCase.Address, err = GetFixedAddress(parsedFileContent.Components, suite.Config) // OpenCageData's expected output relies on fixed addresses
			suite.NoError(err, "GetFixedAddress failed for test case %s", testCase.Name)
			suite.TestCases = append(suite.TestCases, testCase)
		}
	}
}

func (suite *FormatterTestSuite) TestAddressFormatter() {
	for _, testCase := range suite.TestCases {
		formattedAddress, err := FormatAddress(testCase.Address, suite.Config)

		suite.NoError(err)
		suite.Equalf(testCase.ExpectedOutput, formattedAddress, "Test case file: %s", testCase.Name)
	}
}
