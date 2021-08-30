package address_formatter

type Address struct {
	Attention     string
	HouseNumber   string
	House         string
	Road          string
	Hamlet        string
	Village       string
	PostalCity    string
	City          string
	Municipality  string
	County        string
	CountyCode    string
	StateDistrict string
	State         string
	StateCode     string
	Postcode      string
	Suburb        string
	Region        string
	Town          string
	Island        string
	Archipelago   string
	Country       string
	CountryCode   string
	Continent     string
}

type addressMap map[string]string
