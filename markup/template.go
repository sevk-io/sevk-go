package markup

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

const (
	maxIfPasses           = 2000
	maxEachItems          = 1000
	maxEachOutput         = 2_000_000
	maxTemplateWork       = 20_000_000
	outputExceededMessage = "Email block exceeds the maximum template output size"
)

var (
	ifOpenPattern       = regexp.MustCompile(`\{%#if\s+([^%]+)%\}`)
	eachOpenPattern     = regexp.MustCompile(`\{%#each\s+(\w+)(?:\s+as\s+(\w+))?%\}`)
	fallbackPattern     = regexp.MustCompile(`\{%(\w+)\s*\?\?([^%]+)%\}`)
	valuePattern        = regexp.MustCompile(`\{%(\w+)%\}`)
	equalsPattern       = regexp.MustCompile(`^(\w+)\s*==\s*"([^"]*)"$`)
	notEqualsPattern    = regexp.MustCompile(`^(\w+)\s*!=\s*"([^"]*)"$`)
	configEntityPattern = regexp.MustCompile(`(?i)&(?:amp|quot|apos|lt|gt|#39|#x27|#34|#x22);`)
)

type outputMeter struct {
	rendered int
}

func (meter *outputMeter) charge(value string) string {
	if meter.rendered+len(value) > maxEachOutput {
		failRender(outputExceededMessage)
	}
	meter.rendered += len(value)
	return value
}

func isTruthy(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed != ""
	case float64:
		return typed != 0
	case int:
		return typed != 0
	case []any:
		return len(typed) > 0
	default:
		return true
	}
}

func templateText(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case bool:
		if typed {
			return "true"
		}
		return "false"
	case string:
		return typed
	case float64:
		return jsNumber(typed)
	case int:
		return strconv.Itoa(typed)
	case map[string]any:
		return "[object Object]"
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			if _, isObject := item.(map[string]any); isObject {
				parts = append(parts, "[object Object]")
			} else {
				parts = append(parts, templateText(item))
			}
		}
		return strings.Join(parts, ",")
	default:
		return fmt.Sprintf("%v", typed)
	}
}

func evaluateCondition(expression string, config map[string]any) bool {
	expression = strings.TrimSpace(expression)

	if strings.Contains(expression, "||") {
		for _, part := range strings.Split(expression, "||") {
			if evaluateCondition(part, config) {
				return true
			}
		}
		return false
	}

	if strings.Contains(expression, "&&") {
		for _, part := range strings.Split(expression, "&&") {
			if !evaluateCondition(part, config) {
				return false
			}
		}
		return true
	}

	if match := equalsPattern.FindStringSubmatch(expression); match != nil {
		return templateText(config[match[1]]) == match[2]
	}
	if match := notEqualsPattern.FindStringSubmatch(expression); match != nil {
		return templateText(config[match[1]]) != match[2]
	}

	return isTruthy(config[expression])
}

func renderTemplate(template string, config map[string]any) string {
	meter := &outputMeter{}
	result := renderConditionals(template, config)
	result = renderLoops(result, config, meter)

	result = fallbackPattern.ReplaceAllStringFunc(result, func(match string) string {
		parts := fallbackPattern.FindStringSubmatch(match)
		if value, exists := config[parts[1]]; exists && value != nil && value != "" {
			return meter.charge(templateText(value))
		}
		return meter.charge(strings.TrimSpace(parts[2]))
	})

	return valuePattern.ReplaceAllStringFunc(result, func(match string) string {
		return meter.charge(templateText(config[valuePattern.FindStringSubmatch(match)[1]]))
	})
}

func renderConditionals(template string, config map[string]any) string {
	result := template
	work := 0

	for pass := 0; pass < maxIfPasses; pass++ {
		work += len(result)
		if work > maxTemplateWork {
			failRender(tooComplexMessage)
		}

		openings := ifOpenPattern.FindAllStringSubmatchIndex(result, -1)
		if len(openings) == 0 {
			break
		}

		replaced := false
		searchFrom := 0
		for {
			closeIndex := strings.Index(result[searchFrom:], "{%/if%}")
			if closeIndex < 0 {
				break
			}
			closeIndex += searchFrom

			var opening []int
			for _, candidate := range openings {
				if candidate[1] > closeIndex {
					break
				}
				opening = candidate
			}
			if opening == nil {
				searchFrom = closeIndex + len("{%/if%}")
				continue
			}

			condition := result[opening[2]:opening[3]]
			body := result[opening[1]:closeIndex]
			replacement := ""
			if elseIndex := strings.Index(body, "{%else%}"); elseIndex >= 0 {
				if evaluateCondition(condition, config) {
					replacement = body[:elseIndex]
				} else {
					replacement = body[elseIndex+len("{%else%}"):]
				}
			} else if evaluateCondition(condition, config) {
				replacement = body
			}

			result = result[:opening[0]] + replacement + result[closeIndex+len("{%/if%}"):]
			replaced = true
			break
		}

		if !replaced {
			break
		}
	}

	return result
}

