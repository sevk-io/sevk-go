package markup

import (
	"fmt"
	"regexp"
	"strings"
)

func processTag(content, tagName string, processor func(map[string]string, string) string, budget *renderBudget) string {
	result := content
	openPattern := regexp.MustCompile(fmt.Sprintf(`(?i)<%s([^>]*)>`, tagName))
	closePattern := regexp.MustCompile(fmt.Sprintf(`(?i)</%s>`, tagName))
	nestedOpenPattern := regexp.MustCompile(fmt.Sprintf(`(?i)<%s`, tagName))

	for iterations := 0; iterations < maxTagIterations; iterations++ {
		budget.chargeScanWork(len(result))
		locations := openPattern.FindAllStringIndex(result, -1)
		if len(locations) == 0 {
			break
		}

		type replacement struct {
			start int
			end   int
			text  string
		}
		replacements := make([]replacement, 0, len(locations))

		for _, location := range locations {
			start := location[0]
			if len(replacements) > 0 && start < replacements[len(replacements)-1].end {
				continue
			}

			match := openPattern.FindStringSubmatch(result[start:])
			attributes := ""
			if len(match) > 1 {
				attributes = match[1]
			}

			innerStart := location[1]
			closeLocation := closePattern.FindStringIndex(result[innerStart:])
			if closeLocation == nil {
				budget.chargeScanWork(len(result) - innerStart)
				continue
			}
			budget.chargeScanWork(closeLocation[0])

			innerEnd := innerStart + closeLocation[0]
			inner := result[innerStart:innerEnd]
			if nestedOpenPattern.MatchString(inner) {
				continue
			}

			replacements = append(replacements, replacement{
				start: start,
				end:   innerEnd + closeLocation[1] - closeLocation[0],
				text:  processor(parseAttributes(attributes), inner),
			})
		}

		if len(replacements) == 0 {
			break
		}

		var builder strings.Builder
		cursor := 0
		for _, replacement := range replacements {
			builder.WriteString(result[cursor:replacement.start])
			builder.WriteString(replacement.text)
			cursor = replacement.end
			if builder.Len() > maxRenderOutput {
				failRender(fmt.Sprintf("Email markup expands beyond the maximum size of %d bytes", maxRenderOutput))
			}
		}
		builder.WriteString(result[cursor:])
		rebuilt := builder.String()
		if rebuilt == result {
			break
		}
		result = rebuilt
	}

	return result
}

func parseAttributes(source string) map[string]string {
	attributes := make(map[string]string)
	pattern := regexp.MustCompile(`([\w-]+)=(?:"([^"]*)"|'([^']*)')`)

	for _, match := range pattern.FindAllStringSubmatch(source, -1) {
		if len(match) <= 2 {
			continue
		}
		if match[2] != "" {
			attributes[match[1]] = match[2]
		} else {
			attributes[match[1]] = match[3]
		}
	}

	return attributes
}
