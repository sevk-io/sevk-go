package markup

import (
	"fmt"
	"sort"
	"strings"
)

func extractAllStyleAttributes(attributes map[string]string) map[string]string {
	style := make(map[string]string)

	copyAttribute := func(source, target string) {
		if value := attributes[source]; value != "" {
			style[target] = value
		}
	}

	if value := attributes["text-color"]; value != "" {
		style["color"] = value
	} else {
		copyAttribute("color", "color")
	}

	for _, name := range []string{
		"background-color",
		"font-size",
		"font-family",
		"font-weight",
		"line-height",
		"text-align",
		"text-decoration",
		"width",
		"height",
		"max-width",
		"max-height",
		"min-width",
		"min-height",
		"padding",
		"padding-top",
		"padding-right",
		"padding-bottom",
		"padding-left",
		"margin",
		"margin-top",
		"margin-right",
		"margin-bottom",
		"margin-left",
		"border",
		"border-top",
		"border-right",
		"border-bottom",
		"border-left",
		"border-color",
		"border-width",
		"border-style",
		"border-radius",
		"border-top-left-radius",
		"border-top-right-radius",
		"border-bottom-left-radius",
		"border-bottom-right-radius",
	} {
		copyAttribute(name, name)
	}

	if backgroundImage := attributes["background-image"]; backgroundImage != "" {
		style["background-image"] = fmt.Sprintf("url('%s')", backgroundImage)
		setBackgroundAttribute(style, attributes, "background-size", "cover")
		setBackgroundAttribute(style, attributes, "background-position", "center")
		setBackgroundAttribute(style, attributes, "background-repeat", "no-repeat")
	} else {
		copyAttribute("background-size", "background-size")
		copyAttribute("background-position", "background-position")
		copyAttribute("background-repeat", "background-repeat")
	}

	return style
}

func setBackgroundAttribute(
	style map[string]string,
	attributes map[string]string,
	name string,
	fallback string,
) {
	if value := attributes[name]; value != "" {
		style[name] = value
	} else {
		style[name] = fallback
	}
}

var stylePriority = map[string]int{
	"color":                      1,
	"background-color":           2,
	"background-image":           3,
	"background-size":            4,
	"background-position":        5,
	"background-repeat":          6,
	"font-size":                  7,
	"font-family":                8,
	"font-weight":                9,
	"line-height":                10,
	"text-align":                 11,
	"text-decoration":            12,
	"width":                      13,
	"height":                     14,
	"max-width":                  15,
	"max-height":                 16,
	"min-width":                  17,
	"min-height":                 18,
	"padding":                    19,
	"padding-top":                20,
	"padding-right":              21,
	"padding-bottom":             22,
	"padding-left":               23,
	"margin":                     24,
	"margin-top":                 25,
	"margin-right":               26,
	"margin-bottom":              27,
	"margin-left":                28,
	"border":                     29,
	"border-color":               30,
	"border-width":               31,
	"border-style":               32,
	"border-top":                 33,
	"border-right":               34,
	"border-bottom":              35,
	"border-left":                36,
	"border-radius":              37,
	"border-top-left-radius":     38,
	"border-top-right-radius":    39,
	"border-bottom-left-radius":  40,
	"border-bottom-right-radius": 41,
	"border-collapse":            42,
	"list-style-type":            43,
	"vertical-align":             44,
}

func styleToString(style map[string]string) string {
	keys := make([]string, 0, len(style))
	for key := range style {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		left, leftKnown := stylePriority[keys[i]]
		right, rightKnown := stylePriority[keys[j]]
		if !leftKnown {
			left = 1 << 30
		}
		if !rightKnown {
			right = 1 << 30
		}
		if left != right {
			return left < right
		}
		return keys[i] < keys[j]
	})

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		if style[key] == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:%s", key, style[key]))
	}
	return strings.Join(parts, ";")
}
