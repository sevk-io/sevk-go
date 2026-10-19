package sevk

import "time"

// Contact represents a contact in sevk
type Contact struct {
	ID                   string             `json:"id"`
	Email                string             `json:"email"`
	Subscribed           bool               `json:"subscribed"`
	Data                 map[string]any     `json:"data,omitempty"`
	CreatedAt            time.Time          `json:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt"`
	ProjectID            string             `json:"projectId"`
	LastEventAt          *time.Time         `json:"lastEventAt,omitempty"`
	UnsubscribedAt       *time.Time         `json:"unsubscribedAt,omitempty"`
	UnsubscribedBy       *UnsubscribeSource `json:"unsubscribedBy,omitempty"`
	ResubscriptionLocked bool               `json:"resubscriptionLocked"`
	BounceCount          int                `json:"bounceCount"`
	LastBouncedAt        *time.Time         `json:"lastBouncedAt,omitempty"`
}

// CreateContactParams represents parameters for creating a contact
type CreateContactParams struct {
	Email      string         `json:"email"`
	Subscribed *bool          `json:"subscribed,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

// UpdateContactParams represents parameters for updating a contact.
// Data supports explicit null: Null[map[string]any]() clears all metadata.
type UpdateContactParams struct {
	Email      *string                   `json:"email,omitempty"`
	Subscribed *bool                     `json:"subscribed,omitempty"`
	Data       *Nullable[map[string]any] `json:"data,omitempty"`
}

// ListContactsParams represents parameters for listing contacts.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListContactsParams struct {
	Page            *int    `json:"page,omitempty"`
	Limit           *int    `json:"limit,omitempty"`
	After           *string `json:"after,omitempty"`
	Before          *string `json:"before,omitempty"`
	Search          *string `json:"search,omitempty"`
	ExcludeAudience *string `json:"excludeAudience,omitempty"`
	Subscribed      *bool   `json:"subscribed,omitempty"`
}

// BulkUpdateContact represents a single contact update in a bulk operation.
// Data supports explicit null: Null[map[string]any]() clears all metadata.
type BulkUpdateContact struct {
	Email      string                    `json:"email"`
	Subscribed *bool                     `json:"subscribed,omitempty"`
	Data       *Nullable[map[string]any] `json:"data,omitempty"`
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

// ContactActivityItem represents a contact activity entry with joined email/broadcast refs
type ContactActivityItem struct {
	ID          string                `json:"id"`
	Action      ActivityAction        `json:"action"`
	Description *string               `json:"description,omitempty"`
	CreatedAt   time.Time             `json:"createdAt"`
	UserAgent   *string               `json:"userAgent,omitempty"`
	Link        *string               `json:"link,omitempty"`
	Metadata    JSONValue             `json:"metadata,omitempty"`
	Broadcast   *ActivityBroadcastRef `json:"broadcast,omitempty"`
	Email       *ActivityEmailRef     `json:"email,omitempty"`
}

// ActivityBroadcastRef is a slim broadcast reference attached to contact
// activity. Name is nil when the broadcast was deleted.
type ActivityBroadcastRef struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

// ActivityEmailRef is a slim email reference attached to contact activity.
// Subject is nil when the email row has none.
type ActivityEmailRef struct {
	ID      string  `json:"id"`
	Subject *string `json:"subject"`
}

// ImportContact represents a single contact in an import operation
type ImportContact struct {
	Email      string         `json:"email"`
	Subscribed *bool          `json:"subscribed,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

// ImportContactsParams represents parameters for importing contacts
type ImportContactsParams struct {
	Contacts   []ImportContact `json:"contacts"`
	AudienceID *string         `json:"audienceId,omitempty"`
	TopicID    *string         `json:"topicId,omitempty"`
}

// ImportedContact is the slim contact row returned by BulkImport
type ImportedContact struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Subscribed bool   `json:"subscribed"`
}

// ImportContactsResponse represents the response from importing contacts
type ImportContactsResponse struct {
	Created  int               `json:"created"`
	Errors   []string          `json:"errors,omitempty"`
	Contacts []ImportedContact `json:"contacts,omitempty"`
}

// ListContactActivityParams represents parameters for listing contact activity.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListContactActivityParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// AddedContact is the slim contact row returned when adding contacts to an
// audience or topic; only newly linked contacts are included, newest first.
type AddedContact struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Subscribed bool      `json:"subscribed"`
	CreatedAt  time.Time `json:"createdAt"`
}
