package sevk

import (
	"github.com/sevk-io/sevk-go/markup"
)

func Render(markupContent string, headSettings ...*markup.EmailHeadSettings) (string, error) {
	return markup.Render(markupContent, headSettings...)
}
