package sevk_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	sevk "github.com/sevk-io/sevk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseURL string

var testClient *sevk.Client
var testAudience *sevk.Audience

func uniqueID() int64 {
	return time.Now().UnixMilli() + int64(rand.Intn(10000))
}

func TestMain(m *testing.M) {
	apiKey := os.Getenv("SEVK_TEST_API_KEY")
	if apiKey == "" {
		fmt.Println("SEVK_TEST_API_KEY not set, skipping integration tests")
		os.Exit(0)
	}

	baseURL = os.Getenv("SEVK_TEST_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.sevk.io"
	}

	testClient = sevk.NewWithOptions(apiKey, sevk.Options{BaseURL: baseURL})

	// Create a shared test audience for all tests
	audience, err := testClient.Audiences.Create(sevk.CreateAudienceParams{
		Name: fmt.Sprintf("Shared Test Audience %d", uniqueID()),
	})
	if err != nil {
		fmt.Printf("Skipping tests - audience creation failed: %v\n", err)
		os.Exit(0)
	}
	testAudience = audience

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func getClient(t *testing.T) *sevk.Client {
	if testClient == nil {
		t.Skip("Test client not initialized")
	}
	return testClient
}

func getAudience(t *testing.T) *sevk.Audience {
	if testAudience == nil {
		t.Skip("Test audience not initialized")
	}
	return testAudience
}

func skipDomainTests(t *testing.T) {
	if os.Getenv("INCLUDE_DOMAIN_TESTS") != "true" {
		t.Skip("INCLUDE_DOMAIN_TESTS not set")
	}
}

// ============================================
// AUTHENTICATION TESTS
// ============================================

func TestAuthInvalidAPIKey(t *testing.T) {
	client := sevk.NewWithOptions("sevk_invalid_api_key_12345", sevk.Options{BaseURL: baseURL})

	_, err := client.Contacts.List(nil)
	assert.Error(t, err)

	sevkErr, ok := err.(*sevk.Error)
	if ok {
		assert.True(t, sevkErr.IsUnauthorized())
	}
}

func TestAuthEmptyAPIKey(t *testing.T) {
	client := sevk.NewWithOptions("", sevk.Options{BaseURL: baseURL})

	_, err := client.Contacts.List(nil)
	assert.Error(t, err)
}

func TestAuthMalformedAPIKey(t *testing.T) {
	client := sevk.NewWithOptions("invalid_key_format", sevk.Options{BaseURL: baseURL})

	_, err := client.Contacts.List(nil)
	assert.Error(t, err)
}

// ============================================
// CONTACTS TESTS
// ============================================

func TestContactsListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Contacts.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
	assert.GreaterOrEqual(t, result.Page, 1)
}

