package address_formatter

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestUtilityTestSuite(t *testing.T) {
	suite.Run(t, new(UtilityTestSuite))
}

type UtilityTestSuite struct {
	suite.Suite
	Config *Config
}

func (suite *UtilityTestSuite) SetupTest() {
	suite.Config = LoadConfig(ConfigFiles{
		CountriesPath:     "templates/countries/worldwide.yaml",
		ComponentsPath:    "templates/components.yaml",
		StateCodesPath:    "templates/state_codes.yaml",
		CountryToLangPath: "templates/country2lang.yaml",
		CountyCodesPath:   "templates/county_codes.yaml",
		CountryCodesPath:  "templates/country_codes.yaml",
		AbbreviationFiles: "templates/abbreviations/*.yaml",
	})
}

func (suite *UtilityTestSuite) TestAddressToMap() {
	address := Address{
		Attention:     "Lorem",
		House:         "dolor",
		HouseNumber:   "ipsum",
		Road:          "sit",
		Hamlet:        "amet",
		Village:       "consetetur",
		Neighbourhood: "estitius",
		PostalCity:    "sadipscing",
		City:          "elitr",
		CityDistrict:  "vero",
		Municipality:  "sed",
		County:        "diam",
		CountyCode:    "nonumy",
		StateDistrict: "eirmod",
		Postcode:      "invidunt",
		State:         "ut",
		StateCode:     "tempor",
		Region:        "gubergren",
		Suburb:        "labore",
		Town:          "sanctus",
		Island:        "magna",
		Archipelago:   "aliquyam",
		Country:       "erat",
		CountryCode:   "wisi",
		Continent:     "voluptua",
	}

	expectedAddressMap := addressMap{
		"attention":      "Lorem",
		"house_number":   "ipsum",
		"house":          "dolor",
		"road":           "sit",
		"hamlet":         "amet",
		"village":        "consetetur",
		"postal_city":    "sadipscing",
		"city":           "elitr",
		"city_district":  "vero",
		"municipality":   "sed",
		"neighbourhood":  "estitius",
		"county":         "diam",
		"county_code":    "nonumy",
		"state_district": "eirmod",
		"state":          "ut",
		"state_code":     "tempor",
		"postcode":       "invidunt",
		"suburb":         "labore",
		"region":         "gubergren",
		"town":           "sanctus",
		"island":         "magna",
		"archipelago":    "aliquyam",
		"country":        "erat",
		"country_code":   "wisi",
		"continent":      "voluptua",
	}

	addressMap, err := addressToMap(&address)

	suite.NoError(err)
	suite.Equal(addressMap, expectedAddressMap)
}

func (suite *UtilityTestSuite) TestMapToAddress() {
	addressMap := addressMap{
		"attention":      "Lorem",
		"house_number":   "ipsum",
		"house":          "dolor",
		"road":           "sit",
		"hamlet":         "amet",
		"village":        "consetetur",
		"neighbourhood":  "estitius",
		"postal_city":    "sadipscing",
		"city":           "elitr",
		"municipality":   "sed",
		"county":         "diam",
		"county_code":    "nonumy",
		"state_district": "eirmod",
		"state":          "ut",
		"state_code":     "tempor",
		"postcode":       "invidunt",
		"suburb":         "labore",
		"region":         "gubergren",
		"town":           "sanctus",
		"island":         "magna",
		"archipelago":    "aliquyam",
		"country":        "erat",
		"country_code":   "wisi",
		"continent":      "voluptua",
	}
	expectedAddress := &Address{
		Attention:     "Lorem",
		House:         "dolor",
		HouseNumber:   "ipsum",
		Road:          "sit",
		Hamlet:        "amet",
		Village:       "consetetur",
		Neighbourhood: "estitius",
		PostalCity:    "sadipscing",
		City:          "elitr",
		Municipality:  "sed",
		County:        "diam",
		CountyCode:    "nonumy",
		StateDistrict: "eirmod",
		Postcode:      "invidunt",
		State:         "ut",
		StateCode:     "tempor",
		Region:        "gubergren",
		Suburb:        "labore",
		Town:          "sanctus",
		Island:        "magna",
		Archipelago:   "aliquyam",
		Country:       "erat",
		CountryCode:   "wisi",
		Continent:     "voluptua",
	}

	address := MapToAddress(addressMap, suite.Config)

	suite.Equal(address, expectedAddress)
}

func (suite *UtilityTestSuite) TestMapToAddressAliases() {
	// leave some fields out that should get filled with the alias things
	addressMap := addressMap{
		"street_number":             "Lorem", // for house_number
		"building":                  "dolor", // for house
		"street":                    "ipsum", // for road
		"croft":                     "sit",   // for hamlet
		"village":                   "amet",
		"city_district":             "consetetur", // for neighbourhood
		"postal_city":               "estitius",
		"town":                      "sadipscing", // for city
		"local_administrative_area": "elitr",      // for municipality
		"department":                "sed",        // for country
		"state_district":            "diam",
		"partial_postcode":          "nonumy", // for postcode
		"province":                  "eirmod", // for state
		"region":                    "invidunt",
		"island":                    "ut",
		"archipelago":               "tempor",
		"country_name":              "gubergren", // for country
		"country_code":              "labore",
		"continent":                 "sanctus",
	}
	expectedAddress := &Address{
		HouseNumber:   "Lorem",
		House:         "dolor",
		Road:          "ipsum",
		Hamlet:        "sit",
		Village:       "amet",
		Neighbourhood: "consetetur",
		PostalCity:    "estitius",
		City:          "sadipscing",
		CityDistrict:  "consetetur",
		Town:          "sadipscing",
		Municipality:  "elitr",
		County:        "sed",
		StateDistrict: "diam",
		Postcode:      "nonumy",
		State:         "eirmod",
		Region:        "invidunt",
		Island:        "ut",
		Archipelago:   "tempor",
		Country:       "gubergren",
		CountryCode:   "labore",
		Continent:     "sanctus",
	}

	address := MapToAddress(addressMap, suite.Config)

	suite.Equal(address, expectedAddress)
}
