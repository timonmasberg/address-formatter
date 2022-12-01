package addrFmt

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// GetFixedAddress Fixes postcode/country, adds missing state/county/country-code and applies template replacements
// GetFixedAddress Entrypoint for data such as from osm
func GetFixedAddress(addressMap addressMap, config *Config) (*Address, error) {
	addressMap["country_code"] = getFixedCountryCode(addressMap["country_code"])
	// set template before applying aliases to ensure country template is being used
	template := findTemplate(addressMap["country_code"], config.Templates)

	addressMap["country_code"] = determineCountryCode(addressMap["country_code"], template)

	newCountry, hasChangeCountry := template.(map[string]interface{})["change_country"].(string)
	if hasChangeCountry {
		var err error
		addressMap["country"], err = determineCountry(addressMap, newCountry)

		if err != nil {
			return nil, err
		}
	}

	addComponent, hasAddComponent := template.(map[string]interface{})["add_component"].(string)
	if hasAddComponent {
		addTemplateComponents(addressMap, addComponent)
	}

	applySpecialCases(addressMap)

	replacements, hasReplacements := template.(map[string]interface{})["replace"].([]interface{})
	if hasReplacements {
		applyReplacements(addressMap, replacements)
	}

	applyUrlCleanup(addressMap)

	address := MapToAddress(addressMap, config.ComponentAliases, config.UnknownAsAttention)
	cleanupAddress(address, config)

	return address, nil
}

var washingtonCheck = regexp.MustCompile(`(?i)^washington,? d\.?c\.?`)
var postcodeRangeCheck = regexp.MustCompile(`^(\d{5}),\d{5}`)
var multiplePostcodeCheck = regexp.MustCompile(`\d+;\d+`)

// this function is mostly ported from @fragaria/Address-formatter
func cleanupAddress(address *Address, config *Config) {
	if address.Country != "" && address.State != "" {
		if _, err := strconv.ParseInt(address.Country, 10, 64); err == nil {
			address.Country = address.State
			address.State = ""
		}
	}

	if address.StateCode == "" && address.State != "" {
		address.StateCode = getStateCode(address.State, address.CountryCode, config.StateCodes)

		if washingtonCheck.MatchString(address.State) {
			address.StateCode = "DC"
			address.State = "District of Columbia"
			address.City = "Washington"
		}
	}

	if address.CountyCode == "" && address.County != "" {
		address.CountyCode = getCountyCode(address.County, address.CountryCode, config.CountyCodes)
	}

	if address.Postcode != "" {
		if len(address.Postcode) > 20 || multiplePostcodeCheck.MatchString(address.Postcode) {
			address.Postcode = ""

		} else // postcode range
		if matches := postcodeRangeCheck.FindStringSubmatch(address.Postcode); len(matches) > 0 {
			address.Postcode = matches[1]
		}
	}
}

func applyReplacements(address addressMap, replacements []interface{}) {
	for key, value := range address {
		for _, replacement := range replacements {
			r, err := regexp.Compile("^" + key + "=")

			if err != nil {
				log.Printf("Could not replace due to bad regexp: %v", err)
			}

			replacementSrc := replacement.([]interface{})[0].(string)
			replacementVal := replacement.([]interface{})[1].(string)

			if r.MatchString(replacementSrc) {
				valueAfterReplacement := r.ReplaceAllString(replacementSrc, "")

				if value == valueAfterReplacement {
					address[key] = replacementVal
				}
			} else {
				r, err = regexp.Compile(replacementSrc)

				if err != nil {
					log.Printf("Could not replace due to bad regexp: %v", err)
				}

				address[key] = r.ReplaceAllString(address[key], replacementVal)
			}
		}
	}
}

func getFixedCountryCode(countryCode string) string {
	if len(countryCode) != 2 {
		return ""
	}

	countryCode = strings.ToUpper(countryCode)
	// special case for UK
	if countryCode == "UK" {
		countryCode = "GB"
	}

	return countryCode
}

