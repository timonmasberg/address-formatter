package address_formatter

import (
	"errors"
	"reflect"
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
} // (struct name, template name)

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
			return nil, errors.New(fi.Name + " has no corresponding name")
		}
	}

	return addressMap, nil
}
