package markup

import (
	"math"
	"strconv"
	"strings"
)

func parsePadding(style map[string]string) (int, int, int, int) {
	var top, right, bottom, left int

	if padding, ok := style["padding"]; ok {
		parts := strings.Fields(padding)
		switch len(parts) {
		case 1:
			value := parsePixels(parts[0])
			top, right, bottom, left = value, value, value, value
		case 2:
			vertical := parsePixels(parts[0])
			horizontal := parsePixels(parts[1])
			top, right, bottom, left = vertical, horizontal, vertical, horizontal
		case 3:
			top = parsePixels(parts[0])
			right = parsePixels(parts[1])
			bottom = parsePixels(parts[2])
			left = parsePixels(parts[1])
		case 4:
			top = parsePixels(parts[0])
			right = parsePixels(parts[1])
			bottom = parsePixels(parts[2])
			left = parsePixels(parts[3])
		}
	}

	if value, ok := style["padding-top"]; ok {
		top = parsePixels(value)
	}
	if value, ok := style["padding-right"]; ok {
		right = parsePixels(value)
	}
	if value, ok := style["padding-bottom"]; ok {
		bottom = parsePixels(value)
	}
	if value, ok := style["padding-left"]; ok {
		left = parsePixels(value)
	}

	return top, right, bottom, left
}

func parsePixels(value string) int {
	value = strings.TrimSpace(value)
	end := 0
	if strings.HasPrefix(value, "+") || strings.HasPrefix(value, "-") {
		end = 1
	}
	for end < len(value) && value[end] >= '0' && value[end] <= '9' {
		end++
	}
	pixels, _ := strconv.Atoi(value[:end])
	return pixels
}

func pxToPt(pixels int) int {
	return int(math.Floor(float64(pixels*3+2) / 4))
}

func computeFontWidthAndSpaceCount(expectedWidth int) (float64, int) {
	if expectedWidth <= 0 {
		return 0, 0
	}

	spaceCount := int(math.Ceil(float64(expectedWidth) / 10))
	if spaceCount < 1 {
		spaceCount = 1
	}
	if spaceCount > 1000 {
		spaceCount = 1000
	}
	return float64(expectedWidth) / float64(spaceCount) / 2.0, spaceCount
}
