package sevk

import "time"

// AudienceBase is returned by create and update operations, which omit counters.
type AudienceBase struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description,omitempty"`
	Visibility  AudienceVisibility `json:"visibility"`
	ProjectID   string             `json:"projectId"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

// Audience represents a read response with aggregate counters.
type Audience struct {
	ID                      string             `json:"id"`
	Name                    string             `json:"name"`
	Description             *string            `json:"description,omitempty"`
	Visibility              AudienceVisibility `json:"visibility"`
	ProjectID               string             `json:"projectId"`
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
	ContactsCount           int64              `json:"contactsCount"`
	SubscribedContactsCount int64              `json:"subscribedContactsCount"`
	BroadcastsCount         int64              `json:"broadcastsCount"`
	TopicsCount             int64              `json:"topicsCount"`
	SegmentsCount           int64              `json:"segmentsCount"`
}

// CreateAudienceParams represents parameters for creating an audience
type CreateAudienceParams struct {
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Visibility  *AudienceVisibility `json:"visibility,omitempty"`
}

// UpdateAudienceParams represents parameters for updating an audience
type UpdateAudienceParams struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Visibility  *AudienceVisibility `json:"visibility,omitempty"`
}

// ListAudiencesParams represents parameters for listing audiences.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListAudiencesParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
	Search *string `json:"search,omitempty"`
}

// ListAudienceContactsParams represents parameters for listing contacts in an audience.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListAudienceContactsParams struct {
	Page           *int    `json:"page,omitempty"`
	Limit          *int    `json:"limit,omitempty"`
	After          *string `json:"after,omitempty"`
	Before         *string `json:"before,omitempty"`
	Search         *string `json:"search,omitempty"`
	Subscribed     *bool   `json:"subscribed,omitempty"`
	ExcludeTopicID *string `json:"excludeTopicId,omitempty"`
}

// AddContactsToAudienceResponse represents the response from adding contacts to an audience
type AddContactsToAudienceResponse struct {
	Added    int            `json:"added"`
	Contacts []AddedContact `json:"contacts"`
}

// AudienceRemoveContactResponse represents the response from removing a contact from an audience
type AudienceRemoveContactResponse struct {
	AudienceID string `json:"audienceId"`
	ContactID  string `json:"contactId"`
	Removed    bool   `json:"removed"`
}
