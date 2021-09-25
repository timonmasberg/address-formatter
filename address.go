package address_formatter

type Address struct {
	Attention     string
	House         string
	HouseNumber   string
	Road          string
	Hamlet        string
	Village       string
	Neighbourhood string
	PostalCity    string
	City          string
	CityDistrict  string
	Municipality  string
	County        string
	CountyCode    string
	StateDistrict string
	Postcode      string
	State         string
	StateCode     string
	Region        string
	Suburb        string
	Town          string
	Island        string
	Archipelago   string
	Country       string
	CountryCode   string
	Continent     string
}

type addressMap map[string]string
