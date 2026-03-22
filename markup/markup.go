// Package markup provides email markup rendering for Sevk
package markup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// FontConfig represents a font configuration
type FontConfig struct {
	ID   string
	Name string
	URL  string
}

// EmailHeadSettings represents head settings for email generation
type EmailHeadSettings struct {
	Title       string
	PreviewText string
	Styles      string
	Fonts       []FontConfig
	Lang        string
	Dir         string
}

// ParsedEmailContent represents parsed email content
type ParsedEmailContent struct {
	Body         string
	HeadSettings EmailHeadSettings
}

// Render converts Sevk markup to email-compatible HTML
func Render(markupContent string) string {
	return GenerateEmailFromMarkup(markupContent, nil)
}

// GenerateEmailFromMarkup generates email HTML from Sevk markup
func GenerateEmailFromMarkup(htmlContent string, headSettings *EmailHeadSettings) string {
	// Always parse to extract clean body content (strips <mail>/<head> wrapper tags)
	parsed := ParseEmailHTML(htmlContent)
	var settings EmailHeadSettings
	if headSettings != nil {
		settings = *headSettings
	} else {
		settings = parsed.HeadSettings
	}
	contentToProcess := parsed.Body

	normalized := normalizeMarkup(contentToProcess)
	processed := processMarkup(normalized)

	// Resolve lang and dir with defaults
	lang := settings.Lang
	if lang == "" {
		lang = "en"
	}
	dir := settings.Dir
	if dir == "" {
		dir = "ltr"
	}

	// Build head content
	titleTag := ""
	if settings.Title != "" {
		titleTag = fmt.Sprintf("<title>%s</title>", settings.Title)
	}

	fontLinks := generateFontLinks(settings.Fonts)

	customStyles := ""
	if settings.Styles != "" {
		customStyles = fmt.Sprintf(`<style type="text/css">%s</style>`, settings.Styles)
	}

	previewText := ""
	if settings.PreviewText != "" {
		previewText = fmt.Sprintf(`<div style="display:none;font-size:1px;color:#ffffff;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;">%s</div>`, settings.PreviewText)
	}

	return fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html lang="%s" dir="%s" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
<head>
<meta content="text/html; charset=UTF-8" http-equiv="Content-Type"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<meta name="x-apple-disable-message-reformatting"/>
<meta content="IE=edge" http-equiv="X-UA-Compatible"/>
<meta name="format-detection" content="telephone=no,address=no,email=no,date=no,url=no"/>
<!--[if mso]>
<noscript>
<xml>
<o:OfficeDocumentSettings>
<o:AllowPNG/>
<o:PixelsPerInch>96</o:PixelsPerInch>
</o:OfficeDocumentSettings>
</xml>
</noscript>
<![endif]-->
<style type="text/css">
#outlook a { padding: 0; }
body { margin: 0; padding: 0; -webkit-text-size-adjust: 100%%; -ms-text-size-adjust: 100%%; }
table, td { border-collapse: collapse; mso-table-lspace: 0pt; mso-table-rspace: 0pt; }
.sevk-row-table { border-collapse: separate !important; }
img { border: 0; height: auto; line-height: 100%%; outline: none; text-decoration: none; -ms-interpolation-mode: bicubic; }
@media only screen and (max-width: 479px) {
  .sevk-row-table { width: 100%% !important; }
  .sevk-column { display: block !important; width: 100%% !important; max-width: 100%% !important; }
}
</style>
%s
%s
%s
</head>
<body style="margin:0;padding:0;word-spacing:normal;-webkit-text-size-adjust:100%%;-ms-text-size-adjust:100%%;font-family:ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif">
<div aria-roledescription="email" role="article">
%s
%s
</div>
</body>
</html>`, lang, dir, titleTag, fontLinks, customStyles, previewText, processed)
}

func normalizeMarkup(content string) string {
	result := content

	// Replace <link> with <sevk-link>
	if strings.Contains(result, "<link") {
		re := regexp.MustCompile(`<link\s+href=`)
		result = re.ReplaceAllString(result, "<sevk-link href=")
		result = strings.ReplaceAll(result, "</link>", "</sevk-link>")
	}

	if !strings.Contains(result, "<sevk-email") && !strings.Contains(result, "<email") && !strings.Contains(result, "<mail") {
		result = fmt.Sprintf("<mail><body>%s</body></mail>", result)
	}

	return result
}

func generateFontLinks(fonts []FontConfig) string {
	var links []string
	for _, f := range fonts {
		links = append(links, fmt.Sprintf(`<link href="%s" rel="stylesheet" type="text/css" />`, f.URL))
	}
	return strings.Join(links, "\n")
}

// ParseEmailHTML parses email HTML and extracts head settings
func ParseEmailHTML(content string) ParsedEmailContent {
	if strings.Contains(content, "<email>") || strings.Contains(content, "<email ") ||
		strings.Contains(content, "<mail>") || strings.Contains(content, "<mail ") {
		return parseSevkMarkup(content)
	}
	return ParsedEmailContent{
		Body:         content,
		HeadSettings: EmailHeadSettings{},
	}
}

func parseSevkMarkup(content string) ParsedEmailContent {
	var headSettings EmailHeadSettings

	// Parse lang and dir from <mail> or <email> root tag
	rootRe := regexp.MustCompile(`(?i)<(?:email|mail)([^>]*)>`)
	if rootMatch := rootRe.FindStringSubmatch(content); len(rootMatch) > 1 {
		rootAttrs := rootMatch[1]
		langRe := regexp.MustCompile(`(?i)lang=["']([^"']+)["']`)
		if langMatch := langRe.FindStringSubmatch(rootAttrs); len(langMatch) > 1 {
			headSettings.Lang = langMatch[1]
		}
		dirRe := regexp.MustCompile(`(?i)dir=["']([^"']+)["']`)
		if dirMatch := dirRe.FindStringSubmatch(rootAttrs); len(dirMatch) > 1 {
			headSettings.Dir = dirMatch[1]
		}
	}

	// Extract title
	titleRe := regexp.MustCompile(`<title[^>]*>([\s\S]*?)</title>`)
	if matches := titleRe.FindStringSubmatch(content); len(matches) > 1 {
		headSettings.Title = strings.TrimSpace(matches[1])
	}

	// Extract preview
	previewRe := regexp.MustCompile(`<preview[^>]*>([\s\S]*?)</preview>`)
	if matches := previewRe.FindStringSubmatch(content); len(matches) > 1 {
		headSettings.PreviewText = strings.TrimSpace(matches[1])
	}

	// Extract styles
	styleRe := regexp.MustCompile(`<style[^>]*>([\s\S]*?)</style>`)
	if matches := styleRe.FindStringSubmatch(content); len(matches) > 1 {
		headSettings.Styles = strings.TrimSpace(matches[1])
	}

	// Extract fonts
	fontRe := regexp.MustCompile(`<font[^>]*name=["']([^"']*)["'][^>]*url=["']([^"']*)["'][^>]*/?\s*>`)
	fontMatches := fontRe.FindAllStringSubmatch(content, -1)
	for i, match := range fontMatches {
		if len(match) > 2 {
			headSettings.Fonts = append(headSettings.Fonts, FontConfig{
				ID:   fmt.Sprintf("font-%d", i),
				Name: match[1],
				URL:  match[2],
			})
		}
	}

	// Extract body
	bodyRe := regexp.MustCompile(`<body[^>]*>([\s\S]*?)</body>`)
	var body string
	if matches := bodyRe.FindStringSubmatch(content); len(matches) > 1 {
		body = strings.TrimSpace(matches[1])
	} else {
		body = content
		patterns := []string{
			`<email[^>]*>`, `</email>`,
			`<mail[^>]*>`, `</mail>`,
			`<head[^>]*>[\s\S]*?</head>`,
			`<title[^>]*>[\s\S]*?</title>`,
			`<preview[^>]*>[\s\S]*?</preview>`,
			`<style[^>]*>[\s\S]*?</style>`,
			`<font[^>]*>[\s\S]*?</font>`,
			`<font[^>]*/?>`,
		}
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			body = re.ReplaceAllString(body, "")
		}
		body = strings.TrimSpace(body)
	}

	return ParsedEmailContent{Body: body, HeadSettings: headSettings}
}

func processMarkup(content string) string {
	result := content

	// Convert <link> to <sevk-link> (before processing tags)
	if strings.Contains(result, "<link") {
		linkRe := regexp.MustCompile(`<link\s+href=`)
		result = linkRe.ReplaceAllString(result, "<sevk-link href=")
		result = strings.ReplaceAll(result, "</link>", "</sevk-link>")
	}

	// Process block tags BEFORE other tags
	result = processTag(result, "block", func(attrs map[string]string, inner string) string {
		return processBlockTag(attrs, inner)
	})

	// Process section tags
	result = processTag(result, "section", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody>
<tr>
<td>%s</td>
</tr>
</tbody>
</table>`, styleStr, inner)
	})

	// Process row tags
	rowCounter := 0
	currentRowGap := 0
	result = processTag(result, "row", func(attrs map[string]string, inner string) string {
		gap := attrs["gap"]
		if gap == "" {
			gap = "0"
		}
		style := extractAllStyleAttributes(attrs)
		delete(style, "gap")
		styleStr := styleToString(style)
		gapPx := strings.ReplaceAll(gap, "px", "")
		gapNum, _ := strconv.Atoi(gapPx)
		rowID := fmt.Sprintf("sevk-row-%d", rowCounter)
		rowCounter++
		currentRowGap = gapNum

		// Assign equal widths to columns if more than one
		columnCountRe := regexp.MustCompile(`class="sevk-column"`)
		columnCount := len(columnCountRe.FindAllString(inner, -1))
		processedInner := inner
		if columnCount > 1 {
			equalWidth := fmt.Sprintf("%d%%", 100/columnCount)
			colWidthRe := regexp.MustCompile(`<td class="sevk-column" style="([^"]*)"`)
			processedInner = colWidthRe.ReplaceAllStringFunc(processedInner, func(match string) string {
				if strings.Contains(match, "width:") {
					return match
				}
				submatches := colWidthRe.FindStringSubmatch(match)
				if len(submatches) > 1 {
					return fmt.Sprintf(`<td class="sevk-column" style="width:%s;%s"`, equalWidth, submatches[1])
				}
				return match
			})
		}

		gapStyle := ""
		if gapNum > 0 {
			gapStyle = fmt.Sprintf(`<style>@media only screen and (max-width:479px){.%s > tbody > tr > td{margin-bottom:%spx !important;padding-left:0 !important;padding-right:0 !important;}.%s > tbody > tr > td:last-child{margin-bottom:0 !important;}}</style>`, rowID, gapPx, rowID)
		}
		return fmt.Sprintf(`%s<table class="sevk-row-table %s" align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody style="width:100%%">
<tr style="width:100%%">%s</tr>
</tbody>
</table>`, gapStyle, rowID, styleStr, processedInner)
	})

	// Process column tags - apply half-gap padding from parent row
	result = processTag(result, "column", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		if _, ok := style["vertical-align"]; !ok {
			style["vertical-align"] = "top"
		}
		if currentRowGap > 0 {
			halfGap := float64(currentRowGap) / 2.0
			if _, ok := style["padding-left"]; !ok {
				style["padding-left"] = fmt.Sprintf("%gpx", halfGap)
			}
			if _, ok := style["padding-right"]; !ok {
				style["padding-right"] = fmt.Sprintf("%gpx", halfGap)
			}
		}
		styleStr := styleToString(style)
		return fmt.Sprintf(`<td class="sevk-column" style="%s">%s</td>`, styleStr, inner)
	})

	// Process container tags
	result = processTag(result, "container", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		tdStyle := make(map[string]string)
		tableStyle := make(map[string]string)

		// Visual styles go on <td>, layout styles stay on <table>
		visualKeys := map[string]bool{
			"background-color": true, "background-image": true, "background-size": true,
			"background-position": true, "background-repeat": true,
			"border": true, "border-top": true, "border-right": true, "border-bottom": true, "border-left": true,
			"border-color": true, "border-width": true, "border-style": true,
			"border-radius": true, "border-top-left-radius": true, "border-top-right-radius": true,
			"border-bottom-left-radius": true, "border-bottom-right-radius": true,
			"padding": true, "padding-top": true, "padding-right": true, "padding-bottom": true, "padding-left": true,
			"color": true, "font-size": true, "font-family": true, "font-weight": true,
			"text-align": true, "text-decoration": true, "line-height": true,
		}

		for key, value := range style {
			if visualKeys[key] {
				tdStyle[key] = value
			} else {
				tableStyle[key] = value
			}
		}

		// Add border-collapse: separate when border-radius is used
		hasBorderRadius := tdStyle["border-radius"] != "" || tdStyle["border-top-left-radius"] != "" ||
			tdStyle["border-top-right-radius"] != "" || tdStyle["border-bottom-left-radius"] != "" ||
			tdStyle["border-bottom-right-radius"] != ""
		if hasBorderRadius {
			tableStyle["border-collapse"] = "separate"
		}

		// Make fixed widths responsive: width becomes max-width, width set to 100%
		if w, ok := tableStyle["width"]; ok && w != "100%" && w != "auto" {
			if _, hasMax := tableStyle["max-width"]; !hasMax {
				tableStyle["max-width"] = w
			}
			tableStyle["width"] = "100%"
		}

		tableStyleStr := styleToString(tableStyle)
		tdStyleStr := styleToString(tdStyle)
		return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody>
<tr style="width:100%%">
<td style="%s">%s</td>
</tr>
</tbody>
</table>`, tableStyleStr, tdStyleStr, inner)
	})

	// Process heading tags
	result = processTag(result, "heading", func(attrs map[string]string, inner string) string {
		level := attrs["level"]
		if level == "" {
			level = "1"
		}
		style := extractAllStyleAttributes(attrs)
		if _, ok := style["margin"]; !ok {
			style["margin"] = "0"
		}
		styleStr := styleToString(style)
		return fmt.Sprintf(`<h%s style="%s">%s</h%s>`, level, styleStr, inner, level)
	})

	// Process paragraph tags
	result = processTag(result, "paragraph", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		if _, ok := style["margin"]; !ok {
			style["margin"] = "0"
		}
		styleStr := styleToString(style)
		return fmt.Sprintf(`<p style="%s">%s</p>`, styleStr, inner)
	})

	// Process text tags
	result = processTag(result, "text", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<span style="%s">%s</span>`, styleStr, inner)
	})

	// Process button tags with MSO compatibility
	result = processTag(result, "button", func(attrs map[string]string, inner string) string {
		return processButton(attrs, inner)
	})

	// Process image tags
	result = processTag(result, "image", func(attrs map[string]string, _ string) string {
		src := attrs["src"]
		alt := attrs["alt"]
		width := attrs["width"]
		height := attrs["height"]

		style := extractAllStyleAttributes(attrs)
		// Add default image styles
		if _, ok := style["vertical-align"]; !ok {
			style["vertical-align"] = "middle"
		}
		if _, ok := style["max-width"]; !ok {
			style["max-width"] = "100%"
		}
		if _, ok := style["outline"]; !ok {
			style["outline"] = "none"
		}
		if _, ok := style["border"]; !ok {
			style["border"] = "none"
		}
		if _, ok := style["text-decoration"]; !ok {
			style["text-decoration"] = "none"
		}

		styleStr := styleToString(style)
		widthAttr := ""
		if width != "" {
			widthAttr = fmt.Sprintf(` width="%s"`, strings.TrimSuffix(width, "px"))
		}
		heightAttr := ""
		if height != "" {
			heightAttr = fmt.Sprintf(` height="%s"`, strings.TrimSuffix(height, "px"))
		}

		return fmt.Sprintf(`<img src="%s" alt="%s"%s%s style="%s" />`, src, alt, widthAttr, heightAttr, styleStr)
	})

	// Process divider tags
	result = processTag(result, "divider", func(attrs map[string]string, _ string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		classAttr := ""
		if class, ok := attrs["class"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, class)
		} else if className, ok := attrs["className"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, className)
		}
		return fmt.Sprintf(`<hr style="%s"%s />`, styleStr, classAttr)
	})

	// Clean up stray </divider> closing tags (divider is self-closing)
	result = strings.ReplaceAll(result, "</divider>", "")

	// Process link tags
	result = processTag(result, "sevk-link", func(attrs map[string]string, inner string) string {
		href := attrs["href"]
		if href == "" {
			href = "#"
		}
		target := attrs["target"]
		if target == "" {
			target = "_blank"
		}
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<a href="%s" target="%s" style="%s">%s</a>`, href, target, styleStr, inner)
	})

	// Process list tags
	result = processTag(result, "list", func(attrs map[string]string, inner string) string {
		listType := attrs["type"]
		tag := "ul"
		if listType == "ordered" {
			tag = "ol"
		}
		style := extractAllStyleAttributes(attrs)
		if _, ok := style["margin"]; !ok {
			style["margin"] = "0"
		}
		if lst, ok := attrs["list-style-type"]; ok {
			style["list-style-type"] = lst
		}
		styleStr := styleToString(style)
		classAttr := ""
		if class, ok := attrs["class"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, class)
		} else if className, ok := attrs["className"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, className)
		}
		return fmt.Sprintf(`<%s style="%s"%s>%s</%s>`, tag, styleStr, classAttr, inner, tag)
	})

	// Process list item tags
	result = processTag(result, "li", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		classAttr := ""
		if class, ok := attrs["class"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, class)
		} else if className, ok := attrs["className"]; ok {
			classAttr = fmt.Sprintf(` class="%s"`, className)
		}
		return fmt.Sprintf(`<li style="%s"%s>%s</li>`, styleStr, classAttr, inner)
	})

	// Process codeblock tags with chroma syntax highlighting
	result = processTag(result, "codeblock", func(attrs map[string]string, inner string) string {
		return processCodeBlock(attrs, inner)
	})

	// Clean up stray Sevk closing tags that weren't consumed
	strayClosingTags := []string{
		"</container>", "</section>", "</row>", "</column>",
		"</heading>", "</paragraph>", "</text>", "</button>", "</sevk-link>",
	}
	for _, tag := range strayClosingTags {
		result = strings.ReplaceAll(result, tag, "")
	}

	// Clean up wrapper tags
	wrapperPatterns := []string{
		`<sevk-email[^>]*>`, `</sevk-email>`,
		`<sevk-body[^>]*>`, `</sevk-body>`,
		`<email[^>]*>`, `</email>`,
		`<mail[^>]*>`, `</mail>`,
		`<body[^>]*>`, `</body>`,
	}
	for _, pattern := range wrapperPatterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, "")
	}

	return strings.TrimSpace(result)
}

// isTruthy checks if a value is truthy (non-nil, non-empty, non-zero, non-false)
func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case float64:
		return v != 0
	case int:
		return v != 0
	case []interface{}:
		return len(v) > 0
	default:
		return true
	}
}

// evaluateCondition evaluates a condition expression supporting ==, !=, &&, ||
func evaluateCondition(expr string, config map[string]interface{}) bool {
	trimmed := strings.TrimSpace(expr)

	// OR: split on ||, return true if any part is true
	if strings.Contains(trimmed, "||") {
		parts := strings.Split(trimmed, "||")
		for _, part := range parts {
			if evaluateCondition(part, config) {
				return true
			}
		}
		return false
	}

	// AND: split on &&, return true if all parts are true
	if strings.Contains(trimmed, "&&") {
		parts := strings.Split(trimmed, "&&")
		for _, part := range parts {
			if !evaluateCondition(part, config) {
				return false
			}
		}
		return true
	}

	// Equality: key == "value"
	eqRe := regexp.MustCompile(`^(\w+)\s*==\s*"([^"]*)"$`)
	if eqMatch := eqRe.FindStringSubmatch(trimmed); eqMatch != nil {
		val := config[eqMatch[1]]
		return fmt.Sprintf("%v", val) == eqMatch[2]
	}

	// Inequality: key != "value"
	neqRe := regexp.MustCompile(`^(\w+)\s*!=\s*"([^"]*)"$`)
	if neqMatch := neqRe.FindStringSubmatch(trimmed); neqMatch != nil {
		val := config[neqMatch[1]]
		return fmt.Sprintf("%v", val) != neqMatch[2]
	}

	// Simple truthy check
	return isTruthy(config[trimmed])
}

