package sevk

import "time"

const (
	WebhookEventEmailSent            WebhookEventType = "email.sent"
	WebhookEventEmailDelivered       WebhookEventType = "email.delivered"
	WebhookEventEmailOpened          WebhookEventType = "email.opened"
	WebhookEventEmailClicked         WebhookEventType = "email.clicked"
	WebhookEventEmailBounced         WebhookEventType = "email.bounced"
	WebhookEventEmailComplained      WebhookEventType = "email.complained"
	WebhookEventEmailRejected        WebhookEventType = "email.rejected"
	WebhookEventEmailDelayed         WebhookEventType = "email.delayed"
	WebhookEventEmailRenderingFailed WebhookEventType = "email.rendering_failed"
	WebhookEventEmailReceived        WebhookEventType = "email.received"
	WebhookEventContactSubscribed    WebhookEventType = "contact.subscribed"
	WebhookEventContactUnsubscribed  WebhookEventType = "contact.unsubscribed"
	WebhookEventBroadcastStarted     WebhookEventType = "broadcast.started"
	WebhookEventBroadcastFinished    WebhookEventType = "broadcast.finished"
)

// Webhook represents a webhook in sevk
type Webhook struct {
	ID      string             `json:"id"`
	URL     string             `json:"url"`
	Events  []WebhookEventType `json:"events"`
	Enabled bool               `json:"enabled"`
	// Secret is a masked placeholder (••••••••) outside the one-time Create response.
	Secret         *string   `json:"secret,omitempty"`
	FailureCount   int       `json:"failureCount"`
	LastFailureAt  *string   `json:"lastFailureAt,omitempty"`
	LastDeliveryAt *string   `json:"lastDeliveryAt,omitempty"`
	DisabledAt     *string   `json:"disabledAt,omitempty"`
	ProjectID      string    `json:"projectId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// CreatedWebhook is the one-time create response with a usable whsec_ secret.
type CreatedWebhook struct {
	ID             string             `json:"id"`
	URL            string             `json:"url"`
	Events         []WebhookEventType `json:"events"`
	Enabled        bool               `json:"enabled"`
	Secret         string             `json:"secret"`
	FailureCount   int                `json:"failureCount"`
	LastFailureAt  *string            `json:"lastFailureAt"`
	LastDeliveryAt *string            `json:"lastDeliveryAt"`
	DisabledAt     *string            `json:"disabledAt"`
	ProjectID      string             `json:"projectId"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}

// WebhookEventDefinition describes a webhook event type
type WebhookEventDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Category    WebhookCategory `json:"category"`
}

// WebhooksListResponse represents the response for listing webhooks
type WebhooksListResponse struct {
	PaginatedResponse[Webhook]
	AvailableEvents map[string]WebhookEventDefinition `json:"availableEvents,omitempty"`
}

type ListWebhooksParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// CreateWebhookParams represents parameters for creating a webhook
type CreateWebhookParams struct {
	URL     string             `json:"url"`
	Events  []WebhookEventType `json:"events"`
	Enabled *bool              `json:"enabled,omitempty"`
}

// UpdateWebhookParams represents parameters for updating a webhook
type UpdateWebhookParams struct {
	URL     *string             `json:"url,omitempty"`
	Events  *[]WebhookEventType `json:"events,omitempty"`
	Enabled *bool               `json:"enabled,omitempty"`
}

// TestWebhookResponse represents the response from testing a webhook
type TestWebhookResponse struct {
	Success    bool   `json:"success"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
	Message    string `json:"message,omitempty"`
}

type RotateWebhookResponse struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// WebhookEventsListResponse represents the response for listing webhook event types
type WebhookEventsListResponse struct {
	Items  []WebhookEventType                `json:"items"`
	Events map[string]WebhookEventDefinition `json:"events,omitempty"`
}
