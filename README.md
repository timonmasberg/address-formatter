# Golang Address Formatter

This Address Formatter package is able to convert an address into many international formats of postal addresses based on [OpenCage Configuration](https://github.com/OpenCageData) (or custom ones). 
You can use the provided address structure or use a map of address components.

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=timonmasberg_address-formatter&metric=alert_status)](https://sonarcloud.io/dashboard?id=timonmasberg_address-formatter)
[![Go Report Card](https://goreportcard.com/badge/github.com/timonmasberg/address-formatter)](https://goreportcard.com/report/github.com/timonmasberg/address-formatter)
[![codecov](https://codecov.io/gh/timonmasberg/address-formatter/branch/master/graph/badge.svg?token=M18CWXGVL7)](https://codecov.io/gh/timonmasberg/address-formatter)

## Usage
Import 
```go
import (
    "fmt"
    "github.com/timonmasberg/address-formatter"
)
```
Set up your config files or use [OpenCage Configuration](https://github.com/OpenCageData) (`/conf`)<br>
Load config from config files: 
```go
config := addrFmt.LoadConfig(addrFmt.ConfigFiles{
            CountriesPath:     "templates/countries/worldwide.yaml",
            ComponentsPath:    "templates/components.yaml",
            StateCodesPath:    "templates/state_codes.yaml",
            CountryToLangPath: "templates/country2lang.yaml",
            CountyCodesPath:   "templates/county_codes.yaml",
            CountryCodesPath:  "templates/country_codes.yaml",
            AbbreviationFiles: "templates/abbreviations/*.yaml",
        })
```
You can choose between 3 output formats:
1. Array (returns a slice where each entry represents an address component)
2. OneLine (address components joined with a comma in one line)
3. PostalFormat (valid postal format for letters etc with a trailing \n)
```go
config.OutputFormat = addrFmt.PostalFormat

address :=  &addrFmt.Address{
    House:         "Bundestag",
    HouseNumber:   "1",
    Road:          "Platz der Republik",
    City:          "Berlin",
    Postcode:      "11011",
    State:         "Berlin",
    Country:       "Deutschland",
}

formattedAddress, err := addrFmt.FormatAddress(address, config)
if err != nil {
    fmt.Printf("Failed to format address: %v", err)
} else {
    fmt.Println(formattedAddress)
    /*
        Bundestag
        Platz der Republik 1
        11011 Berlin
        Deutschland
    */
}
```

If you have data from sources such as OSM you probably have a map of unknown data. This package can cleanup the map by using data from the configurations and turn it into an Address structure. You can also use `MapToAddress` directly if you are aware of the quality. Just make sure you add your component names to the component aliases as this package uses the names from OpenCageData.

```go
    	// with unknown data such as from osm
addressMap := make(map[string]string)
addressMap["house"] = "Bundestag"
addressMap["house_number"] = "1"
addressMap["road"] = "Platz der Republik"
addressMap["postcode"] = "11011"
addressMap["state"] = "Berlin"
addressMap["country"] = "Deutschland"

address, err = addrFmt.GetFixedAddress(addressMap, config)
if err != nil {
    fmt.Printf("Fixing address failed: %v", err)
}

```
If you want to treat every unknown component name as an attention entry, set `UnknownAsAttention` to true

```go
config.UnknownAsAttention = true

addressMap["mutti"] = "Angela Merkel"
addressMap["vati"] = "Frank-Walter Steinmeier"

address, err = addrFmt.GetFixedAddress(addressMap, config)
if err != nil {
    fmt.Printf("Fixing address failed: %v", err)
} else {
    fmt.Println(address.Attention)
    // Angela Merkel, Frank-Walter Steinmeier
}
```

<i>Abbreviations tbd.</i>

## Testing
Load the config files from the submodule with `copy-templates.cmd`.
Testing the formatter relies on testcase files. 
You can execute `copy-testcases.cmd` to use the OpenCageData testcases, or you can run the tests with your own. 
Just create them in the testcases folder with the same structure.

## License
[MIT](https://choosealicense.com/licenses/mit/)
