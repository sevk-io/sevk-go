package markup

import (
	"regexp"
	"strings"
)

var selfClosingBlockPattern = regexp.MustCompile(`(?i)<block([^>]*?)/>`)

func processMarkup(content string) (string, string) {
	assertMarkupDepth(content)
	budget := &renderBudget{}

	result := processTag(content, "block", processBlockTag, budget)
	result = selfClosingBlockPattern.ReplaceAllStringFunc(result, func(match string) string {
		return processBlockTag(parseAttributes(selfClosingBlockPattern.FindStringSubmatch(match)[1]), "")
	})
	assertMarkupDepth(result)

	result = regexp.MustCompile(`(?i)<link(\s+(?:[\w-]+=(?:"[^"]*"|'[^']*')\s+)*)href=`).ReplaceAllString(result, "<sevk-link${1}href=")
	result = regexp.MustCompile(`(?i)</link>`).ReplaceAllString(result, "</sevk-link>")

	result = processTag(result, "section", processSection, budget)
	result = processTag(result, "column", processColumn, budget)

	rowCounter := 0
	gapRules := []string{}
	result = processTag(result, "row", func(attributes map[string]string, inner string) string {
		return processRow(attributes, inner, &rowCounter, &gapRules)
	}, budget)
	result = processTag(result, "container", processContainer, budget)

	result = processTag(result, "heading", processHeading, budget)
	result = processTag(result, "paragraph", processParagraph, budget)
	result = processTag(result, "text", processText, budget)
	result = processTag(result, "button", processButton, budget)

	imagePattern := regexp.MustCompile(`(?i)<image([^>]*?)/?>`)
	result = imagePattern.ReplaceAllStringFunc(result, func(match string) string {
		parts := imagePattern.FindStringSubmatch(match)
		return processImage(parseAttributes(parts[1]))
	})

	dividerPattern := regexp.MustCompile(`(?i)<divider([^>]*?)/?>`)
	result = dividerPattern.ReplaceAllStringFunc(result, func(match string) string {
		parts := dividerPattern.FindStringSubmatch(match)
		return processDivider(parseAttributes(parts[1]))
	})

	result = regexp.MustCompile(`(?i)</divider>`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`(?i)</image>`).ReplaceAllString(result, "")
	result = processTag(result, "sevk-link", processLink, budget)
	result = processTag(result, "list", processList, budget)
	result = processTag(result, "li", processListItem, budget)
	result = processTag(result, "codeblock", func(attributes map[string]string, inner string) string {
		return processCodeBlock(attributes, inner, budget)
	}, budget)

	result = regexp.MustCompile(`(?i)</(?:container|section|row|column|heading|paragraph|text|button|sevk-link)>`).ReplaceAllString(result, "")

	for _, pattern := range []string{
		`(?i)<sevk-email[^>]*>`, `(?i)</sevk-email>`,
		`(?i)<sevk-body[^>]*>`, `(?i)</sevk-body>`,
		`(?i)<email[^>]*>`, `(?i)</email>`,
		`(?i)<mail[^>]*>`, `(?i)</mail>`,
		`(?i)<body[^>]*>`, `(?i)</body>`,
	} {
		result = regexp.MustCompile(pattern).ReplaceAllString(result, "")
	}

	gapStyles := ""
	if len(gapRules) > 0 {
		gapStyles = "@media only screen and (max-width:479px){" +
			strings.Join(gapRules, "") +
			"}"
	}

	return strings.TrimSpace(result), gapStyles
}