// renderTemplate processes template syntax with config values
func renderTemplate(template string, config map[string]interface{}) string {
	result := template

	// Process {%#each array as alias%}...{%/each%} loops
	eachRe := regexp.MustCompile(`\{%#each\s+(\w+)(?:\s+as\s+(\w+))?%\}(.*?)\{%/each%\}`)
	result = eachRe.ReplaceAllStringFunc(result, func(match string) string {
		sub := eachRe.FindStringSubmatch(match)
		if len(sub) < 4 {
			return ""
		}
		arrayKey := sub[1]
		alias := sub[2]
		if alias == "" {
			alias = "item"
		}
		body := sub[3]

		arr, ok := config[arrayKey].([]interface{})
		if !ok || len(arr) == 0 {
			return ""
		}

		var parts []string
		for _, item := range arr {
			itemStr := body
			if itemMap, ok := item.(map[string]interface{}); ok {
				// Replace {%alias.prop%} with item property values
				propRe := regexp.MustCompile(`\{%` + regexp.QuoteMeta(alias) + `\.(\w+)(?:\s*\?\?\s*([^%]+))?%\}`)
				itemStr = propRe.ReplaceAllStringFunc(itemStr, func(propMatch string) string {
					propSub := propRe.FindStringSubmatch(propMatch)
					if len(propSub) < 2 {
						return ""
					}
					propKey := propSub[1]
					fallback := ""
					if len(propSub) > 2 {
						fallback = strings.TrimSpace(propSub[2])
					}
					if v, exists := itemMap[propKey]; exists && isTruthy(v) {
						return fmt.Sprintf("%v", v)
					}
					return fallback
				})
			}
			parts = append(parts, itemStr)
		}
		return strings.Join(parts, "")
	})

	// Process nested {%#if key%}...{%/if%} from innermost outward.
	// We find the last {%#if in the string, then match its closing {%/if%},
	// ensuring we process innermost blocks first.
	ifOpenRe := regexp.MustCompile(`\{%#if\s+([^%]+)%\}`)

	for i := 0; i < 100; i++ {
		// Find all {%#if ...%} positions
		allOpens := ifOpenRe.FindAllStringIndex(result, -1)
		if len(allOpens) == 0 {
			break
		}

		// Process the last (innermost) opening tag
		processed := false
		for idx := len(allOpens) - 1; idx >= 0; idx-- {
			openLoc := allOpens[idx]
			openMatch := ifOpenRe.FindStringSubmatch(result[openLoc[0]:])
			key := openMatch[1]
			afterOpen := openLoc[1]

			// Find the first {%/if%} after this opening tag
			closeIdx := strings.Index(result[afterOpen:], "{%/if%}")
			if closeIdx < 0 {
				continue
			}
			bodyEnd := afterOpen + closeIdx
			body := result[afterOpen:bodyEnd]
			fullEnd := bodyEnd + len("{%/if%}")

			// If body contains another {%#if, skip -- not innermost
			if strings.Contains(body, "{%#if ") {
				continue
			}

			// Check for {%else%} in body
			elseIdx := strings.Index(body, "{%else%}")
			var replacement string
			if elseIdx >= 0 {
				trueBody := body[:elseIdx]
				falseBody := body[elseIdx+len("{%else%}"):]
				if evaluateCondition(key, config) {
					replacement = trueBody
				} else {
					replacement = falseBody
				}
			} else {
				if evaluateCondition(key, config) {
					replacement = body
				} else {
					replacement = ""
				}
			}

			result = result[:openLoc[0]] + replacement + result[fullEnd:]
			processed = true
			break
		}

		if !processed {
			break
		}
	}

	// Process {%variable ?? fallback%} with fallback
	fallbackRe := regexp.MustCompile(`\{%(\w+)\s*\?\?\s*([^%]+)%\}`)
	result = fallbackRe.ReplaceAllStringFunc(result, func(match string) string {
		sub := fallbackRe.FindStringSubmatch(match)
		key := sub[1]
		fallback := strings.TrimSpace(sub[2])
		if v, exists := config[key]; exists && isTruthy(v) {
			return fmt.Sprintf("%v", v)
		}
		return fallback
	})

	// Process simple {%variable%} injection
	simpleRe := regexp.MustCompile(`\{%(\w+)%\}`)
	result = simpleRe.ReplaceAllStringFunc(result, func(match string) string {
		sub := simpleRe.FindStringSubmatch(match)
		key := sub[1]
		if v, exists := config[key]; exists {
			return fmt.Sprintf("%v", v)
		}
		return ""
	})

	return result
}

