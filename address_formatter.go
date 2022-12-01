package addrFmt

import (
	"errors"
	"github.com/cbroglie/mustache"
	"html"
	"log"
	"regexp"
	"strings"
)

var commonCountryCodeAliases = map[string]string{"UK": "GB"}
var validReplacementComponents = []string{"state"}
var requiredAddressProperties = []string{"road", "postcode"}

type replacement struct {
	pattern *regexp.Regexp
	replace string
}

var replacements = []replacement{
	{pattern: regexp.MustCompile(`[\},\s]+$`), replace: ""},
	{pattern: regexp.MustCompile(`(?m)^[,\s]+`), replace: ""},
	{pattern: regexp.MustCompile(`(?m)^- `), replace: ""},                      // line starting with dash due to a parameter missing
	{pattern: regexp.MustCompile(`,\s*,`), replace: ", "},                      // multiple commas to one
	{pattern: regexp.MustCompile(`[[:blank:]]+,[[:blank:]]+/`), replace: ", "}, // one horiz whitespace behind comma
	{pattern: regexp.MustCompile(`[[:blank:]][[:blank:]]+`), replace: " "},     // multiple horiz whitespace to one
	{pattern: regexp.MustCompile(`[[:blank:]]\n`), replace: "\n"},              // horiz whitespace, newline to newline
	{pattern: regexp.MustCompile(`\n,`), replace: "\n"},                        // newline comma to just newline
	{pattern: regexp.MustCompile(`,,+`), replace: ","},                         // multiple commas to one
	{pattern: regexp.MustCompile(`,\n`), replace: "\n"},                        // comma newline to just newline
	{pattern: regexp.MustCompile(`\n[[:blank:]]+`), replace: "\n"},             // newline plus space to newline
	{pattern: regexp.MustCompile(`\n\n+`), replace: "\n"},                      // multiple newline to one
}

// FormatAddress formats an Address object based on it
func FormatAddress(address *Address, config *Config) (interface{}, error) {
	// ease up the Address into a map to make it accessible via index
	addressMap, err := addressToMap(address)

	if err != nil {
		return nil, err
	}

	template := findTemplate(address.CountryCode, config.Templates)
	render, err := applyTemplate(addressMap, template, config.Templates)
	if err != nil {
		return nil, err
	}

	return getOutput(render, config.OutputFormat)
}

func applyTemplate(addressMap addressMap, template template, templates map[string]template) (string, error) {
	templateText := chooseTemplateText(addressMap, template, templates)

	render, _ := mustache.Render(templateText, getRenderInput(addressMap))
	// unescape render to enforce official mustache HTML escaping rules
	render = html.UnescapeString(render)
	// todo: postformat replacements rely on a clean render but can mess it up again... (constraint by OpenCageData)
	var err error
	render, err = cleanupRender(render)
	if err != nil {
		return "", err
	}
	render = applyPostformatReplacements(render, template)
	render, err = cleanupRender(render)
	if err != nil {
		return "", err
	}

	return render, nil
}

var possibilitiesRegExp = regexp.MustCompile(`\s*\|\|\s*`)

func getRenderInput(addressMap addressMap) map[string]interface{} {
	input := make(map[string]interface{})

	for k, v := range addressMap {
		input[k] = v
	}

	input["first"] = func(t string, f func(string) (string, error)) (string, error) {
		t, _ = f(t)
		possibilities := possibilitiesRegExp.Split(t, -1)

		for _, possibility := range possibilities {
			if possibility != "" {
				return possibility, nil
			}
		}

		return "", nil
	}

	return input
}

func chooseTemplateText(address addressMap, template template, templates map[string]template) string {
	var templateText string

	templateValue, isTemplateMap := template.(map[string]interface{})

	if isTemplateMap {
		missingPropertyCount := 0
		for _, requiredProperty := range requiredAddressProperties {
			if _, hasProperty := address[requiredProperty]; !hasProperty {
				missingPropertyCount++
			}
		}

		if missingPropertyCount == len(requiredAddressProperties) {
			if fallbackTemplate, hasFallbackTemplate := templateValue["fallback_template"]; hasFallbackTemplate {
				templateText = fallbackTemplate.(string)
			} else if defaultTemplate, hasDefaultTemplate := templates["default"]; hasDefaultTemplate {
				templateText = defaultTemplate.(map[string]interface{})["fallback_template"].(string)
			}
		} else // has country specific template
		if addressTemplate, hasAddressTemplate := templateValue["address_template"]; hasAddressTemplate {
			templateText = addressTemplate.(string)
		} else // has default template
		if defaultTemplate, hasDefaultTemplate := templates["default"]; hasDefaultTemplate {
			templateText = defaultTemplate.(map[string]interface{})["address_template"].(string)
		}
	} else {
		templateText = template.(string)
	}

	return templateText
}

func cleanupRender(render string) (string, error) {
	for _, replacement := range replacements {
		render = replacement.pattern.ReplaceAllString(render, replacement.replace)

		render = dedupe(strings.Split(render, "\n"), "\n", func(s string) string {
			return dedupe(strings.Split(s, ", "), ", ", func(s string) string {
				return s
			})
		})
	}

	return strings.TrimSpace(render), nil
}

func dedupe(chunks []string, glue string, modifier func(s string) string) string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if strings.ToLower(chunk) == "new york" {
			seen[chunk] = true
			result = append(result, chunk)
		} else if seenChunk, hasChunk := seen[chunk]; !hasChunk || !seenChunk {
			seen[chunk] = true
			result = append(result, modifier(chunk))
		}
	}

	return strings.Join(result, glue)
}

func applyPostformatReplacements(render string, template template) string {
	if postformatReplacements, hasReplacements := template.(map[string]interface{})["postformat_replace"].([]interface{}); hasReplacements {
		for _, replacement := range postformatReplacements {
			r, err := regexp.Compile(replacement.([]interface{})[0].(string))
			if err != nil {
				log.Printf("Could not replace due to bad regexp: %v", err)
			}

			render = r.ReplaceAllString(render, replacement.([]interface{})[1].(string))
		}
	}

	return render
}

func getOutput(render string, outputFormat OutputFormat) (interface{}, error) {
	switch outputFormat {
	case Array:
		return strings.Split(render, "\n"), nil
	case OneLine:
		return strings.Replace(render, "\n", ", ", -1), nil
	case PostalFormat:
		return render + "\n", nil
	default:
		return nil, errors.New("invalid output format")
	}
}