func determineCountryCode(countryCode string, template template) string {
	if len(countryCode) != 2 {
		return ""
	}

	newCountryCode, hasUseCountry := template.(map[string]interface{})["use_country"].(string)

	if hasUseCountry {
		countryCode = strings.ToUpper(newCountryCode)
	}

	if alias, hasAlias := getCountryCodeAlias(countryCode); hasAlias {
		countryCode = alias
	}

	return countryCode
}

func addTemplateComponents(addressMap addressMap, addComponent string) {
	if strings.Contains(addComponent, "=") {
		kv := strings.Split(addComponent, "=")

		for _, validReplacementComponent := range validReplacementComponents {
			if kv[0] == validReplacementComponent {
				addressMap[kv[0]] = kv[1]
				break
			}
		}
	}
}

var countryCheck = regexp.MustCompile(`\$(\w*)`)

func determineCountry(addressMap addressMap, newCountry string) (string, error) {
	matches := countryCheck.FindStringSubmatch(newCountry)

	if matches != nil {
		component := matches[1]
		componentVal, hasComponent := addressMap[component]

		r, err := regexp.Compile(`\$` + component)

		if err != nil {
			return "", errors.New(fmt.Sprintf("Could not compile regex for component %s: %v", componentVal, err))
		}

		if hasComponent {
			newCountry = r.ReplaceAllString(newCountry, componentVal)
		} else {
			newCountry = r.ReplaceAllString(newCountry, "")
		}
	}

	return newCountry, nil
}

var sintMaartenCheck = regexp.MustCompile("(?i)sint maarten")
var arubaCheck = regexp.MustCompile("(?i)aruba")

// special conditions taken from other processors such as @fragaria/Address-formatter
func applySpecialCases(addressMap addressMap) {
	if addressMap["country_code"] == "NL" {
		if addressMap["state"] == "Curaçao" {
			addressMap["country_code"] = "CW"
			addressMap["country"] = "Curaçao"
		} else if isMatching := sintMaartenCheck.MatchString(addressMap["state"]); isMatching {
			addressMap["country_code"] = "SX"
			addressMap["country"] = "Sint Maarten"
		} else if isMatching := arubaCheck.MatchString(addressMap["state"]); isMatching {
			addressMap["country_code"] = "AW"
			addressMap["country"] = "Aruba"
		}
	}
}

func getCountryCodeAlias(countryCode string) (string, bool) {
	for countryCodeToCheck, countryCodeAlias := range commonCountryCodeAliases {
		if countryCodeToCheck == countryCode {
			return countryCodeAlias, true
		}
	}

	return "", false
}

func getCountyCode(county string, countryCode string, countyCodes map[string]map[string]interface{}) string {
	county = strings.ToUpper(county)

	if _, hasCountyCodes := countyCodes[countryCode]; hasCountyCodes {
		for code, v := range countyCodes[countryCode] {
			if countyOfCode, isString := v.(string); isString {
				if strings.ToUpper(countyOfCode) == county {
					return code
				}
			} else if variants, hasVariants := v.(map[string]interface{}); hasVariants {
				for _, countyVariant := range variants {
					if strings.ToUpper(countyVariant.(string)) == county {
						return code
					}
				}
			}
		}
	}

	return ""
}

func getStateCode(state string, countryCode string, stateCodes map[string]map[string]interface{}) string {
	state = strings.ToUpper(state)

	if _, hasStateCodes := stateCodes[countryCode]; hasStateCodes {
		for code, v := range stateCodes[countryCode] {
			if stateOfCode, isString := v.(string); isString {
				if strings.ToUpper(stateOfCode) == state {
					return code
				}
			} else if variants, hasVariants := v.(map[string]interface{}); hasVariants {
				for _, stateVariant := range variants {
					if strings.ToUpper(stateVariant.(string)) == state {
						return code
					}
				}
			}
		}
	}

	return ""
}

var urlCheck = regexp.MustCompile(`^https?://`)

func applyUrlCleanup(addressMap addressMap) {
	for k, v := range addressMap {
		if urlCheck.MatchString(v) {
			delete(addressMap, k)
		}
	}
}
