package sevk

import "time"

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
}

// Contact represents a contact in Sevk
type Contact struct {
	ID         string                 `json:"id"`
	Email      string                 `json:"email"`
	FirstName  *string                `json:"firstName,omitempty"`
	LastName   *string                `json:"lastName,omitempty"`
	Subscribed bool                   `json:"subscribed"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

// CreateContactParams represents parameters for creating a contact
type CreateContactParams struct {
	Email      string                 `json:"email"`
	FirstName  *string                `json:"firstName,omitempty"`
	LastName   *string                `json:"lastName,omitempty"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateContactParams represents parameters for updating a contact
type UpdateContactParams struct {
	FirstName  *string                `json:"firstName,omitempty"`
	LastName   *string                `json:"lastName,omitempty"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ListContactsParams represents parameters for listing contacts
type ListContactsParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Search *string `json:"search,omitempty"`
}

// Audience represents an audience in Sevk
type Audience struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	UsersCanSee string    `json:"usersCanSee"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateAudienceParams represents parameters for creating an audience
type CreateAudienceParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	UsersCanSee *string `json:"usersCanSee,omitempty"`
}

// UpdateAudienceParams represents parameters for updating an audience
type UpdateAudienceParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	UsersCanSee *string `json:"usersCanSee,omitempty"`
}

// ListAudiencesParams represents parameters for listing audiences
type ListAudiencesParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// Template represents an email template in Sevk
type Template struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateTemplateParams represents parameters for creating a template
type CreateTemplateParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateTemplateParams represents parameters for updating a template
type UpdateTemplateParams struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

// ListTemplatesParams represents parameters for listing templates
type ListTemplatesParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Search *string `json:"search,omitempty"`
}

// Broadcast represents a broadcast in Sevk
type Broadcast struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Subject     string     `json:"subject"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
	SentAt      *time.Time `json:"sentAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ListBroadcastsParams represents parameters for listing broadcasts
type ListBroadcastsParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	Search *string `json:"search,omitempty"`
	Status *string `json:"status,omitempty"`
}

// Domain represents a domain in Sevk
type Domain struct {
	ID         string     `json:"id"`
	Domain     string     `json:"domain"`
	Verified   bool       `json:"verified"`
	VerifiedAt *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// DomainsListResponse represents the response for listing domains
type DomainsListResponse struct {
	Items []Domain `json:"items"`
}

// ListDomainsParams represents parameters for listing domains
type ListDomainsParams struct {
	Verified *bool `json:"verified,omitempty"`
}

// CreateDomainParams represents parameters for creating a domain
type CreateDomainParams struct {
	Domain     string  `json:"domain"`
	Email      string  `json:"email"`
	From       string  `json:"from"`
	SenderName string  `json:"senderName"`
	Region     *string `json:"region,omitempty"`
}

// DnsRecord represents a DNS record for domain verification
type DnsRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	TTL      *int   `json:"ttl,omitempty"`
	Priority *int   `json:"priority,omitempty"`
	Status   string `json:"status,omitempty"`
}

// DnsRecordsResponse represents the response for getting DNS records
type DnsRecordsResponse struct {
	Items []DnsRecord `json:"items"`
}


// Topic represents a topic in Sevk
type Topic struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	AudienceID  string    `json:"audienceId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateTopicParams represents parameters for creating a topic
type CreateTopicParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// UpdateTopicParams represents parameters for updating a topic
type UpdateTopicParams struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// ListTopicsParams represents parameters for listing topics
type ListTopicsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// Segment represents a segment in Sevk
type Segment struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Rules      []map[string]interface{} `json:"rules"`
	Operator   string                   `json:"operator"`
	AudienceID string                   `json:"audienceId"`
	CreatedAt  time.Time                `json:"createdAt"`
	UpdatedAt  time.Time                `json:"updatedAt"`
}

// CreateSegmentParams represents parameters for creating a segment
type CreateSegmentParams struct {
	Name     string                   `json:"name"`
	Rules    []map[string]interface{} `json:"rules"`
	Operator string                   `json:"operator"`
}

// UpdateSegmentParams represents parameters for updating a segment
type UpdateSegmentParams struct {
	Name     *string                  `json:"name,omitempty"`
	Rules    []map[string]interface{} `json:"rules,omitempty"`
	Operator *string                  `json:"operator,omitempty"`
}

// ListSegmentsParams represents parameters for listing segments
type ListSegmentsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// SubscribeParams represents parameters for subscribing
type SubscribeParams struct {
	Email      string   `json:"email"`
	AudienceID string   `json:"audienceId"`
	TopicIDs   []string `json:"topicIds,omitempty"`
}

// UnsubscribeParams represents parameters for unsubscribing
type UnsubscribeParams struct {
	Email string `json:"email"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`     // Base64 encoded
	ContentType string `json:"contentType"` // MIME type
}

// SendEmailParams represents parameters for sending an email
type SendEmailParams struct {
	To          interface{}       `json:"to"` // string or []string
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	From        string            `json:"from,omitempty"`
	FromName    *string           `json:"fromName,omitempty"`
	ReplyTo     *string           `json:"replyTo,omitempty"`
	Text        *string           `json:"text,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
}

