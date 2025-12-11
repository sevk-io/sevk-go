// Package markup provides email markup rendering for Sevk
package markup

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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
}

// ParsedEmailContent represents parsed email content
type ParsedEmailContent struct {
	Body         string
	HeadSettings EmailHeadSettings
}

// GenerateEmailFromMarkup generates email HTML from Sevk markup
func GenerateEmailFromMarkup(htmlContent string, headSettings *EmailHeadSettings) string {
	var contentToProcess string
	var settings EmailHeadSettings

	if headSettings != nil {
		contentToProcess = htmlContent
		settings = *headSettings
	} else {
		parsed := ParseEmailHTML(htmlContent)
		contentToProcess = parsed.Body
		settings = parsed.HeadSettings
	}

	normalized := normalizeMarkup(contentToProcess)
	processed := processMarkup(normalized)

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
<html lang="en" dir="ltr">
<head>
<meta content="text/html; charset=UTF-8" http-equiv="Content-Type"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
%s
%s
%s
</head>
<body style="margin:0;padding:0;font-family:ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;background-color:#ffffff">
%s
%s
</body>
</html>`, titleTag, fontLinks, customStyles, previewText, processed)
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
	result = processTag(result, "row", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody style="width:100%%">
<tr style="width:100%%">%s</tr>
</tbody>
</table>`, styleStr, inner)
	})

	// Process column tags
	result = processTag(result, "column", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<td style="%s">%s</td>`, styleStr, inner)
	})

	// Process container tags
	result = processTag(result, "container", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody>
<tr style="width:100%%">
<td>%s</td>
</tr>
</tbody>
</table>`, styleStr, inner)
	})

	// Process heading tags
	result = processTag(result, "heading", func(attrs map[string]string, inner string) string {
		level := attrs["level"]
		if level == "" {
			level = "1"
		}
		style := extractAllStyleAttributes(attrs)
		styleStr := styleToString(style)
		return fmt.Sprintf(`<h%s style="%s">%s</h%s>`, level, styleStr, inner, level)
	})

	// Process paragraph tags
	result = processTag(result, "paragraph", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
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
			widthAttr = fmt.Sprintf(` width="%s"`, width)
		}
		heightAttr := ""
		if height != "" {
			heightAttr = fmt.Sprintf(` height="%s"`, height)
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

	// Process codeblock tags
	result = processTag(result, "codeblock", func(attrs map[string]string, inner string) string {
		style := extractAllStyleAttributes(attrs)
		if _, ok := style["width"]; !ok {
			style["width"] = "100%"
		}
		if _, ok := style["box-sizing"]; !ok {
			style["box-sizing"] = "border-box"
		}
		styleStr := styleToString(style)
		escaped := strings.ReplaceAll(strings.ReplaceAll(inner, "<", "&lt;"), ">", "&gt;")
		return fmt.Sprintf(`<pre style="%s"><code>%s</code></pre>`, styleStr, escaped)
	})

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
	openPattern := fmt.Sprintf(`<%s([^>]*)>`, tagName)
	closeTag := fmt.Sprintf("</%s>", tagName)
	openRe := regexp.MustCompile(openPattern)

	// Process from innermost to outermost by finding the last opening tag
	// that has a matching close tag without any nested same tags
	maxIterations := 10000 // Safety limit
	iterations := 0

	for iterations < maxIterations {
		iterations++

		// Find all opening tags
		allLocs := openRe.FindAllStringIndex(result, -1)
		if allLocs == nil || len(allLocs) == 0 {
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

			// Find the next close tag after this opening tag
			closePos := strings.Index(result[innerStart:], closeTag)
			if closePos == -1 {
				continue
			}

			innerEnd := innerStart + closePos
			inner := result[innerStart:innerEnd]

			// Check if there's another opening tag inside
			if openRe.MatchString(inner) {
				// This tag has nested same tags, skip it
				continue
			}

			// This is an innermost tag, process it
			attrs := parseAttributes(attrsStr)
			replacement := processor(attrs, inner)
			end := innerEnd + len(closeTag)

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
	re := regexp.MustCompile(`([\w-]+)=["']([^"']*)["']`)
	matches := re.FindAllStringSubmatch(attrsStr, -1)
	for _, match := range matches {
		if len(match) > 2 {
			attrs[match[1]] = match[2]
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
