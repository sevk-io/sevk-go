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

// Send sends an email
func (r *EmailsResource) Send(params SendEmailParams) (*SendEmailResponse, error) {
	var result SendEmailResponse
	if err := r.client.post("/emails", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
