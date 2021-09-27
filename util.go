package address_formatter

import (
	"errors"
	"reflect"
	"strings"
)

var addressMemberNameMapping = map[string]string{
	"Attention":     "attention",
	"HouseNumber":   "house_number",
	"House":         "house",
	"Road":          "road",
	"Village":       "village",
	"Suburb":        "suburb",
	"City":          "city",
	"County":        "county",
	"CountyCode":    "county_code",
	"Postcode":      "postcode",
	"StateDistrict": "state_district",
	"State":         "state",
	"StateCode":     "state_code",
	"Region":        "region",
	"Island":        "island",
	"Country":       "country",
	"CountryCode":   "country_code",
	"Continent":     "continent",
	"Archipelago":   "archipelago",
	"Municipality":  "municipality",
	"Town":          "town",
	"Hamlet":        "hamlet",
	"PostalCity":    "postal_city",
	"Neighbourhood": "neighbourhood",
	"CityDistrict":  "city_district",
} // (struct Name, template Name)

// convert Address to a map of names used in OpenCageData templates and their value
func addressToMap(address *Address) (addressMap, error) {
	v := reflect.ValueOf(*address)
	addressMap := make(map[string]string, v.NumField())
	addressType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := addressType.Field(i)

		mapFieldName, hasMapping := addressMemberNameMapping[fi.Name]

		if hasMapping {
			addressMap[mapFieldName] = v.Field(i).String()
		} else {
			return nil, errors.New(fi.Name + " has no corresponding Name")
		}
	}

	return addressMap, nil
}

// MapToAddress Convert map of address components used in OpenCageData templates and their aliases into an Address struct
func MapToAddress(addressMap addressMap, componentAliases map[string]string, unknownAsAttention bool) *Address {
	// replace common aliases with their main keys used in templates
	addressMap = applyComponentAliases(addressMap, componentAliases)

	// invert addressMemberNameMapping to map component names to Address struct fields
	componentNameAddressFieldMapping := getNameAddressFieldMapping()

	var address Address
	av := reflect.ValueOf(&address).Elem()

	unknownFieldValues := make([]string, 0)

	for k, v := range addressMap {
		name, hasCorrespondingField := componentNameAddressFieldMapping[k]

		if hasCorrespondingField {
			av.FieldByName(name).Set(reflect.ValueOf(v))
		} else // has no corresponding field and is also not an alias => attention
		if _, hasAlias := componentAliases[k]; unknownAsAttention && !hasAlias {
			unknownFieldValues = append(unknownFieldValues, v)
		}
	}

	if attention, hasAttention := addressMap["attention"]; hasAttention {
		address.Attention = attention
	} else {
		address.Attention = strings.Join(unknownFieldValues, ", ")
	}

	return &address
}

// this duplicates values from the alias to the given component name mapping
func applyComponentAliases(addressMap addressMap, componentAliases map[string]string) addressMap {
	for k, v := range addressMap {
		if alias, hasAlias := componentAliases[k]; hasAlias {
			if _, aliasAlreadyGiven := addressMap[alias]; !aliasAlreadyGiven || addressMap[alias] == "" {
				addressMap[alias] = v
			}
		}
	}

	return addressMap
}

func getNameAddressFieldMapping() map[string]string {
	componentNameAddressFieldMapping := make(map[string]string, len(addressMemberNameMapping))
	for k, v := range addressMemberNameMapping {
		componentNameAddressFieldMapping[v] = k
	}

	return componentNameAddressFieldMapping
}

func findTemplate(countryCode string, templates map[string]interface{}) template {
	template, hasTemplate := templates[countryCode]

	if hasTemplate {
		return template
	}

	return templates["default"]
}
