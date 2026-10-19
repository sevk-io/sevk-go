package sevk_test

import (
	"os"
	"testing"

	sevk "github.com/sevk-io/sevk-go"
	"github.com/stretchr/testify/require"
)

func TestSendsAMarkupEmailThroughTheSevkAPI(t *testing.T) {
	apiKey := os.Getenv("SEVK_API_KEY")
	from := os.Getenv("SEVK_TEST_FROM")
	to := os.Getenv("SEVK_TEST_TO")
	if apiKey == "" || from == "" || to == "" {
		t.Skip("SEVK_API_KEY, SEVK_TEST_FROM and SEVK_TEST_TO are required")
	}

	client := sevk.NewWithOptions(apiKey, sevk.Options{BaseURL: os.Getenv("SEVK_BASE_URL")})
	response, err := client.Emails.Send(sevk.SendEmailParams{
		From:    from,
		To:      sevk.StringOrList{to},
		Subject: "Sevk SDK send test",
		Markup:  "<section><paragraph>Hello from the Sevk Go SDK</paragraph></section>",
	})
	require.NoError(t, err)

	single, ok := response.(sevk.SendEmailIDResponse)
	require.True(t, ok)
	require.NotEmpty(t, single.ID)
}