// processBlockTag parses block config and renders the template
func processBlockTag(attrs map[string]string, inner string) string {
	template := strings.TrimSpace(inner)
	if template == "" {
		template = attrs["template"]
	}
	if template == "" {
		return ""
	}
	configStr := attrs["config"]
	if configStr == "" {
		configStr = "{}"
	}
	configStr = strings.ReplaceAll(configStr, "'", "\"")
	configStr = strings.ReplaceAll(configStr, "&quot;", "\"")
	configStr = strings.ReplaceAll(configStr, "&amp;", "&")
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		config = map[string]interface{}{}
	}
	return renderTemplate(template, config)
}

// resolveChromaStyle maps theme attribute values to chroma style names.
// "oneDark" maps to "dracula", other values are tried as-is with a "monokai" fallback.
func resolveChromaStyle(themeName string) *chroma.Style {
	chromaName := themeName
	switch themeName {
	case "oneDark":
		chromaName = "dracula"
	case "oneLight":
		chromaName = "monokailight"
	case "vscDarkPlus":
		chromaName = "monokai"
	}
	s := styles.Get(chromaName)
	if s == styles.Fallback {
		s = styles.Get("monokai")
	}
	return s
}

// processCodeBlock renders a <codeblock> with chroma syntax highlighting.
// It produces email-safe HTML with inline styles (no CSS classes) and wraps
// each source line in <p style="margin:0;min-height:1em"> to match the Node SDK output.
func processCodeBlock(attrs map[string]string, code string) string {
	if code == "" {
		return "<pre><code></code></pre>"
	}

	language := attrs["language"]
	if language == "" {
		language = "javascript"
	}
	themeName := attrs["theme"]
	if themeName == "" {
		themeName = "oneDark"
	}

	customStyle := extractAllStyleAttributes(attrs)

	// Resolve lexer
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}

	// Fallback: plain pre/code when language is not recognised
	if lexer == nil {
		style := map[string]string{
			"background-color": "#282c34",
			"color":            "#abb2bf",
			"border-radius":    "0.3em",
			"padding":          "1em",
			"white-space":      "pre",
			"font-family":      "'Fira Code', 'Fira Mono', Menlo, Consolas, 'DejaVu Sans Mono', monospace",
			"font-size":        "13px",
			"line-height":      "1.5",
			"width":            "100%",
			"box-sizing":       "border-box",
		}
		for k, v := range customStyle {
			style[k] = v
		}
		escaped := strings.ReplaceAll(strings.ReplaceAll(code, "<", "&lt;"), ">", "&gt;")
		return fmt.Sprintf(`<pre style="%s"><code>%s</code></pre>`, styleToString(style), escaped)
	}

	lexer = chroma.Coalesce(lexer)

	// Resolve style / theme
	chromaStyle := resolveChromaStyle(themeName)

	// Use chroma HTML formatter with inline styles (email-safe, no CSS classes).
	// WithClasses(false) is the default, which produces inline styles instead of CSS classes.
	formatter := html.New(html.WithClasses(false))

	// Tokenize
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		escaped := strings.ReplaceAll(strings.ReplaceAll(code, "<", "&lt;"), ">", "&gt;")
		return fmt.Sprintf(`<pre><code>%s</code></pre>`, escaped)
	}

	// Render with chroma into a buffer
	var buf bytes.Buffer
	if err := formatter.Format(&buf, chromaStyle, iterator); err != nil {
		escaped := strings.ReplaceAll(strings.ReplaceAll(code, "<", "&lt;"), ">", "&gt;")
		return fmt.Sprintf(`<pre><code>%s</code></pre>`, escaped)
	}

	// Chroma emits a full <pre …><span …><code>…</code></span></pre> block.
	// We need to extract the inner highlighted HTML, split it into source lines,
	// and re-wrap each line in <p style="margin:0;min-height:1em">.

	highlighted := buf.String()

	// Strip the outer <pre …> wrapper that chroma produces so we can rebuild it.
	// Chroma output looks like: <pre ...><code ...>CONTENT</code></pre>
	innerContent := highlighted
	if idx := strings.Index(innerContent, "<code"); idx >= 0 {
		// skip to after the closing > of <code …>
		end := strings.Index(innerContent[idx:], ">")
		if end >= 0 {
			innerContent = innerContent[idx+end+1:]
		}
	}
	if idx := strings.LastIndex(innerContent, "</code>"); idx >= 0 {
		innerContent = innerContent[:idx]
	}
	// Also strip any trailing </pre> or </span></pre>
	innerContent = strings.TrimSuffix(innerContent, "</pre>")
	innerContent = strings.TrimSuffix(innerContent, "</span>")

	// Split by newlines (chroma keeps \n characters in the output)
	lines := strings.Split(innerContent, "\n")
	// Remove trailing empty line that comes from a final \n
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	var linesHTML []string
	for _, line := range lines {
		if line == "" {
			line = " " // preserve empty lines
		}
		linesHTML = append(linesHTML, fmt.Sprintf(`<p style="margin:0;min-height:1em">%s</p>`, line))
	}

	// Build pre style
	preStyle := map[string]string{
		"background-color": "#282c34",
		"color":            "#abb2bf",
		"border-radius":    "0.3em",
		"padding":          "1em",
		"white-space":      "pre",
		"font-family":      "'Fira Code', 'Fira Mono', Menlo, Consolas, 'DejaVu Sans Mono', monospace",
		"font-size":        "13px",
		"line-height":      "1.5",
		"width":            "100%",
		"box-sizing":       "border-box",
		"overflow":         "auto",
		"margin":           "0.5em 0",
	}

	// Override background/color from chroma style entry if available
	bg := chromaStyle.Get(chroma.Background)
	if bg.Background.IsSet() {
		preStyle["background-color"] = bg.Background.String()
	}
	if bg.Colour.IsSet() {
		preStyle["color"] = bg.Colour.String()
	}

	for k, v := range customStyle {
		preStyle[k] = v
	}

	return fmt.Sprintf(`<pre style="%s"><code>%s</code></pre>`, styleToString(preStyle), strings.Join(linesHTML, ""))
}

