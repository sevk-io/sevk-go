package markup

import (
	"fmt"
	"strings"
)

const (
	maxMarkupDepth          = 64
	maxTagIterations        = 200
	maxRenderOutput         = 5_000_000
	maxTagScanWork          = 30_000_000
	maxCodeblockLength      = 2500
	maxTotalCodeblockLength = 2500
	tooComplexMessage       = "Email markup is too complex to render"
)

var depthTagNames = map[string]bool{
	"block": true, "section": true, "column": true, "row": true, "container": true, "heading": true,
	"paragraph": true, "text": true, "button": true, "link": true, "sevk-link": true, "list": true,
	"li": true, "codeblock": true,
}

type RenderError struct {
	message string
}

func (err *RenderError) Error() string {
	return err.message
}

type renderBudget struct {
	scanWork    int
	highlighted int
}

func (budget *renderBudget) chargeScanWork(work int) {
	budget.scanWork += work
	if budget.scanWork > maxTagScanWork {
		failRender(tooComplexMessage)
	}
}

func failRender(message string) {
	panic(&RenderError{message: message})
}

func isTagNameByte(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || char == '-'
}

func assertMarkupDepth(content string) {
	depths := map[string]int{}
	index := 0

	for index < len(content) {
		open := strings.IndexByte(content[index:], '<')
		if open < 0 {
			break
		}
		open += index
		close := strings.IndexByte(content[open+1:], '>')
		if close < 0 {
			break
		}
		close += open + 1
		index = close + 1

		position := open + 1
		closing := content[position] == '/'
		if closing {
			position++
		}
		nameEnd := position
		for nameEnd < close && isTagNameByte(content[nameEnd]) {
			nameEnd++
		}
		name := strings.ToLower(content[position:nameEnd])
		if !depthTagNames[name] {
			continue
		}

		if closing {
			if depths[name] > 0 {
				depths[name]--
			}
		} else if content[close-1] != '/' {
			depths[name]++
			if depths[name] > maxMarkupDepth {
				failRender(fmt.Sprintf("Email markup nesting exceeds the maximum depth of %d", maxMarkupDepth))
			}
		}
	}
}
