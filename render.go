// Package sevk provides the official Go SDK for Sevk - Email Marketing Platform
package sevk

import (
	"github.com/sevk-io/sevk-go/markup"
)

// Render renders Sevk markup to email-compatible HTML
// This is a convenience wrapper around markup.GenerateEmailFromMarkup
func Render(markupContent string) string {
	if markupContent == "" {
		return ""
	}
	return markup.GenerateEmailFromMarkup(markupContent, nil)
}