// processButton processes button with MSO compatibility (like Node.js)
func processButton(attrs map[string]string, inner string) string {
	href := attrs["href"]
	if href == "" {
		href = "#"
	}
	style := extractAllStyleAttributes(attrs)

	// Parse padding
	paddingTop, paddingRight, paddingBottom, paddingLeft := parsePadding(style)

	y := paddingTop + paddingBottom
	textRaise := pxToPt(y)

	plFontWidth, plSpaceCount := computeFontWidthAndSpaceCount(paddingLeft)
	prFontWidth, prSpaceCount := computeFontWidthAndSpaceCount(paddingRight)

	buttonStyle := map[string]string{
		"line-height":     "100%",
		"text-decoration": "none",
		"display":         "inline-block",
		"max-width":       "100%",
		"mso-padding-alt": "0px",
	}

	// Merge with extracted styles
	for k, v := range style {
		buttonStyle[k] = v
	}

	// Override padding with parsed values
	buttonStyle["padding-top"] = fmt.Sprintf("%dpx", paddingTop)
	buttonStyle["padding-right"] = fmt.Sprintf("%dpx", paddingRight)
	buttonStyle["padding-bottom"] = fmt.Sprintf("%dpx", paddingBottom)
	buttonStyle["padding-left"] = fmt.Sprintf("%dpx", paddingLeft)

	styleStr := styleToString(buttonStyle)

	leftMsoSpaces := strings.Repeat("&#8202;", plSpaceCount)
	rightMsoSpaces := strings.Repeat("&#8202;", prSpaceCount)

	return fmt.Sprintf(
		`<a href="%s" target="_blank" style="%s"><!--[if mso]><i style="mso-font-width:%d%%;mso-text-raise:%d" hidden>%s</i><![endif]--><span style="max-width:100%%;display:inline-block;line-height:120%%;mso-padding-alt:0px;mso-text-raise:%d">%s</span><!--[if mso]><i style="mso-font-width:%d%%" hidden>%s&#8203;</i><![endif]--></a>`,
		href,
		styleStr,
		int(plFontWidth*100),
		textRaise,
		leftMsoSpaces,
		pxToPt(paddingBottom),
		inner,
		int(prFontWidth*100),
		rightMsoSpaces,
	)
}

