package sevk

// SubscribeParams represents parameters for subscribing
type SubscribeParams struct {
	Email      string  `json:"email"`
	AudienceID *string `json:"audienceId,omitempty"`
	// Pointer preserves the API distinction between omitted data and an explicit empty object.
	Data *map[string]any `json:"data,omitempty"`
}

// UnsubscribeParams represents parameters for unsubscribing
type UnsubscribeParams struct {
	Email       *string `json:"email,omitempty"`
	ContactID   *string `json:"contactId,omitempty"`
	BroadcastID *string `json:"broadcastId,omitempty"`
}

// SubscribeResponseContact represents the contact returned by Subscribe
type SubscribeResponseContact struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Subscribed bool           `json:"subscribed"`
	Data       map[string]any `json:"data,omitempty"`
}

// SubscribeResponseAudience represents the audience returned by Subscribe
type SubscribeResponseAudience struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SubscribeResponse represents the response from Subscribe
type SubscribeResponse struct {
	Contact           SubscribeResponseContact   `json:"contact"`
	Audience          *SubscribeResponseAudience `json:"audience,omitempty"`
	AlreadySubscribed *bool                      `json:"alreadySubscribed,omitempty"`
	Message           *string                    `json:"message,omitempty"`
}

// UnsubscribeResponseContact represents the contact returned by Unsubscribe
type UnsubscribeResponseContact struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Subscribed bool   `json:"subscribed"`
}

// UnsubscribeResponse represents the response from Unsubscribe
type UnsubscribeResponse struct {
	Contact             UnsubscribeResponseContact `json:"contact"`
	AlreadyUnsubscribed *bool                      `json:"alreadyUnsubscribed,omitempty"`
}