func TestContactsListPagination(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 5
	result, err := client.Contacts.List(&sevk.ListContactsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
}

func TestContactsCreateRequired(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("test-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)
	assert.NotEmpty(t, contact.ID)
	assert.Equal(t, email, contact.Email)
}

func TestContactsGetByID(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("get-test-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	fetched, err := client.Contacts.Get(contact.ID)
	require.NoError(t, err)
	assert.Equal(t, contact.ID, fetched.ID)
	assert.Equal(t, email, fetched.Email)
}

func TestContactsUpdate(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("update-test-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	subscribed := false
	updated, err := client.Contacts.Update(contact.ID, sevk.UpdateContactParams{Subscribed: &subscribed})
	require.NoError(t, err)
	assert.Equal(t, contact.ID, updated.ID)
	assert.Equal(t, false, updated.Subscribed)
}

func TestContactsDelete(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("delete-test-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	err = client.Contacts.Delete(contact.ID)
	require.NoError(t, err)

	_, err = client.Contacts.Get(contact.ID)
	assert.Error(t, err)
	sevkErr, ok := err.(*sevk.Error)
	if ok {
		assert.True(t, sevkErr.IsNotFound())
	}
}

func TestContactsErrorNonExistent(t *testing.T) {
	client := getClient(t)

	_, err := client.Contacts.Get("non-existent-id")
	assert.Error(t, err)
	sevkErr, ok := err.(*sevk.Error)
	if ok {
		assert.True(t, sevkErr.IsNotFound())
	}
}

func TestContactsBulkUpdate(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("bulk-update-%d@example.com", uniqueID())
	_, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	// Wait a bit to avoid rate limit conflicts with import endpoint
	time.Sleep(2 * time.Second)

	subscribed := false
	result, err := client.Contacts.BulkUpdate(sevk.BulkUpdateContactsParams{
		Contacts: []sevk.BulkUpdateContact{
			{Email: email, Subscribed: &subscribed},
		},
	})
	if err != nil {
		sevkErr, ok := err.(*sevk.Error)
		if ok && sevkErr.StatusCode == 429 {
			t.Skip("Skipping due to rate limit on bulk update endpoint")
		}
		require.NoError(t, err)
	}
	assert.NotNil(t, result)
}

func TestContactsGetEvents(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("events-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	result, err := client.Contacts.GetEvents(contact.ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestContactsImport(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("import-test-%d@example.com", uniqueID())
	result, err := client.Contacts.Import(sevk.ImportContactsParams{
		Contacts: []sevk.ImportContact{
			{Email: email},
		},
	})
	if err != nil {
		sevkErr, ok := err.(*sevk.Error)
		if ok && sevkErr.StatusCode == 429 {
			t.Skip("Skipping due to rate limit on import endpoint")
		}
		require.NoError(t, err)
	}
	assert.NotNil(t, result)
}

// ============================================
// AUDIENCES TESTS
// ============================================

func TestAudiencesListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Audiences.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
	assert.GreaterOrEqual(t, result.Page, 1)
}

func TestAudiencesCreateRequired(t *testing.T) {
	client := getClient(t)

	name := fmt.Sprintf("Test Audience %d", uniqueID())
	audience, err := client.Audiences.Create(sevk.CreateAudienceParams{Name: name})
	require.NoError(t, err)
	assert.NotEmpty(t, audience.ID)
	assert.Equal(t, name, audience.Name)

	// Cleanup
	client.Audiences.Delete(audience.ID)
}

func TestAudiencesCreateAllFields(t *testing.T) {
	client := getClient(t)

	name := fmt.Sprintf("Full Audience %d", uniqueID())
	description := "Test description"
	usersCanSee := "PUBLIC"
	audience, err := client.Audiences.Create(sevk.CreateAudienceParams{
		Name:        name,
		Description: &description,
		UsersCanSee: &usersCanSee,
	})
	require.NoError(t, err)
	assert.Equal(t, name, audience.Name)

	// Cleanup
	client.Audiences.Delete(audience.ID)
}

func TestAudiencesGetByID(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	fetched, err := client.Audiences.Get(audience.ID)
	require.NoError(t, err)
	assert.Equal(t, audience.ID, fetched.ID)
}

func TestAudiencesUpdate(t *testing.T) {
	client := getClient(t)

	name := fmt.Sprintf("Update Audience %d", uniqueID())
	audience, err := client.Audiences.Create(sevk.CreateAudienceParams{Name: name})
	require.NoError(t, err)

	newName := fmt.Sprintf("Updated Audience %d", uniqueID())
	updated, err := client.Audiences.Update(audience.ID, sevk.UpdateAudienceParams{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// Cleanup
	client.Audiences.Delete(audience.ID)
}

func TestAudiencesAddContacts(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	email := fmt.Sprintf("add-contact-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	err = client.Audiences.AddContacts(audience.ID, []string{contact.ID})
	require.NoError(t, err)
}

func TestAudiencesListContacts(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	result, err := client.Audiences.ListContacts(audience.ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Items)
}

func TestAudiencesRemoveContact(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	email := fmt.Sprintf("audience-remove-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	err = client.Audiences.AddContacts(audience.ID, []string{contact.ID})
	require.NoError(t, err)

	err = client.Audiences.RemoveContact(audience.ID, contact.ID)
	require.NoError(t, err)

	// Verify removal by listing contacts
	result, err := client.Audiences.ListContacts(audience.ID, nil)
	require.NoError(t, err)
	for _, c := range result.Items {
		assert.NotEqual(t, contact.ID, c.ID)
	}
}

func TestAudiencesDelete(t *testing.T) {
	client := getClient(t)

	name := fmt.Sprintf("Delete Test %d", uniqueID())
	audience, err := client.Audiences.Create(sevk.CreateAudienceParams{Name: name})
	require.NoError(t, err)

	err = client.Audiences.Delete(audience.ID)
	require.NoError(t, err)

	_, err = client.Audiences.Get(audience.ID)
	assert.Error(t, err)
}

// ============================================
// TEMPLATES TESTS
// ============================================

func TestTemplatesListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Templates.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
}

func TestTemplatesCreateRequired(t *testing.T) {
	client := getClient(t)

	title := fmt.Sprintf("Test Template %d", uniqueID())
	content := "<p>Hello {{name}}</p>"
	template, err := client.Templates.Create(sevk.CreateTemplateParams{Title: title, Content: content})
	require.NoError(t, err)
	assert.NotEmpty(t, template.ID)
	assert.Equal(t, title, template.Title)
	assert.Equal(t, content, template.Content)

	// Cleanup
	client.Templates.Delete(template.ID)
}

func TestTemplatesGetByID(t *testing.T) {
	client := getClient(t)

	title := fmt.Sprintf("Get Template %d", uniqueID())
	template, err := client.Templates.Create(sevk.CreateTemplateParams{Title: title, Content: "<p>Test</p>"})
	require.NoError(t, err)

	fetched, err := client.Templates.Get(template.ID)
	require.NoError(t, err)
	assert.Equal(t, template.ID, fetched.ID)

	// Cleanup
	client.Templates.Delete(template.ID)
}

func TestTemplatesUpdate(t *testing.T) {
	client := getClient(t)

	title := fmt.Sprintf("Test Template %d", uniqueID())
	template, err := client.Templates.Create(sevk.CreateTemplateParams{Title: title, Content: "<p>Test</p>"})
	require.NoError(t, err)

	newTitle := fmt.Sprintf("Updated Template %d", uniqueID())
	updated, err := client.Templates.Update(template.ID, sevk.UpdateTemplateParams{Title: &newTitle})
	require.NoError(t, err)
	assert.Equal(t, newTitle, updated.Title)

	// Cleanup
	client.Templates.Delete(template.ID)
}

func TestTemplatesDuplicate(t *testing.T) {
	client := getClient(t)

	title := fmt.Sprintf("Test Template %d", uniqueID())
	template, err := client.Templates.Create(sevk.CreateTemplateParams{Title: title, Content: "<p>Test</p>"})
	require.NoError(t, err)

	duplicated, err := client.Templates.Duplicate(template.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, duplicated.ID)
	assert.NotEqual(t, template.ID, duplicated.ID)

	// Cleanup
	client.Templates.Delete(template.ID)
	client.Templates.Delete(duplicated.ID)
}

func TestTemplatesDelete(t *testing.T) {
	client := getClient(t)

	title := fmt.Sprintf("Delete Test %d", uniqueID())
	template, err := client.Templates.Create(sevk.CreateTemplateParams{Title: title, Content: "<p>Test</p>"})
	require.NoError(t, err)

	err = client.Templates.Delete(template.ID)
	require.NoError(t, err)

	_, err = client.Templates.Get(template.ID)
	assert.Error(t, err)
}

// ============================================
// BROADCASTS TESTS
// ============================================

func TestBroadcastsListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Broadcasts.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
}

func TestBroadcastsListPagination(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 5
	result, err := client.Broadcasts.List(&sevk.ListBroadcastsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
}

func TestBroadcastsListSearch(t *testing.T) {
	client := getClient(t)

	search := "nonexistent"
	result, err := client.Broadcasts.List(&sevk.ListBroadcastsParams{Search: &search})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBroadcastsGetStatus(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 1
	result, err := client.Broadcasts.List(&sevk.ListBroadcastsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)

	if len(result.Items) == 0 {
		t.Skip("No broadcasts available to test GetStatus")
	}

	status, err := client.Broadcasts.GetStatus(result.Items[0].ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
}

func TestBroadcastsGetEmails(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 1
	result, err := client.Broadcasts.List(&sevk.ListBroadcastsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)

	if len(result.Items) == 0 {
		t.Skip("No broadcasts available to test GetEmails")
	}

	emails, err := client.Broadcasts.GetEmails(result.Items[0].ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, emails)
	assert.NotNil(t, emails.Items)
}

func TestBroadcastsEstimateCost(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 1
	result, err := client.Broadcasts.List(&sevk.ListBroadcastsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)

	if len(result.Items) == 0 {
		t.Skip("No broadcasts available to test EstimateCost")
	}

	estimate, err := client.Broadcasts.EstimateCost(result.Items[0].ID)
	require.NoError(t, err)
	assert.NotNil(t, estimate)
}

func TestBroadcastsListActive(t *testing.T) {
	client := getClient(t)

	result, err := client.Broadcasts.ListActive()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBroadcastsCreate(t *testing.T) {
	client := getClient(t)

	// Get a domain from the project to use for broadcast
	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast create")
	}
	domainID := domains.Items[0].ID

	name := fmt.Sprintf("Test Broadcast %d", uniqueID())
	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       name,
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test broadcast body</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, broadcast.ID)
	assert.Equal(t, "DRAFT", broadcast.Status)

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsGet(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast get")
	}
	domainID := domains.Items[0].ID

	name := fmt.Sprintf("Get Broadcast %d", uniqueID())
	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       name,
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	fetched, err := client.Broadcasts.Get(broadcast.ID)
	require.NoError(t, err)
	assert.Equal(t, broadcast.ID, fetched.ID)
	assert.Equal(t, "Test Subject", fetched.Subject)

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsUpdate(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast update")
	}
	domainID := domains.Items[0].ID

	name := fmt.Sprintf("Update Broadcast %d", uniqueID())
	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       name,
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	newName := fmt.Sprintf("Updated Broadcast %d", uniqueID())
	updated, err := client.Broadcasts.Update(broadcast.ID, sevk.UpdateBroadcastParams{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, broadcast.ID, updated.ID)

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsGetAnalytics(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast analytics")
	}
	domainID := domains.Items[0].ID

	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       fmt.Sprintf("Analytics Broadcast %d", uniqueID()),
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	analytics, err := client.Broadcasts.GetAnalytics(broadcast.ID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsSendTest(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast send test")
	}
	domainID := domains.Items[0].ID

	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       fmt.Sprintf("SendTest Broadcast %d", uniqueID()),
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	err = client.Broadcasts.SendTest(broadcast.ID, sevk.SendTestParams{
		Emails: []string{"test@example.com"},
	})
	// May fail if domain is unverified, which is expected
	if err != nil {
		sevkErr, ok := err.(*sevk.Error)
		if ok {
			assert.NotEmpty(t, sevkErr.Message)
		}
	}

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsSendErrorDraft(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast send error")
	}
	domainID := domains.Items[0].ID

	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       fmt.Sprintf("Send Error Broadcast %d", uniqueID()),
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	_, sendErr := client.Broadcasts.Send(broadcast.ID)
	// Expected to fail if broadcast is not ready to send
	if sendErr != nil {
		sevkErr, ok := sendErr.(*sevk.Error)
		if ok {
			assert.NotEmpty(t, sevkErr.Message)
		}
	}

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsCancelErrorDraft(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast cancel error")
	}
	domainID := domains.Items[0].ID

	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       fmt.Sprintf("Cancel Error Broadcast %d", uniqueID()),
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	_, cancelErr := client.Broadcasts.Cancel(broadcast.ID)
	// Expected to fail if broadcast is not in a cancellable state
	if cancelErr != nil {
		sevkErr, ok := cancelErr.(*sevk.Error)
		if ok {
			assert.NotEmpty(t, sevkErr.Message)
		}
	}

	// Cleanup
	client.Broadcasts.Delete(broadcast.ID)
}

func TestBroadcastsDelete(t *testing.T) {
	client := getClient(t)

	domains, err := client.Domains.List(nil)
	require.NoError(t, err)
	if len(domains.Items) == 0 {
		t.Skip("No domains available to test broadcast delete")
	}
	domainID := domains.Items[0].ID

	senderName := "Test Sender"
	targetType := "ALL"
	broadcast, err := client.Broadcasts.Create(sevk.CreateBroadcastParams{
		DomainID:   domainID,
		Name:       fmt.Sprintf("Delete Broadcast %d", uniqueID()),
		Subject:    "Test Subject",
		Body:       "<section><paragraph>Test</paragraph></section>",
		SenderName: &senderName,
		TargetType: &targetType,
	})
	require.NoError(t, err)

	err = client.Broadcasts.Delete(broadcast.ID)
	require.NoError(t, err)

	_, err = client.Broadcasts.Get(broadcast.ID)
	assert.Error(t, err)
}

// ============================================
// DOMAINS TESTS
// ============================================

func TestDomainsListStructure(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	result, err := client.Domains.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDomainsListVerified(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	verified := true
	result, err := client.Domains.List(&sevk.ListDomainsParams{Verified: &verified})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDomainsCreate(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	subdomain := fmt.Sprintf("test-%d.example.com", uniqueID())
	domain, err := client.Domains.Create(sevk.CreateDomainParams{
		Domain: subdomain,
		Email:  fmt.Sprintf("test@%s", subdomain),
		From:   "noreply",
		SenderName: "Test Sender",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, domain.ID)
	assert.Equal(t, subdomain, domain.Domain)

	// Cleanup
	client.Domains.Delete(domain.ID)
}

func TestDomainsGet(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	subdomain := fmt.Sprintf("get-test-%d.example.com", uniqueID())
	domain, err := client.Domains.Create(sevk.CreateDomainParams{
		Domain: subdomain,
		Email:  fmt.Sprintf("test@%s", subdomain),
		From:   "noreply",
		SenderName: "Test Sender",
	})
	require.NoError(t, err)

	fetched, err := client.Domains.Get(domain.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.ID, fetched.ID)

	// Cleanup
	client.Domains.Delete(domain.ID)
}

func TestDomainsGetDnsRecords(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	subdomain := fmt.Sprintf("dns-test-%d.example.com", uniqueID())
	domain, err := client.Domains.Create(sevk.CreateDomainParams{
		Domain: subdomain,
		Email:  fmt.Sprintf("test@%s", subdomain),
		From:   "noreply",
		SenderName: "Test Sender",
	})
	require.NoError(t, err)

	records, err := client.Domains.GetDnsRecords(domain.ID)
	require.NoError(t, err)
	assert.NotNil(t, records)

	// Cleanup
	client.Domains.Delete(domain.ID)
}

func TestDomainsGetRegions(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	regions, err := client.Domains.GetRegions()
	require.NoError(t, err)
	assert.NotNil(t, regions)
}

func TestDomainsVerify(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	subdomain := fmt.Sprintf("verify-test-%d.example.com", uniqueID())
	domain, err := client.Domains.Create(sevk.CreateDomainParams{
		Domain: subdomain,
		Email:  fmt.Sprintf("test@%s", subdomain),
		From:   "noreply",
		SenderName: "Test Sender",
	})
	require.NoError(t, err)

	// Verification is expected to fail for test domains without proper DNS records
	_, verifyErr := client.Domains.Verify(domain.ID)
	if verifyErr != nil {
		sevkErr, ok := verifyErr.(*sevk.Error)
		if ok {
			assert.NotEmpty(t, sevkErr.Message)
		}
	}

	// Cleanup
	client.Domains.Delete(domain.ID)
}

func TestDomainsDelete(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	subdomain := fmt.Sprintf("delete-test-%d.example.com", uniqueID())
	domain, err := client.Domains.Create(sevk.CreateDomainParams{
		Domain: subdomain,
		Email:  fmt.Sprintf("test@%s", subdomain),
		From:   "noreply",
		SenderName: "Test Sender",
	})
	require.NoError(t, err)

	err = client.Domains.Delete(domain.ID)
	require.NoError(t, err)

	_, err = client.Domains.Get(domain.ID)
	assert.Error(t, err)
}

func TestDomainsUpdate(t *testing.T) {
	skipDomainTests(t)
	client := getClient(t)

	result, err := client.Domains.List(nil)
	require.NoError(t, err)

	if len(result.Items) == 0 {
		t.Skip("No domains available to test Update")
	}

	domain := result.Items[0]
	clickTracking := true
	updated, err := client.Domains.Update(domain.ID, sevk.UpdateDomainRequest{ClickTracking: &clickTracking})
	require.NoError(t, err)
	assert.Equal(t, domain.ID, updated.ID)
}

// ============================================
// TOPICS TESTS
// ============================================

func TestTopicsList(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	listResult, err := client.Topics.List(audience.ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, listResult)
}

func TestTopicsCreate(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("Test Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)
	assert.NotEmpty(t, topic.ID)
	assert.Equal(t, topicName, topic.Name)

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

func TestTopicsGetByID(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("Get Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	fetched, err := client.Topics.Get(audience.ID, topic.ID)
	require.NoError(t, err)
	assert.Equal(t, topic.ID, fetched.ID)

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

func TestTopicsUpdate(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("Update Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	newName := fmt.Sprintf("Updated Topic %d", uniqueID())
	updated, err := client.Topics.Update(audience.ID, topic.ID, sevk.UpdateTopicParams{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

func TestTopicsDelete(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("Delete Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	err = client.Topics.Delete(audience.ID, topic.ID)
	require.NoError(t, err)

	_, err = client.Topics.Get(audience.ID, topic.ID)
	assert.Error(t, err)
}

func TestTopicsAddContacts(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("AddContacts Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	email := fmt.Sprintf("topic-add-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	// Add contact to audience first
	err = client.Audiences.AddContacts(audience.ID, []string{contact.ID})
	require.NoError(t, err)

	// Add contact to topic
	err = client.Topics.AddContacts(audience.ID, topic.ID, sevk.AddTopicContactsParams{
		ContactIDs: []string{contact.ID},
	})
	require.NoError(t, err)

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

func TestTopicsRemoveContact(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("RemoveContact Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	email := fmt.Sprintf("topic-remove-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	// Add contact to audience and topic
	err = client.Audiences.AddContacts(audience.ID, []string{contact.ID})
	require.NoError(t, err)
	err = client.Topics.AddContacts(audience.ID, topic.ID, sevk.AddTopicContactsParams{
		ContactIDs: []string{contact.ID},
	})
	require.NoError(t, err)

	// Remove contact from topic
	err = client.Topics.RemoveContact(audience.ID, topic.ID, contact.ID)
	require.NoError(t, err)

	// Verify removal by listing contacts in the topic
	result, err := client.Topics.ListContacts(audience.ID, topic.ID, nil)
	require.NoError(t, err)
	for _, c := range result.Items {
		assert.NotEqual(t, contact.ID, c.ID)
	}

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

func TestTopicsListContacts(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	topicName := fmt.Sprintf("ListContacts Topic %d", uniqueID())
	topic, err := client.Topics.Create(audience.ID, sevk.CreateTopicParams{Name: topicName})
	require.NoError(t, err)

	result, err := client.Topics.ListContacts(audience.ID, topic.ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Items)
	assert.GreaterOrEqual(t, result.Total, 0)

	// Cleanup
	client.Topics.Delete(audience.ID, topic.ID)
}

// ============================================
// SEGMENTS TESTS
// ============================================

func TestSegmentsList(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	listResult, err := client.Segments.List(audience.ID, nil)
	require.NoError(t, err)
	assert.NotNil(t, listResult)
}

func TestSegmentsCreate(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	segmentName := fmt.Sprintf("Test Segment %d", uniqueID())
	rules := []map[string]interface{}{
		{"field": "email", "operator": "contains", "value": "@example.com"},
	}
	segment, err := client.Segments.Create(audience.ID, sevk.CreateSegmentParams{
		Name:     segmentName,
		Rules:    rules,
		Operator: "AND",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, segment.ID)
	assert.Equal(t, segmentName, segment.Name)

	// Cleanup
	client.Segments.Delete(audience.ID, segment.ID)
}

func TestSegmentsGetByID(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	segmentName := fmt.Sprintf("Get Segment %d", uniqueID())
	rules := []map[string]interface{}{
		{"field": "email", "operator": "contains", "value": "@test.com"},
	}
	segment, err := client.Segments.Create(audience.ID, sevk.CreateSegmentParams{
		Name:     segmentName,
		Rules:    rules,
		Operator: "AND",
	})
	require.NoError(t, err)

	fetched, err := client.Segments.Get(audience.ID, segment.ID)
	require.NoError(t, err)
	assert.Equal(t, segment.ID, fetched.ID)

	// Cleanup
	client.Segments.Delete(audience.ID, segment.ID)
}

func TestSegmentsUpdate(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	segmentName := fmt.Sprintf("Update Segment %d", uniqueID())
	rules := []map[string]interface{}{
		{"field": "email", "operator": "contains", "value": "@update.com"},
	}
	segment, err := client.Segments.Create(audience.ID, sevk.CreateSegmentParams{
		Name:     segmentName,
		Rules:    rules,
		Operator: "AND",
	})
	require.NoError(t, err)

	newName := fmt.Sprintf("Updated Segment %d", uniqueID())
	updated, err := client.Segments.Update(audience.ID, segment.ID, sevk.UpdateSegmentParams{Name: &newName})
	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)

	// Cleanup
	client.Segments.Delete(audience.ID, segment.ID)
}

func TestSegmentsCalculate(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	segmentName := fmt.Sprintf("Calculate Segment %d", uniqueID())
	rules := []map[string]interface{}{
		{"field": "email", "operator": "contains", "value": "@example.com"},
	}
	segment, err := client.Segments.Create(audience.ID, sevk.CreateSegmentParams{
		Name:     segmentName,
		Rules:    rules,
		Operator: "AND",
	})
	require.NoError(t, err)

	result, err := client.Segments.Calculate(audience.ID, segment.ID)
	if err != nil {
		// Rate limit (429) is expected if other segment tests ran recently
		sevkErr, ok := err.(*sevk.Error)
		if ok && sevkErr.StatusCode == 429 {
			t.Skip("Skipping due to rate limit on segment calculate endpoint")
		}
		require.NoError(t, err)
	}
	assert.NotNil(t, result)

	// Cleanup
	client.Segments.Delete(audience.ID, segment.ID)
}

func TestSegmentsPreview(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	result, err := client.Segments.Preview(audience.ID, sevk.CreateSegmentParams{
		Name:     "Preview Segment",
		Rules:    []map[string]interface{}{{"field": "email", "operator": "contains", "value": "@example.com"}},
		Operator: "AND",
	})
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestSegmentsDelete(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	segmentName := fmt.Sprintf("Delete Segment %d", uniqueID())
	rules := []map[string]interface{}{
		{"field": "email", "operator": "contains", "value": "@delete.com"},
	}
	segment, err := client.Segments.Create(audience.ID, sevk.CreateSegmentParams{
		Name:     segmentName,
		Rules:    rules,
		Operator: "AND",
	})
	require.NoError(t, err)

	err = client.Segments.Delete(audience.ID, segment.ID)
	require.NoError(t, err)

	_, err = client.Segments.Get(audience.ID, segment.ID)
	assert.Error(t, err)
}

// ============================================
// SUBSCRIPTIONS TESTS
// ============================================

func TestSubscriptionsSubscribe(t *testing.T) {
	client := getClient(t)
	audience := getAudience(t)

	email := fmt.Sprintf("subscribe-test-%d@example.com", uniqueID())
	err := client.Subscriptions.Subscribe(sevk.SubscribeParams{
		Email:      email,
		AudienceID: audience.ID,
	})
	require.NoError(t, err)
}

func TestSubscriptionsUnsubscribe(t *testing.T) {
	client := getClient(t)

	email := fmt.Sprintf("unsubscribe-test-%d@example.com", uniqueID())
	contact, err := client.Contacts.Create(sevk.CreateContactParams{Email: email})
	require.NoError(t, err)

	err = client.Subscriptions.Unsubscribe(sevk.UnsubscribeParams{Email: email})
	require.NoError(t, err)

	updatedContact, err := client.Contacts.Get(contact.ID)
	require.NoError(t, err)
	assert.Equal(t, false, updatedContact.Subscribed)
}

// ============================================
// WEBHOOKS TESTS
// ============================================

func TestWebhooksCRUDCycle(t *testing.T) {
	client := getClient(t)

	// Create
	url := fmt.Sprintf("https://example.com/webhook-%d", uniqueID())
	webhook, err := client.Webhooks.Create(sevk.CreateWebhookParams{
		URL:    url,
		Events: []string{"contact.subscribed"},
	})
	require.NoError(t, err)
	assert.NotEmpty(t, webhook.ID)
	assert.Equal(t, url, webhook.URL)

	// Get
	fetched, err := client.Webhooks.Get(webhook.ID)
	require.NoError(t, err)
	assert.Equal(t, webhook.ID, fetched.ID)

	// Update
	newURL := fmt.Sprintf("https://example.com/webhook-updated-%d", uniqueID())
	updated, err := client.Webhooks.Update(webhook.ID, sevk.UpdateWebhookParams{URL: &newURL})
	require.NoError(t, err)
	assert.Equal(t, webhook.ID, updated.ID)
	assert.Equal(t, newURL, updated.URL)

	// Delete
	err = client.Webhooks.Delete(webhook.ID)
	require.NoError(t, err)

	_, err = client.Webhooks.Get(webhook.ID)
	assert.Error(t, err)
}

func TestWebhooksList(t *testing.T) {
	client := getClient(t)

	result, err := client.Webhooks.List()
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestWebhooksTest(t *testing.T) {
	client := getClient(t)

	url := fmt.Sprintf("https://example.com/webhook-test-%d", uniqueID())
	webhook, err := client.Webhooks.Create(sevk.CreateWebhookParams{
		URL:    url,
		Events: []string{"contact.subscribed"},
	})
	require.NoError(t, err)

	result, err := client.Webhooks.Test(webhook.ID)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Cleanup
	client.Webhooks.Delete(webhook.ID)
}

func TestWebhooksListEvents(t *testing.T) {
	client := getClient(t)

	events, err := client.Webhooks.ListEvents()
	require.NoError(t, err)
	assert.NotNil(t, events)
}

// ============================================
// EVENTS TESTS
// ============================================

func TestEventsListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Events.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
	assert.GreaterOrEqual(t, result.Page, 1)
}

func TestEventsListPagination(t *testing.T) {
	client := getClient(t)

	page := 1
	limit := 5
	result, err := client.Events.List(&sevk.ListEventsParams{Page: &page, Limit: &limit})
	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
}

func TestEventsStats(t *testing.T) {
	client := getClient(t)

	stats, err := client.Events.Stats()
	require.NoError(t, err)
	assert.NotNil(t, stats)
}

// ============================================
// USAGE TESTS
// ============================================

func TestGetUsage(t *testing.T) {
	client := getClient(t)

	usage, err := client.GetUsage()
	require.NoError(t, err)
	assert.NotNil(t, usage)
	assert.GreaterOrEqual(t, usage.AudienceLimit, 0)
	assert.GreaterOrEqual(t, usage.ContactLimit, 0)
	assert.GreaterOrEqual(t, usage.BroadcastLimit, 0)
}

// ============================================
// EMAILS TESTS
// ============================================

func TestEmailsRejectUnverifiedDomain(t *testing.T) {
	client := getClient(t)

	_, err := client.Emails.Send(sevk.SendEmailParams{
		To:      "test@example.com",
		Subject: "Test Email",
		HTML:    "<p>Hello</p>",
		From:    "no-reply@unverified-domain.com",
	})
	assert.Error(t, err)
}

func TestEmailsRejectDomainNotOwned(t *testing.T) {
	client := getClient(t)

	_, err := client.Emails.Send(sevk.SendEmailParams{
		To:      "test@example.com",
		Subject: "Test Email",
		HTML:    "<p>Hello</p>",
		From:    "no-reply@not-owned-domain.com",
	})
	assert.Error(t, err)
}

func TestEmailsRejectInvalidFromAddress(t *testing.T) {
	client := getClient(t)

	_, err := client.Emails.Send(sevk.SendEmailParams{
		To:      "test@example.com",
		Subject: "Test Email",
		HTML:    "<p>Hello</p>",
		From:    "invalid-email-without-domain",
	})
	assert.Error(t, err)
}

func TestEmailsErrorMessageDomainVerification(t *testing.T) {
	client := getClient(t)

	_, err := client.Emails.Send(sevk.SendEmailParams{
		To:      "test@example.com",
		Subject: "Test Email",
		HTML:    "<p>Hello</p>",
		From:    "no-reply@test-domain.com",
	})
	assert.Error(t, err)
}

func TestEmailsGetNonExistent(t *testing.T) {
	client := getClient(t)

	_, err := client.Emails.Get("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
	sevkErr, ok := err.(*sevk.Error)
	if ok {
		assert.True(t, sevkErr.IsNotFound())
	}
}

func TestEmailsBulkRejectUnverifiedDomain(t *testing.T) {
	client := getClient(t)

	result, err := client.Emails.SendBulk([]sevk.SendEmailParams{
		{
			To:      "test1@example.com",
			Subject: "Bulk Test 1",
			HTML:    "<p>Hello 1</p>",
			From:    "no-reply@unverified-domain.com",
		},
		{
			To:      "test2@example.com",
			Subject: "Bulk Test 2",
			HTML:    "<p>Hello 2</p>",
			From:    "no-reply@unverified-domain.com",
		},
	})
	// The API may return an error (e.g., balance check) or return 200 with failures
	if err != nil {
		// Error at the HTTP level is acceptable (e.g., balance check, account not enabled)
		assert.Error(t, err)
	} else {
		// If no HTTP error, the response should show failures for unverified domain
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Failed)
		assert.Equal(t, 0, result.Success)
		assert.NotEmpty(t, result.Errors)
	}
}

// ============================================
// ERROR HANDLING TESTS
// ============================================

func TestError404(t *testing.T) {
	client := getClient(t)

	_, err := client.Contacts.Get("non-existent-id-12345")
	assert.Error(t, err)

	sevkErr, ok := err.(*sevk.Error)
	if ok {
		assert.True(t, sevkErr.IsNotFound())
	}
}

func TestErrorValidation(t *testing.T) {
	client := getClient(t)

	_, err := client.Contacts.Create(sevk.CreateContactParams{Email: "invalid-email"})
	assert.Error(t, err)
}

// ============================================
// MARKUP RENDERER TESTS
// ============================================

func TestRenderDocumentStructure(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "<!DOCTYPE html")
	assert.Contains(t, html, "<html")
	assert.Contains(t, html, "<head>")
	assert.Contains(t, html, "<body")
	assert.Contains(t, html, "</html>")
}

func TestRenderMetaTags(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "charset")
	assert.Contains(t, html, "viewport")
}

func TestRenderTitle(t *testing.T) {
	markupStr := `<email><head><title>Test Email</title></head><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "<title>Test Email</title>")
}

func TestRenderPreviewText(t *testing.T) {
	markupStr := `<email><head><preview>Preview text here</preview></head><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "Preview text here")
	assert.Contains(t, html, "display:none")
}

func TestRenderCustomStyles(t *testing.T) {
	markupStr := `<email><head><style>.custom { color: red; }</style></head><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, ".custom { color: red; }")
}

func TestRenderFontLinks(t *testing.T) {
	markupStr := `<email><head><font name="Roboto" url="https://fonts.googleapis.com/css?family=Roboto" /></head><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "fonts.googleapis.com")
}

func TestRenderEmptyMarkup(t *testing.T) {
	html := sevk.Render("")
	// Empty markup returns empty string - this is expected behavior
	assert.Equal(t, "", html)
}

func TestRenderBodyStyles(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "margin:0")
	assert.Contains(t, html, "padding:0")
	assert.Contains(t, html, "font-family")
}

func TestRenderLangAttribute(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, `lang="en"`)
}

func TestRenderDirAttribute(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, `dir="ltr"`)
}

func TestRenderContentType(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "Content-Type")
	assert.Contains(t, html, "text/html")
}

func TestRenderXHTMLDoctype(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "XHTML 1.0 Transitional")
}

func TestRenderMultipleFonts(t *testing.T) {
	markupStr := `<email><head>
		<font name="Roboto" url="https://fonts.googleapis.com/css?family=Roboto" />
		<font name="Open Sans" url="https://fonts.googleapis.com/css?family=Open+Sans" />
	</head><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "Roboto")
	assert.Contains(t, html, "Open+Sans")
}

func TestRenderNoBackgroundColor(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.NotContains(t, html, "background-color:#ffffff")
}

func TestRenderMailTag(t *testing.T) {
	markupStr := `<mail><body></body></mail>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "<!DOCTYPE html")
	assert.Contains(t, html, "<body")
}

func TestRenderComplexMarkup(t *testing.T) {
	markupStr := `<email>
		<head>
			<title>Complex Email</title>
			<preview>This is a preview</preview>
			<style>.test { color: blue; }</style>
		</head>
		<body></body>
	</email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "Complex Email")
	assert.Contains(t, html, "This is a preview")
	assert.Contains(t, html, ".test { color: blue; }")
}
