package address_formatter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddressToMap(t *testing.T) {
	address := Address{
		Attention:     "Lorem",
		HouseNumber:   "ipsum",
		House:         "dolor",
		Road:          "sit",
		Hamlet:        "amet",
		Village:       "consetetur",
		PostalCity:    "sadipscing",
		City:          "elitr",
		Municipality:  "sed",
		County:        "diam",
		CountyCode:    "nonumy",
		StateDistrict: "eirmod",
		State:         "ut",
		StateCode:     "tempor",
		Postcode:      "invidunt",
		Suburb:        "labore",
		Region:        "gubergren",
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

	addressMap, err := addressToMap(&address)

	assert.Nil(t, err)
	assert.Equal(t, addressMap, expectedAddressMap)
}