// parsePadding parses padding values from style
func parsePadding(style map[string]string) (int, int, int, int) {
	if padding, ok := style["padding"]; ok {
		parts := strings.Fields(padding)
		switch len(parts) {
		case 1:
			val := parsePx(parts[0])
			return val, val, val, val
		case 2:
			vertical := parsePx(parts[0])
			horizontal := parsePx(parts[1])
			return vertical, horizontal, vertical, horizontal
		case 4:
			return parsePx(parts[0]), parsePx(parts[1]), parsePx(parts[2]), parsePx(parts[3])
		}
	}

	pt := parsePx(style["padding-top"])
	pr := parsePx(style["padding-right"])
	pb := parsePx(style["padding-bottom"])
	pl := parsePx(style["padding-left"])
	return pt, pr, pb, pl
}

func parsePx(s string) int {
	s = strings.TrimSuffix(s, "px")
	val, _ := strconv.Atoi(s)
	return val
}

// pxToPt converts px to pt for MSO
func pxToPt(px int) int {
	return (px * 3) / 4
}

// computeFontWidthAndSpaceCount computes font width and space count for MSO padding
func computeFontWidthAndSpaceCount(expectedWidth int) (float64, int) {
	if expectedWidth == 0 {
		return 0, 0
	}

	smallestSpaceCount := 0
	maxFontWidth := 5.0

	for {
		var requiredFontWidth float64
		if smallestSpaceCount > 0 {
			requiredFontWidth = float64(expectedWidth) / float64(smallestSpaceCount) / 2.0
		} else {
			requiredFontWidth = math.Inf(1)
		}

		if requiredFontWidth <= maxFontWidth {
			return requiredFontWidth, smallestSpaceCount
		}
		smallestSpaceCount++
	}
}

