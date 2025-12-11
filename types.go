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
	Domains []Domain `json:"domains"`
}

// ListDomainsParams represents parameters for listing domains
type ListDomainsParams struct {
	Verified *bool `json:"verified,omitempty"`
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

// SendEmailParams represents parameters for sending an email
type SendEmailParams struct {
	To       string  `json:"to"`
	Subject  string  `json:"subject"`
	HTML     string  `json:"html"`
	From     string  `json:"from"`
	FromName *string `json:"fromName,omitempty"`
	ReplyTo  *string `json:"replyTo,omitempty"`
	Text     *string `json:"text,omitempty"`
}

// SendEmailResponse represents the response from sending an email
type SendEmailResponse struct {
	MessageID string `json:"messageId"`
}

// APIError represents an error response from the API
type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
