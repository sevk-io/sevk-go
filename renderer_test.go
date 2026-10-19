package sevk_test

import (
	"strings"
	"testing"

	sevk "github.com/sevk-io/sevk-go"
	"github.com/stretchr/testify/require"
)

func TestRendersMarkupIntoAnHTMLEmailDocument(t *testing.T) {
	html, err := sevk.Render("<section><paragraph>Hello world</paragraph></section>")
	require.NoError(t, err)

	require.True(t, strings.HasPrefix(html, "<!DOCTYPE html PUBLIC"))
	require.Contains(t, html, "Hello world")
	require.NotContains(t, html, "<paragraph>")
}