func processTag(content, tagName string, processor func(map[string]string, string) string) string {
	result := content
	openPattern := fmt.Sprintf(`(?i)<%s([^>]*)>`, tagName)
	closeRe := regexp.MustCompile(fmt.Sprintf(`(?i)</%s>`, tagName))
	nestedOpenRe := regexp.MustCompile(fmt.Sprintf(`(?i)<%s`, tagName))
	openRe := regexp.MustCompile(openPattern)

	// Process from innermost to outermost by finding the last opening tag
	// that has a matching close tag without any nested same tags
	maxIterations := 10000 // Safety limit
	iterations := 0

	for iterations < maxIterations {
		iterations++

		// Find all opening tags
		allLocs := openRe.FindAllStringIndex(result, -1)
		if len(allLocs) == 0 {
			break
		}

		// Find the innermost tag (one that has no nested same tags)
		processed := false
		for i := len(allLocs) - 1; i >= 0; i-- {
			loc := allLocs[i]
			start := loc[0]
			innerStart := loc[1]

			// Get attributes
			match := openRe.FindStringSubmatch(result[start:])
			attrsStr := ""
			if len(match) > 1 {
				attrsStr = match[1]
			}

			// Find the next close tag after this opening tag using regex
			// to avoid Unicode length-change issues with strings.ToLower/strings.Index
			closeLoc := closeRe.FindStringIndex(result[innerStart:])
			if closeLoc == nil {
				continue
			}

			innerEnd := innerStart + closeLoc[0]
			closeTagLen := closeLoc[1] - closeLoc[0]
			inner := result[innerStart:innerEnd]

			// Check if there's another opening tag inside
			if nestedOpenRe.MatchString(inner) {
				// This tag has nested same tags, skip it
				continue
			}

			// This is an innermost tag, process it
			attrs := parseAttributes(attrsStr)
			replacement := processor(attrs, inner)
			end := innerEnd + closeTagLen

			result = result[:start] + replacement + result[end:]
			processed = true
			break
		}

		if !processed {
			// No more tags to process
			break
		}
	}

	return result
}