// Email represents a sent email retrieved from the API
type Email struct {
	ID        string     `json:"id"`
	To        string     `json:"to"`
	From      string     `json:"from"`
	Subject   string     `json:"subject"`
	HTML      string     `json:"html,omitempty"`
	Text      string     `json:"text,omitempty"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	SentAt    *time.Time `json:"sentAt,omitempty"`
}

// SendEmailResponse represents the response from sending an email
type SendEmailResponse struct {
	ID  string   `json:"id,omitempty"`
	IDs []string `json:"ids,omitempty"`
}

// BulkEmailError represents an error in bulk email sending
type BulkEmailError struct {
	Index int    `json:"index"`
	Email string `json:"email"`
	Error string `json:"error"`
}

// BulkEmailResponse represents the response from sending bulk emails
type BulkEmailResponse struct {
	Success int              `json:"success"`
	Failed  int              `json:"failed"`
	IDs     []string         `json:"ids"`
	Errors  []BulkEmailError `json:"errors,omitempty"`
}

// BulkUpdateContact represents a single contact update in a bulk operation
type BulkUpdateContact struct {
	Email      string                 `json:"email"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// BulkUpdateContactsParams represents parameters for bulk updating contacts
type BulkUpdateContactsParams struct {
	Contacts []BulkUpdateContact `json:"contacts"`
}

// BulkUpdateContactsResponse represents the response from bulk updating contacts
type BulkUpdateContactsResponse struct {
	Updated        int      `json:"updated"`
	NotFound       int      `json:"notFound"`
	Errors         []string `json:"errors,omitempty"`
	UpdatedEmails  []string `json:"updatedEmails,omitempty"`
	NotFoundEmails []string `json:"notFoundEmails,omitempty"`
}

// ImportContact represents a single contact in an import operation
type ImportContact struct {
	Email      string                 `json:"email"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// ImportContactsParams represents parameters for importing contacts
type ImportContactsParams struct {
	Contacts   []ImportContact `json:"contacts"`
	AudienceID *string         `json:"audienceId,omitempty"`
}

// ImportContactsResponse represents the response from importing contacts
type ImportContactsResponse struct {
	Created  int         `json:"created"`
	Errors   []string    `json:"errors,omitempty"`
	Contacts []Contact   `json:"contacts,omitempty"`
}

// ContactEvent represents an event associated with a contact
type ContactEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Action      string    `json:"action"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ListContactEventsParams represents parameters for listing contact events
type ListContactEventsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// ListAudienceContactsParams represents parameters for listing contacts in an audience
type ListAudienceContactsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// RegionsResponse represents the response for listing domain regions
type RegionsResponse struct {
	Items []string `json:"items"`
}

// CreateBroadcastParams represents parameters for creating a broadcast
type CreateBroadcastParams struct {
	Name        string  `json:"name"`
	Subject     string  `json:"subject"`
	Body        string  `json:"body"`
	Style       *string `json:"style,omitempty"`
	TargetType  *string `json:"targetType,omitempty"`
	AudienceID  *string `json:"audienceId,omitempty"`
	TopicID     *string `json:"topicId,omitempty"`
	SegmentID   *string `json:"segmentId,omitempty"`
	SenderName  *string `json:"senderName,omitempty"`
	DomainID    string  `json:"domainId"`
	ScheduledAt *string `json:"scheduledAt,omitempty"`
}

// UpdateBroadcastParams represents parameters for updating a broadcast
type UpdateBroadcastParams struct {
	Name        *string `json:"name,omitempty"`
	Subject     *string `json:"subject,omitempty"`
	Body        *string `json:"body,omitempty"`
	Style       *string `json:"style,omitempty"`
	TargetType  *string `json:"targetType,omitempty"`
	AudienceID  *string `json:"audienceId,omitempty"`
	TopicID     *string `json:"topicId,omitempty"`
	SegmentID   *string `json:"segmentId,omitempty"`
	SenderName  *string `json:"senderName,omitempty"`
	DomainID    *string `json:"domainId,omitempty"`
	ScheduledAt *string `json:"scheduledAt,omitempty"`
}

// SendTestParams represents parameters for sending a test broadcast email
type SendTestParams struct {
	Emails []string `json:"emails"`
}

// BroadcastAnalytics represents analytics data for a broadcast
type BroadcastAnalytics struct {
	Total      int `json:"total"`
	Sent       int `json:"sent"`
	Delivered  int `json:"delivered"`
	Opened     int `json:"opened"`
	Clicked    int `json:"clicked"`
	Bounced    int `json:"bounced"`
	Complained int `json:"complained"`
}

// SegmentCalculation represents the result of a segment calculation
type SegmentCalculation struct {
	Count int       `json:"count"`
	Items []Contact `json:"items,omitempty"`
}

// AddTopicContactsParams represents parameters for adding contacts to a topic
type AddTopicContactsParams struct {
	ContactIDs []string `json:"contactIds"`
}

// UpdateDomainRequest represents parameters for updating a domain
type UpdateDomainRequest struct {
	Email         *string `json:"email,omitempty"`
	From          *string `json:"from,omitempty"`
	SenderName    *string `json:"senderName,omitempty"`
	ClickTracking *bool   `json:"clickTracking,omitempty"`
	OpenTracking  *bool   `json:"openTracking,omitempty"`
}

// BroadcastStatus represents the status of a broadcast
type BroadcastStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Total  int    `json:"total"`
	Sent   int    `json:"sent"`
	Failed int    `json:"failed"`
}

// BroadcastEmail represents an email sent as part of a broadcast
type BroadcastEmail struct {
	ID        string     `json:"id"`
	To        string     `json:"to"`
	Status    string     `json:"status"`
	SentAt    *time.Time `json:"sentAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// ListBroadcastEmailsParams represents parameters for listing broadcast emails
type ListBroadcastEmailsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// ActiveBroadcastsResponse represents the response for listing active broadcasts
type ActiveBroadcastsResponse struct {
	Items []Broadcast `json:"items"`
	Total int         `json:"total"`
}

// BroadcastCostEstimate represents the cost estimate for a broadcast
type BroadcastCostEstimate struct {
	Recipients int     `json:"recipients"`
	Cost       float64 `json:"cost"`
}

// ListTopicContactsParams represents parameters for listing contacts in a topic
type ListTopicContactsParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// Webhook represents a webhook in Sevk
type Webhook struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Active    bool      `json:"active"`
	Secret    *string   `json:"secret,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// WebhooksListResponse represents the response for listing webhooks
type WebhooksListResponse struct {
	Items []Webhook `json:"items"`
}

// CreateWebhookParams represents parameters for creating a webhook
type CreateWebhookParams struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active *bool    `json:"active,omitempty"`
}

// UpdateWebhookParams represents parameters for updating a webhook
type UpdateWebhookParams struct {
	URL    *string  `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Active *bool    `json:"active,omitempty"`
}

// TestWebhookResponse represents the response from testing a webhook
type TestWebhookResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message,omitempty"`
}

// WebhookEventsListResponse represents the response for listing webhook event types
type WebhookEventsListResponse struct {
	Items  []string               `json:"items"`
	Events map[string]interface{} `json:"events,omitempty"`
}

// Event represents an event in Sevk
type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	ContactID   *string                `json:"contactId,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Description *string                `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
}

// ListEventsParams represents parameters for listing events
type ListEventsParams struct {
	Page      *int    `json:"page,omitempty"`
	Limit     *int    `json:"limit,omitempty"`
	Type      *string `json:"type,omitempty"`
	Action    *string `json:"action,omitempty"`
	ContactID *string `json:"contactId,omitempty"`
}

// EventStats represents event statistics
type EventStats struct {
	Total     int            `json:"total"`
	ByType    map[string]int `json:"byType"`
	ByAction  map[string]int `json:"byAction"`
}

// APIError represents an error response from the API
type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
