package sevk

import "fmt"

// ContactsResource handles contact operations
type ContactsResource struct {
	client *Client
}

// List returns a paginated list of contacts
func (r *ContactsResource) List(params *ListContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/contacts"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if params.Search != nil {
			query += fmt.Sprintf("search=%s&", *params.Search)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Contact]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a contact by ID
func (r *ContactsResource) Get(id string) (*Contact, error) {
	var result Contact
	if err := r.client.get("/contacts/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new contact
func (r *ContactsResource) Create(params CreateContactParams) (*Contact, error) {
	var result Contact
	if err := r.client.post("/contacts", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a contact
func (r *ContactsResource) Update(id string, params UpdateContactParams) (*Contact, error) {
	var result Contact
	if err := r.client.put("/contacts/"+id, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a contact
func (r *ContactsResource) Delete(id string) error {
	return r.client.delete("/contacts/" + id)
}

// BulkUpdate updates multiple contacts at once
func (r *ContactsResource) BulkUpdate(params BulkUpdateContactsParams) (*BulkUpdateContactsResponse, error) {
	var result BulkUpdateContactsResponse
	if err := r.client.put("/contacts/bulk-update", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Import imports contacts from CSV content
func (r *ContactsResource) Import(params ImportContactsParams) (*ImportContactsResponse, error) {
	var result ImportContactsResponse
	if err := r.client.post("/contacts/import", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEvents returns events for a contact
func (r *ContactsResource) GetEvents(id string, params *ListContactEventsParams) (*PaginatedResponse[ContactEvent], error) {
	path := "/contacts/" + id + "/events"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[ContactEvent]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AudiencesResource handles audience operations
type AudiencesResource struct {
	client *Client
}

// List returns a paginated list of audiences
func (r *AudiencesResource) List(params *ListAudiencesParams) (*PaginatedResponse[Audience], error) {
	path := "/audiences"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Audience]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns an audience by ID
func (r *AudiencesResource) Get(id string) (*Audience, error) {
	var result Audience
	if err := r.client.get("/audiences/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new audience
func (r *AudiencesResource) Create(params CreateAudienceParams) (*Audience, error) {
	var result Audience
	if err := r.client.post("/audiences", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an audience
func (r *AudiencesResource) Update(id string, params UpdateAudienceParams) (*Audience, error) {
	var result Audience
	if err := r.client.put("/audiences/"+id, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes an audience
func (r *AudiencesResource) Delete(id string) error {
	return r.client.delete("/audiences/" + id)
}

// AddContacts adds contacts to an audience
func (r *AudiencesResource) AddContacts(audienceID string, contactIDs []string) error {
	body := map[string]interface{}{
		"contactIds": contactIDs,
	}
	return r.client.post("/audiences/"+audienceID+"/contacts", body, nil)
}

// ListContacts returns a paginated list of contacts in an audience
func (r *AudiencesResource) ListContacts(audienceID string, params *ListAudienceContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/audiences/" + audienceID + "/contacts"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Contact]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RemoveContact removes a contact from an audience
func (r *AudiencesResource) RemoveContact(audienceID, contactID string) error {
	return r.client.delete("/audiences/" + audienceID + "/contacts/" + contactID)
}

// TemplatesResource handles template operations
type TemplatesResource struct {
	client *Client
}

// List returns a paginated list of templates
func (r *TemplatesResource) List(params *ListTemplatesParams) (*PaginatedResponse[Template], error) {
	path := "/templates"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if params.Search != nil {
			query += fmt.Sprintf("search=%s&", *params.Search)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Template]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a template by ID
func (r *TemplatesResource) Get(id string) (*Template, error) {
	var result Template
	if err := r.client.get("/templates/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new template
func (r *TemplatesResource) Create(params CreateTemplateParams) (*Template, error) {
	var result Template
	if err := r.client.post("/templates", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a template
func (r *TemplatesResource) Update(id string, params UpdateTemplateParams) (*Template, error) {
	var result Template
	if err := r.client.put("/templates/"+id, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a template
func (r *TemplatesResource) Delete(id string) error {
	return r.client.delete("/templates/" + id)
}

// Duplicate duplicates a template
func (r *TemplatesResource) Duplicate(id string) (*Template, error) {
	var result Template
	if err := r.client.post("/templates/"+id+"/duplicate", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// BroadcastsResource handles broadcast operations
type BroadcastsResource struct {
	client *Client
}

// List returns a paginated list of broadcasts
func (r *BroadcastsResource) List(params *ListBroadcastsParams) (*PaginatedResponse[Broadcast], error) {
	path := "/broadcasts"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if params.Search != nil {
			query += fmt.Sprintf("search=%s&", *params.Search)
		}
		if params.Status != nil {
			query += fmt.Sprintf("status=%s&", *params.Status)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Broadcast]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a broadcast by ID
func (r *BroadcastsResource) Get(id string) (*Broadcast, error) {
	var result Broadcast
	if err := r.client.get("/broadcasts/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new broadcast
func (r *BroadcastsResource) Create(params CreateBroadcastParams) (*Broadcast, error) {
	var result Broadcast
	if err := r.client.post("/broadcasts", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a broadcast
func (r *BroadcastsResource) Update(id string, params UpdateBroadcastParams) (*Broadcast, error) {
	var result Broadcast
	if err := r.client.put("/broadcasts/"+id, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a broadcast
func (r *BroadcastsResource) Delete(id string) error {
	return r.client.delete("/broadcasts/" + id)
}

// Send sends a broadcast
func (r *BroadcastsResource) Send(id string) (*Broadcast, error) {
	var result Broadcast
	if err := r.client.post("/broadcasts/"+id+"/send", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Cancel cancels a broadcast
func (r *BroadcastsResource) Cancel(id string) (*Broadcast, error) {
	var result Broadcast
	if err := r.client.post("/broadcasts/"+id+"/cancel", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendTest sends a test email for a broadcast
func (r *BroadcastsResource) SendTest(id string, params SendTestParams) error {
	return r.client.post("/broadcasts/"+id+"/test", params, nil)
}

// GetAnalytics returns analytics for a broadcast
func (r *BroadcastsResource) GetAnalytics(id string) (*BroadcastAnalytics, error) {
	var result BroadcastAnalytics
	if err := r.client.get("/broadcasts/"+id+"/analytics", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetStatus returns the status of a broadcast
func (r *BroadcastsResource) GetStatus(id string) (*BroadcastStatus, error) {
	var result BroadcastStatus
	if err := r.client.get("/broadcasts/"+id+"/status", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEmails returns a paginated list of emails sent as part of a broadcast
func (r *BroadcastsResource) GetEmails(id string, params *ListBroadcastEmailsParams) (*PaginatedResponse[BroadcastEmail], error) {
	path := "/broadcasts/" + id + "/emails"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[BroadcastEmail]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EstimateCost returns the estimated cost for sending a broadcast
func (r *BroadcastsResource) EstimateCost(id string) (*BroadcastCostEstimate, error) {
	var result BroadcastCostEstimate
	if err := r.client.get("/broadcasts/"+id+"/estimate-cost", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListActive returns a list of active broadcasts
func (r *BroadcastsResource) ListActive() (*ActiveBroadcastsResponse, error) {
	var result ActiveBroadcastsResponse
	if err := r.client.get("/broadcasts/active", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DomainsResource handles domain operations
type DomainsResource struct {
	client *Client
}

// List returns a list of domains
func (r *DomainsResource) List(params *ListDomainsParams) (*DomainsListResponse, error) {
	path := "/domains"
	if params != nil && params.Verified != nil {
		path += fmt.Sprintf("?verified=%t", *params.Verified)
	}

	var result DomainsListResponse
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a domain by ID
func (r *DomainsResource) Get(id string) (*Domain, error) {
	var result Domain
	if err := r.client.get("/domains/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new domain
func (r *DomainsResource) Create(params CreateDomainParams) (*Domain, error) {
	var result Domain
	if err := r.client.post("/domains", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a domain
func (r *DomainsResource) Delete(id string) error {
	return r.client.delete("/domains/" + id)
}

// Update updates a domain
func (r *DomainsResource) Update(id string, data UpdateDomainRequest) (*Domain, error) {
	var result Domain
	if err := r.client.put("/domains/"+id, data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Verify triggers verification for a domain
func (r *DomainsResource) Verify(id string) (*Domain, error) {
	var result Domain
	if err := r.client.post("/domains/"+id+"/verify", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetDnsRecords returns the DNS records for a domain
func (r *DomainsResource) GetDnsRecords(id string) (*DnsRecordsResponse, error) {
	var result DnsRecordsResponse
	if err := r.client.get("/domains/"+id+"/dns-records", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRegions returns the available domain regions
func (r *DomainsResource) GetRegions() (*RegionsResponse, error) {
	var result RegionsResponse
	if err := r.client.get("/domains/regions", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// TopicsResource handles topic operations
type TopicsResource struct {
	client *Client
}

// List returns a paginated list of topics for an audience
func (r *TopicsResource) List(audienceID string, params *ListTopicsParams) (*PaginatedResponse[Topic], error) {
	path := "/audiences/" + audienceID + "/topics"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Topic]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a topic by ID
func (r *TopicsResource) Get(audienceID, topicID string) (*Topic, error) {
	var result Topic
	if err := r.client.get("/audiences/"+audienceID+"/topics/"+topicID, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new topic
func (r *TopicsResource) Create(audienceID string, params CreateTopicParams) (*Topic, error) {
	var result Topic
	if err := r.client.post("/audiences/"+audienceID+"/topics", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a topic
func (r *TopicsResource) Update(audienceID, topicID string, params UpdateTopicParams) (*Topic, error) {
	var result Topic
	if err := r.client.put("/audiences/"+audienceID+"/topics/"+topicID, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a topic
func (r *TopicsResource) Delete(audienceID, topicID string) error {
	return r.client.delete("/audiences/" + audienceID + "/topics/" + topicID)
}

// AddContacts adds contacts to a topic
func (r *TopicsResource) AddContacts(audienceID, topicID string, params AddTopicContactsParams) error {
	return r.client.post("/audiences/"+audienceID+"/topics/"+topicID+"/contacts", params, nil)
}

// RemoveContact removes a contact from a topic
func (r *TopicsResource) RemoveContact(audienceID, topicID, contactID string) error {
	return r.client.delete("/audiences/" + audienceID + "/topics/" + topicID + "/contacts/" + contactID)
}

// ListContacts returns a paginated list of contacts in a topic
func (r *TopicsResource) ListContacts(audienceID, topicID string, params *ListTopicContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/audiences/" + audienceID + "/topics/" + topicID + "/contacts"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Contact]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SegmentsResource handles segment operations
type SegmentsResource struct {
	client *Client
}

// List returns a paginated list of segments for an audience
func (r *SegmentsResource) List(audienceID string, params *ListSegmentsParams) (*PaginatedResponse[Segment], error) {
	path := "/audiences/" + audienceID + "/segments"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Segment]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a segment by ID
func (r *SegmentsResource) Get(audienceID, segmentID string) (*Segment, error) {
	var result Segment
	if err := r.client.get("/audiences/"+audienceID+"/segments/"+segmentID, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new segment
func (r *SegmentsResource) Create(audienceID string, params CreateSegmentParams) (*Segment, error) {
	var result Segment
	if err := r.client.post("/audiences/"+audienceID+"/segments", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a segment
func (r *SegmentsResource) Update(audienceID, segmentID string, params UpdateSegmentParams) (*Segment, error) {
	var result Segment
	if err := r.client.put("/audiences/"+audienceID+"/segments/"+segmentID, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a segment
func (r *SegmentsResource) Delete(audienceID, segmentID string) error {
	return r.client.delete("/audiences/" + audienceID + "/segments/" + segmentID)
}

// Calculate calculates a segment's matching contacts
func (r *SegmentsResource) Calculate(audienceID, segmentID string) (*SegmentCalculation, error) {
	var result SegmentCalculation
	if err := r.client.get("/audiences/"+audienceID+"/segments/"+segmentID+"/calculate", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Preview previews how many contacts match given segment rules before saving
func (r *SegmentsResource) Preview(audienceID string, params CreateSegmentParams) (*SegmentCalculation, error) {
	var result SegmentCalculation
	if err := r.client.post("/audiences/"+audienceID+"/segments/preview", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SubscriptionsResource handles subscription operations
type SubscriptionsResource struct {
	client *Client
}

// Subscribe subscribes an email to an audience
func (r *SubscriptionsResource) Subscribe(params SubscribeParams) error {
	return r.client.post("/subscriptions/subscribe", params, nil)
}

// Unsubscribe unsubscribes an email
func (r *SubscriptionsResource) Unsubscribe(params UnsubscribeParams) error {
	return r.client.post("/subscriptions/unsubscribe", params, nil)
}

// EmailsResource handles email operations
type EmailsResource struct {
	client *Client
}

// Get returns an email by ID
func (r *EmailsResource) Get(id string) (*Email, error) {
	var result Email
	if err := r.client.get("/emails/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Send sends an email
func (r *EmailsResource) Send(params SendEmailParams) (*SendEmailResponse, error) {
	var result SendEmailResponse
	if err := r.client.post("/emails", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendBulk sends multiple emails in bulk (max 100)
func (r *EmailsResource) SendBulk(emails []SendEmailParams) (*BulkEmailResponse, error) {
	body := map[string]interface{}{
		"emails": emails,
	}
	var result BulkEmailResponse
	if err := r.client.post("/emails/bulk", body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// WebhooksResource handles webhook operations
type WebhooksResource struct {
	client *Client
}

// List returns a list of webhooks
func (r *WebhooksResource) List() (*WebhooksListResponse, error) {
	var result WebhooksListResponse
	if err := r.client.get("/webhooks", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get returns a webhook by ID
func (r *WebhooksResource) Get(id string) (*Webhook, error) {
	var result Webhook
	if err := r.client.get("/webhooks/"+id, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new webhook
func (r *WebhooksResource) Create(params CreateWebhookParams) (*Webhook, error) {
	var result Webhook
	if err := r.client.post("/webhooks", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a webhook
func (r *WebhooksResource) Update(id string, params UpdateWebhookParams) (*Webhook, error) {
	var result Webhook
	if err := r.client.put("/webhooks/"+id, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a webhook
func (r *WebhooksResource) Delete(id string) error {
	return r.client.delete("/webhooks/" + id)
}

// Test sends a test event to a webhook
func (r *WebhooksResource) Test(id string) (*TestWebhookResponse, error) {
	var result TestWebhookResponse
	if err := r.client.post("/webhooks/"+id+"/test", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListEvents returns the available webhook event types
func (r *WebhooksResource) ListEvents() (*WebhookEventsListResponse, error) {
	var result WebhookEventsListResponse
	if err := r.client.get("/webhooks/events", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EventsResource handles event operations
type EventsResource struct {
	client *Client
}

// List returns a paginated list of events
func (r *EventsResource) List(params *ListEventsParams) (*PaginatedResponse[Event], error) {
	path := "/events"
	if params != nil {
		query := ""
		if params.Page != nil {
			query += fmt.Sprintf("page=%d&", *params.Page)
		}
		if params.Limit != nil {
			query += fmt.Sprintf("limit=%d&", *params.Limit)
		}
		if params.Type != nil {
			query += fmt.Sprintf("type=%s&", *params.Type)
		}
		if params.Action != nil {
			query += fmt.Sprintf("action=%s&", *params.Action)
		}
		if params.ContactID != nil {
			query += fmt.Sprintf("contactId=%s&", *params.ContactID)
		}
		if query != "" {
			path += "?" + query[:len(query)-1]
		}
	}

	var result PaginatedResponse[Event]
	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stats returns event statistics
func (r *EventsResource) Stats() (*EventStats, error) {
	var result EventStats
	if err := r.client.get("/events/stats", &result); err != nil {
		return nil, err
	}
	return &result, nil
}