func renderLoops(template string, config map[string]any, meter *outputMeter) string {
	if !strings.Contains(template, "{%/each%}") {
		return template
	}

	var output strings.Builder
	from := 0
	work := 0

	for from <= len(template) {
		opening := eachOpenPattern.FindStringSubmatchIndex(template[from:])
		if opening == nil {
			output.WriteString(template[from:])
			break
		}

		bodyStart := from + opening[1]
		closeIndex := strings.Index(template[bodyStart:], "{%/each%}")
		if closeIndex < 0 {
			output.WriteString(template[from:])
			break
		}
		closeIndex += bodyStart

		output.WriteString(template[from : from+opening[0]])
		items, isList := config[template[from+opening[2]:from+opening[3]]].([]any)

		if isList && len(items) > 0 {
			if len(items) > maxEachItems {
				failRender(fmt.Sprintf("Email block exceeds the maximum of %d loop items", maxEachItems))
			}

			alias := "this"
			if opening[4] >= 0 {
				alias = template[from+opening[4] : from+opening[5]]
			}
			propertyPattern := regexp.MustCompile(`\{%` + regexp.QuoteMeta(alias) + `\.(\w+)%\}`)
			body := template[bodyStart:closeIndex]

			for _, item := range items {
				work += len(body) + 16
				if work > maxTemplateWork {
					failRender(tooComplexMessage)
				}
				output.WriteString(renderLoopItem(body, item, config, propertyPattern, meter))
			}
		}

		from = closeIndex + len("{%/each%}")
	}

	return output.String()
}

func renderLoopItem(body string, item any, config map[string]any, propertyPattern *regexp.Regexp, meter *outputMeter) string {
	fields, _ := item.(map[string]any)
	itemOut := 0
	charge := func(value string) string {
		if meter.rendered+itemOut+len(value) > maxEachOutput {
			failRender(outputExceededMessage)
		}
		itemOut += len(value)
		return value
	}

	rendered := propertyPattern.ReplaceAllStringFunc(body, func(match string) string {
		return charge(templateText(itemProperty(item, fields, propertyPattern.FindStringSubmatch(match)[1])))
	})
	rendered = valuePattern.ReplaceAllStringFunc(rendered, func(match string) string {
		return charge(templateText(config[valuePattern.FindStringSubmatch(match)[1]]))
	})

	meter.rendered += len(rendered)
	if meter.rendered > maxEachOutput {
		failRender(outputExceededMessage)
	}
	return rendered
}

func itemProperty(item any, fields map[string]any, property string) any {
	if fields != nil {
		return fields[property]
	}
	if property != "length" {
		return nil
	}
	switch typed := item.(type) {
	case string:
		return len(utf16.Encode([]rune(typed)))
	case []any:
		return len(typed)
	}
	return nil
}

func processBlockTag(attributes map[string]string, inner string) string {
	template := strings.TrimSpace(inner)
	if template == "" {
		template = attributes["template"]
	}
	if template == "" {
		return ""
	}

	configJSON := attributes["config"]
	if configJSON == "" {
		configJSON = "{}"
	}
	configJSON = configEntityPattern.ReplaceAllStringFunc(configJSON, func(entity string) string {
		switch strings.ToLower(entity) {
		case "&amp;":
			return "&"
		case "&quot;", "&#34;", "&#x22;":
			return `"`
		case "&apos;", "&#39;", "&#x27;":
			return "'"
		case "&lt;":
			return "<"
		case "&gt;":
			return ">"
		default:
			return entity
		}
	})

	var config map[string]any
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		if err := json.Unmarshal([]byte(strings.ReplaceAll(configJSON, "'", `"`)), &config); err != nil {
			config = map[string]any{}
		}
	}
	if config == nil {
		config = map[string]any{}
	}

	return renderTemplate(template, config)
}