func parseAttributes(attrsStr string) map[string]string {
	attrs := make(map[string]string)
	re := regexp.MustCompile(`([\w-]+)=(?:"([^"]*)"|'([^']*)')`)
	matches := re.FindAllStringSubmatch(attrsStr, -1)
	for _, match := range matches {
		if len(match) > 2 {
			if match[2] != "" {
				attrs[match[1]] = match[2]
			} else {
				attrs[match[1]] = match[3]
			}
		}
	}
	return attrs
}

// extractAllStyleAttributes extracts all style attributes from element attributes (like Node.js extractStyleAttributes)
func extractAllStyleAttributes(attrs map[string]string) map[string]string {
	style := make(map[string]string)

	// Typography attributes
	if v, ok := attrs["text-color"]; ok {
		style["color"] = v
	} else if v, ok := attrs["color"]; ok {
		style["color"] = v
	}
	if v, ok := attrs["background-color"]; ok {
		style["background-color"] = v
	}
	if v, ok := attrs["font-size"]; ok {
		style["font-size"] = v
	}
	if v, ok := attrs["font-family"]; ok {
		style["font-family"] = v
	}
	if v, ok := attrs["font-weight"]; ok {
		style["font-weight"] = v
	}
	if v, ok := attrs["line-height"]; ok {
		style["line-height"] = v
	}
	if v, ok := attrs["text-align"]; ok {
		style["text-align"] = v
	}
	if v, ok := attrs["text-decoration"]; ok {
		style["text-decoration"] = v
	}

	// Dimensions
	if v, ok := attrs["width"]; ok {
		style["width"] = v
	}
	if v, ok := attrs["height"]; ok {
		style["height"] = v
	}
	if v, ok := attrs["max-width"]; ok {
		style["max-width"] = v
	}
	if v, ok := attrs["max-height"]; ok {
		style["max-height"] = v
	}
	if v, ok := attrs["min-width"]; ok {
		style["min-width"] = v
	}
	if v, ok := attrs["min-height"]; ok {
		style["min-height"] = v
	}

	// Spacing - Padding
	if v, ok := attrs["padding"]; ok {
		style["padding"] = v
	} else {
		if v, ok := attrs["padding-top"]; ok {
			style["padding-top"] = v
		}
		if v, ok := attrs["padding-right"]; ok {
			style["padding-right"] = v
		}
		if v, ok := attrs["padding-bottom"]; ok {
			style["padding-bottom"] = v
		}
		if v, ok := attrs["padding-left"]; ok {
			style["padding-left"] = v
		}
	}

	// Spacing - Margin
	if v, ok := attrs["margin"]; ok {
		style["margin"] = v
	} else {
		if v, ok := attrs["margin-top"]; ok {
			style["margin-top"] = v
		}
		if v, ok := attrs["margin-right"]; ok {
			style["margin-right"] = v
		}
		if v, ok := attrs["margin-bottom"]; ok {
			style["margin-bottom"] = v
		}
		if v, ok := attrs["margin-left"]; ok {
			style["margin-left"] = v
		}
	}

	// Borders
	if v, ok := attrs["border"]; ok {
		style["border"] = v
	} else {
		if v, ok := attrs["border-top"]; ok {
			style["border-top"] = v
		}
		if v, ok := attrs["border-right"]; ok {
			style["border-right"] = v
		}
		if v, ok := attrs["border-bottom"]; ok {
			style["border-bottom"] = v
		}
		if v, ok := attrs["border-left"]; ok {
			style["border-left"] = v
		}
		if v, ok := attrs["border-color"]; ok {
			style["border-color"] = v
		}
		if v, ok := attrs["border-width"]; ok {
			style["border-width"] = v
		}
		if v, ok := attrs["border-style"]; ok {
			style["border-style"] = v
		}
	}

	// Border Radius
	if v, ok := attrs["border-radius"]; ok {
		style["border-radius"] = v
	} else {
		if v, ok := attrs["border-top-left-radius"]; ok {
			style["border-top-left-radius"] = v
		}
		if v, ok := attrs["border-top-right-radius"]; ok {
			style["border-top-right-radius"] = v
		}
		if v, ok := attrs["border-bottom-left-radius"]; ok {
			style["border-bottom-left-radius"] = v
		}
		if v, ok := attrs["border-bottom-right-radius"]; ok {
			style["border-bottom-right-radius"] = v
		}
	}

	// Background image
	if bgImg, ok := attrs["background-image"]; ok {
		style["background-image"] = fmt.Sprintf("url('%s')", bgImg)
		if v, ok := attrs["background-size"]; ok {
			style["background-size"] = v
		} else {
			style["background-size"] = "cover"
		}
		if v, ok := attrs["background-position"]; ok {
			style["background-position"] = v
		} else {
			style["background-position"] = "center"
		}
		if v, ok := attrs["background-repeat"]; ok {
			style["background-repeat"] = v
		} else {
			style["background-repeat"] = "no-repeat"
		}
	} else {
		if v, ok := attrs["background-size"]; ok {
			style["background-size"] = v
		}
		if v, ok := attrs["background-position"]; ok {
			style["background-position"] = v
		}
		if v, ok := attrs["background-repeat"]; ok {
			style["background-repeat"] = v
		}
	}

	return style
}

// styleToString converts style map to inline style string
func styleToString(style map[string]string) string {
	var parts []string
	for k, v := range style {
		parts = append(parts, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(parts, ";")
}
