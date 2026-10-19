package markup

import (
	"regexp"
	"strings"
)

var (
	divOpenPattern     = regexp.MustCompile(`(?i)<div`)
	divClosePattern    = regexp.MustCompile(`(?i)</div>`)
	hiddenStylePattern = regexp.MustCompile(`(?i)style=["'][^"']*display:\s*none`)
)

func findHiddenPreviewDiv(html string) (string, string, bool) {
	position := 0

	for {
		location := divOpenPattern.FindStringIndex(html[position:])
		if location == nil {
			return "", "", false
		}
		open := position + location[0]
		close := strings.IndexByte(html[open+4:], '>')
		if close < 0 {
			return "", "", false
		}
		close += open + 4

		if hiddenStylePattern.MatchString(html[open : close+1]) {
			closing := divClosePattern.FindStringIndex(html[close+1:])
			if closing == nil {
				return "", "", false
			}
			return html[close+1 : close+1+closing[0]], html[open : close+1+closing[1]], true
		}

		position = close + 1
	}
}
