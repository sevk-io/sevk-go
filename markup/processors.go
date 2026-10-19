package markup

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

func processSection(attributes map[string]string, inner string) string {
	style := extractAllStyleAttributes(attributes)
	cellStyle := map[string]string{}

	if value, ok := style["padding"]; ok {
		cellStyle["padding"] = value
		delete(style, "padding")
	}
	if value, ok := style["text-align"]; ok {
		cellStyle["text-align"] = value
		delete(style, "text-align")
	}
	for _, key := range []string{
		"border-radius",
		"border-top-left-radius",
		"border-top-right-radius",
		"border-bottom-left-radius",
		"border-bottom-right-radius",
	} {
		if _, ok := style[key]; ok {
			style["border-collapse"] = "separate"
			break
		}
	}

	return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody>
<tr>
<td style="%s">%s</td>
</tr>
</tbody>
</table>`, styleToString(style), styleToString(cellStyle), inner)
}

var containerCellStyles = map[string]bool{
	"background-color": true, "background-image": true, "background-size": true,
	"background-position": true, "background-repeat": true,
	"border": true, "border-top": true, "border-right": true, "border-bottom": true, "border-left": true,
	"border-color": true, "border-width": true, "border-style": true,
	"border-radius": true, "border-top-left-radius": true, "border-top-right-radius": true,
	"border-bottom-left-radius": true, "border-bottom-right-radius": true,
	"padding": true, "padding-top": true, "padding-right": true, "padding-bottom": true, "padding-left": true,
}

func processContainer(attributes map[string]string, inner string) string {
	style := extractAllStyleAttributes(attributes)
	cellStyle := make(map[string]string)
	tableStyle := make(map[string]string)

	for key, value := range style {
		if containerCellStyles[key] {
			cellStyle[key] = value
		} else {
			tableStyle[key] = value
		}
	}

	hasBorderRadius := cellStyle["border-radius"] != "" ||
		cellStyle["border-top-left-radius"] != "" ||
		cellStyle["border-top-right-radius"] != "" ||
		cellStyle["border-bottom-left-radius"] != "" ||
		cellStyle["border-bottom-right-radius"] != ""
	if hasBorderRadius {
		tableStyle["border-collapse"] = "separate"
	}

	if width, ok := tableStyle["width"]; ok && width != "100%" && width != "auto" {
		if _, hasMaxWidth := tableStyle["max-width"]; !hasMaxWidth {
			tableStyle["max-width"] = width
		}
		tableStyle["width"] = "100%"
	}

	return fmt.Sprintf(`<table width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" style="%s">
<tbody>
<tr style="width:100%%">
<td style="%s">%s</td>
</tr>
</tbody>
</table>`, styleToString(tableStyle), styleToString(cellStyle), inner)
}

func processColumn(attributes map[string]string, inner string) string {
	style := extractAllStyleAttributes(attributes)
	if _, ok := style["vertical-align"]; !ok {
		style["vertical-align"] = "top"
	}
	return fmt.Sprintf(`<td class="sevk-column" style="%s">%s</td>`, styleToString(style), inner)
}

func processRow(
	attributes map[string]string,
	inner string,
	rowCounter *int,
	gapRules *[]string,
) string {
	gap := attributes["gap"]
	if gap == "" {
		gap = "0"
	}
	style := extractAllStyleAttributes(attributes)
	delete(style, "gap")

	gapPixels := strings.Replace(gap, "px", "", 1)
	gapSize := parsePixels(gapPixels)
	rowID := fmt.Sprintf("sevk-row-%d", *rowCounter)
	*rowCounter++

	columnPattern := regexp.MustCompile(`class="sevk-column"`)
	columnCount := len(columnPattern.FindAllString(inner, -1))
	if columnCount > 1 {
		equalWidth := fmt.Sprintf("%d%%", 100/columnCount)
		columnWidthPattern := regexp.MustCompile(`<td class="sevk-column" style="([^"]*)"`)
		inner = columnWidthPattern.ReplaceAllStringFunc(inner, func(match string) string {
			if strings.Contains(match, "width:") {
				return match
			}
			parts := columnWidthPattern.FindStringSubmatch(match)
			if len(parts) > 1 {
				return fmt.Sprintf(`<td class="sevk-column" style="width:%s;%s"`, equalWidth, parts[1])
			}
			return match
		})
	}

	if gapSize > 0 {
		spacer := fmt.Sprintf(
			`</td><td class="sevk-gap" style="width:%spx;min-width:%spx" width="%s"></td><td class="sevk-column"`,
			gapPixels,
			gapPixels,
			gapPixels,
		)
		inner = regexp.MustCompile(`</td>\s*<td class="sevk-column"`).ReplaceAllString(inner, spacer)
		*gapRules = append(
			*gapRules,
			fmt.Sprintf(`.%s .sevk-gap{display:none !important;}.%s > tbody > tr > td.sevk-column{display:block !important;width:100%% !important;margin-bottom:%spx !important;}.%s > tbody > tr > td.sevk-column:last-of-type{margin-bottom:0 !important;}`, rowID, rowID, gapPixels, rowID),
		)
	}

	return fmt.Sprintf(`<table align="center" width="100%%" border="0" cellPadding="0" cellSpacing="0" role="presentation" class="sevk-row-table %s" style="%s">
<tbody style="width:100%%">
<tr style="width:100%%">%s</tr>
</tbody>
</table>`, rowID, styleToString(style), inner)
}

