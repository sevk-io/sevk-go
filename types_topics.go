package sevk

import "time"

// TopicBase is returned by create and update operations, which omit counters.
type TopicBase struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	AudienceID  string    `json:"audienceId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Topic represents a read response with aggregate counters.
type Topic struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Description             *string   `json:"description,omitempty"`
	AudienceID              string    `json:"audienceId"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
	ContactsCount           int64     `json:"contactsCount"`
	SubscribedContactsCount int64     `json:"subscribedContactsCount"`
	BroadcastsCount         int64     `json:"broadcastsCount"`
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

// ListTopicsParams represents parameters for listing topics.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListTopicsParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
	Search *string `json:"search,omitempty"`
}

// AddTopicContactsParams represents parameters for adding contacts to a topic
type AddTopicContactsParams struct {
	ContactIDs []string `json:"contactIds"`
}

// ListTopicContactsParams represents parameters for listing contacts in a topic.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListTopicContactsParams struct {
	Page       *int    `json:"page,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
	After      *string `json:"after,omitempty"`
	Before     *string `json:"before,omitempty"`
	Search     *string `json:"search,omitempty"`
	Subscribed *bool   `json:"subscribed,omitempty"`
}

// AddContactsToTopicResponse represents the response from adding contacts to a topic
type AddContactsToTopicResponse struct {
	Added    int            `json:"added"`
	Contacts []AddedContact `json:"contacts"`
}

// TopicRemoveContactResponse represents the response from removing a contact from a topic
type TopicRemoveContactResponse struct {
	TopicID   string `json:"topicId"`
	ContactID string `json:"contactId"`
	Removed   bool   `json:"removed"`
}
