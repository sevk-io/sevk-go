package sevk_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	sevk "github.com/sevk-io/sevk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:4000"

var testClient *sevk.Client
var testAudience *sevk.Audience

func uniqueID() int64 {
	return time.Now().UnixMilli() + int64(rand.Intn(10000))
}

func TestMain(m *testing.M) {
	// Setup: Create test environment once before all tests
	httpClient := &http.Client{Timeout: 30 * time.Second}
	unique := uniqueID()

	// 1. Register a new test user
	testEmail := fmt.Sprintf("sdk-test-%d@test.example.com", unique)
	testPassword := "TestPassword123!"

	registerBody, _ := json.Marshal(map[string]string{
		"email":    testEmail,
		"password": testPassword,
	})

	registerRes, err := httpClient.Post(baseURL+"/auth/register", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		fmt.Printf("Skipping tests - setup failed: %v\n", err)
		os.Exit(0)
	}
	defer registerRes.Body.Close()

	if registerRes.StatusCode != http.StatusOK && registerRes.StatusCode != http.StatusCreated {
		fmt.Printf("Skipping tests - registration failed: %d\n", registerRes.StatusCode)
		os.Exit(0)
	}

	var registerData map[string]interface{}
	json.NewDecoder(registerRes.Body).Decode(&registerData)
	token := registerData["token"].(string)

	// 2. Create Project
	projectBody, _ := json.Marshal(map[string]string{
		"name":         "Test Project",
		"slug":         fmt.Sprintf("test-project-%d", unique),
		"supportEmail": "support@test.com",
	})

	projectReq, _ := http.NewRequest("POST", baseURL+"/projects", bytes.NewReader(projectBody))
	projectReq.Header.Set("Authorization", "Bearer "+token)
	projectReq.Header.Set("Content-Type", "application/json")

	projectRes, err := httpClient.Do(projectReq)
	if err != nil {
		fmt.Printf("Skipping tests - project creation failed: %v\n", err)
		os.Exit(0)
	}
	defer projectRes.Body.Close()

	var projectData map[string]interface{}
	json.NewDecoder(projectRes.Body).Decode(&projectData)
	projectID := projectData["project"].(map[string]interface{})["id"].(string)

	// 3. Create API Key
	apiKeyBody, _ := json.Marshal(map[string]interface{}{
		"title":      "Test Key",
		"fullAccess": true,
	})

	apiKeyReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/projects/%s/api-keys", baseURL, projectID), bytes.NewReader(apiKeyBody))
	apiKeyReq.Header.Set("Authorization", "Bearer "+token)
	apiKeyReq.Header.Set("Content-Type", "application/json")

	apiKeyRes, err := httpClient.Do(apiKeyReq)
	if err != nil {
		fmt.Printf("Skipping tests - API key creation failed: %v\n", err)
		os.Exit(0)
	}
	defer apiKeyRes.Body.Close()

	var apiKeyData map[string]interface{}
	json.NewDecoder(apiKeyRes.Body).Decode(&apiKeyData)
	apiKey := apiKeyData["apiKey"].(map[string]interface{})["key"].(string)

	testClient = sevk.NewWithOptions(apiKey, sevk.Options{BaseURL: baseURL})

	// 4. Create a shared test audience for all tests
	audience, err := testClient.Audiences.Create(sevk.CreateAudienceParams{
		Name: fmt.Sprintf("Shared Test Audience %d", unique),
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

// ============================================
// DOMAINS TESTS
// ============================================

func TestDomainsListStructure(t *testing.T) {
	client := getClient(t)

	result, err := client.Domains.List(nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDomainsListVerified(t *testing.T) {
	client := getClient(t)

	verified := true
	result, err := client.Domains.List(&sevk.ListDomainsParams{Verified: &verified})
	require.NoError(t, err)
	assert.NotNil(t, result)
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

func TestRenderBackgroundColor(t *testing.T) {
	markupStr := `<email><body></body></email>`
	html := sevk.Render(markupStr)

	assert.Contains(t, html, "background-color")
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