func processHeading(attributes map[string]string, inner string) string {
	level := attributes["level"]
	if level == "" {
		level = "1"
	}
	style := extractAllStyleAttributes(attributes)
	if _, ok := style["margin"]; !ok {
		style["margin"] = "0"
	}
	return fmt.Sprintf(`<h%s style="%s">%s</h%s>`, level, styleToString(style), inner, level)
}

func processParagraph(attributes map[string]string, inner string) string {
	style := extractAllStyleAttributes(attributes)
	if _, ok := style["margin"]; !ok {
		style["margin"] = "0"
	}
	return fmt.Sprintf(`<p style="%s">%s</p>`, styleToString(style), inner)
}

func processText(attributes map[string]string, inner string) string {
	return fmt.Sprintf(
		`<span style="%s">%s</span>`,
		styleToString(extractAllStyleAttributes(attributes)),
		inner,
	)
}

func processLink(attributes map[string]string, inner string) string {
	href := attributes["href"]
	if href == "" {
		href = "#"
	}
	target := attributes["target"]
	if target == "" {
		target = "_blank"
	}
	return fmt.Sprintf(
		`<a href="%s" target="%s" style="%s">%s</a>`,
		href,
		target,
		styleToString(extractAllStyleAttributes(attributes)),
		inner,
	)
}

// Outlook ignores normal anchor padding, so conditional MSO spacing mirrors it.
func processButton(attributes map[string]string, inner string) string {
	href := attributes["href"]
	if href == "" {
		href = "#"
	}
	style := extractAllStyleAttributes(attributes)
	top, right, bottom, left := parsePadding(style)
	leftFontWidth, leftSpaceCount := computeFontWidthAndSpaceCount(left)
	rightFontWidth, rightSpaceCount := computeFontWidthAndSpaceCount(right)

	buttonStyle := map[string]string{
		"line-height":     "100%",
		"text-decoration": "none",
		"display":         "inline-block",
		"max-width":       "100%",
		"mso-padding-alt": "0px",
	}
	for key, value := range style {
		buttonStyle[key] = value
	}
	buttonStyle["padding-top"] = fmt.Sprintf("%dpx", top)
	buttonStyle["padding-right"] = fmt.Sprintf("%dpx", right)
	buttonStyle["padding-bottom"] = fmt.Sprintf("%dpx", bottom)
	buttonStyle["padding-left"] = fmt.Sprintf("%dpx", left)

	return fmt.Sprintf(
		`<a href="%s" target="_blank" style="%s"><!--[if mso]><i style="mso-font-width:%d%%;mso-text-raise:%d" hidden>%s</i><![endif]--><span style="max-width:100%%;display:inline-block;line-height:120%%;mso-padding-alt:0px;mso-text-raise:%d">%s</span><!--[if mso]><i style="mso-font-width:%d%%" hidden>%s&#8203;</i><![endif]--></a>`,
		href,
		styleToString(buttonStyle),
		int(math.Floor(leftFontWidth*100+0.5)),
		pxToPt(top+bottom),
		strings.Repeat("&#8202;", leftSpaceCount),
		pxToPt(bottom),
		inner,
		int(math.Floor(rightFontWidth*100+0.5)),
		strings.Repeat("&#8202;", rightSpaceCount),
	)
}

func processImage(attributes map[string]string) string {
	source := attributes["src"]
	alt := attributes["alt"]
	style := extractAllStyleAttributes(attributes)
	if _, ok := style["vertical-align"]; !ok {
		style["vertical-align"] = "middle"
	}
	if _, ok := style["max-width"]; !ok {
		style["max-width"] = "100%"
	}

	width := ""
	if value := attributes["width"]; value != "" {
		width = fmt.Sprintf(` width="%s"`, strings.Replace(value, "px", "", 1))
	}
	height := ""
	if value := attributes["height"]; value != "" {
		height = fmt.Sprintf(` height="%s"`, strings.Replace(value, "px", "", 1))
	}

	return fmt.Sprintf(
		`<img src="%s" alt="%s"%s%s style="%s" />`,
		source,
		alt,
		width,
		height,
		styleToString(style),
	)
}

func processDivider(attributes map[string]string) string {
	return fmt.Sprintf(
		`<hr style="%s"%s />`,
		styleToString(extractAllStyleAttributes(attributes)),
		classAttribute(attributes),
	)
}

func processList(attributes map[string]string, inner string) string {
	tag := "ul"
	if attributes["type"] == "ordered" {
		tag = "ol"
	}
	style := extractAllStyleAttributes(attributes)
	if _, ok := style["margin"]; !ok {
		style["margin"] = "0"
	}
	if listStyle, ok := attributes["list-style-type"]; ok {
		style["list-style-type"] = listStyle
	}
	return fmt.Sprintf(
		`<%s style="%s"%s>%s</%s>`,
		tag,
		styleToString(style),
		classAttribute(attributes),
		inner,
		tag,
	)
}

func processListItem(attributes map[string]string, inner string) string {
	return fmt.Sprintf(
		`<li style="%s"%s>%s</li>`,
		styleToString(extractAllStyleAttributes(attributes)),
		classAttribute(attributes),
		inner,
	)
}

func classAttribute(attributes map[string]string) string {
	className := attributes["class"]
	if className == "" {
		className = attributes["className"]
	}
	if className == "" {
		return ""
	}
	return fmt.Sprintf(` class="%s"`, className)
}
