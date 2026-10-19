package sevk

import (
	"encoding/json"
	"time"
)

type StringOrList []string

func (value StringOrList) MarshalJSON() ([]byte, error) {
	if len(value) == 1 {
		return json.Marshal(value[0])
	}
	return json.Marshal([]string(value))
}

func (value *StringOrList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*value = StringOrList{single}
		return nil
	}

	var multiple []string
	if err := json.Unmarshal(data, &multiple); err != nil {
		return err
	}
	*value = multiple
	return nil
}

// EmailAttachment represents an email attachment.
type EmailAttachment struct {
	Filename    string  `json:"filename"`
	Content     string  `json:"content"`
	ContentType string  `json:"contentType"`
	CID         *string `json:"cid,omitempty"`
}

// EmailTracking represents per-email tracking overrides
type EmailTracking struct {
	Open  *bool `json:"open,omitempty"`
	Click *bool `json:"click,omitempty"`
}

// SendEmailParams represents parameters for sending an email.
// Markup is rendered to HTML by the SDK locally and stripped before
// the request hits the wire, so it never appears in the JSON payload.
type SendEmailParams struct {
	From        string            `json:"from"`
	To          StringOrList      `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        *string           `json:"text,omitempty"`
	Markup      string            `json:"-"`
	ReplyTo     StringOrList      `json:"reply_to,omitempty"`
	CC          StringOrList      `json:"cc,omitempty"`
	BCC         StringOrList      `json:"bcc,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Tracking    *EmailTracking    `json:"tracking,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	ScheduledAt *string           `json:"scheduledAt,omitempty"`
}

// updateScheduledEmailParams moves an unattempted scheduled email.
type UpdateScheduledEmailParams struct {
	ScheduledAt string `json:"scheduledAt"`
}

// scheduledEmailMutationResponse identifies the mutated email.
type ScheduledEmailMutationResponse struct {
	ID string `json:"id"`
}

type SendBulkEmailParams struct {
	From        string            `json:"from"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        *string           `json:"text,omitempty"`
	Markup      string            `json:"-"`
	ReplyTo     StringOrList      `json:"reply_to,omitempty"`
	CC          StringOrList      `json:"cc,omitempty"`
	BCC         StringOrList      `json:"bcc,omitempty"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	Tracking    *EmailTracking    `json:"tracking,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// Email represents a sent email retrieved from the API
type Email struct {
	ID                 string      `json:"id"`
	MessageID          *string     `json:"messageId,omitempty"`
	ProviderMessageID  *string     `json:"providerMessageId,omitempty"`
	From               *string     `json:"from"`
	To                 *string     `json:"to"`
	Subject            *string     `json:"subject"`
	Body               *string     `json:"body,omitempty"`
	Text               *string     `json:"text,omitempty"`
	Status             EmailStatus `json:"status"`
	Type               EmailType   `json:"type"`
	ContactID          *string     `json:"contactId,omitempty"`
	BroadcastID        *string     `json:"broadcastId,omitempty"`
	BroadcastVariantID *string     `json:"broadcastVariantId,omitempty"`
	ProjectID          string      `json:"projectId"`
	UserAgent          *string     `json:"userAgent,omitempty"`
	OpenedAt           *time.Time  `json:"openedAt,omitempty"`
	ClickedAt          *time.Time  `json:"clickedAt,omitempty"`
	UnsubscribedAt     *time.Time  `json:"unsubscribedAt,omitempty"`
	ScheduledAt        *time.Time  `json:"scheduledAt,omitempty"`
	CancelledAt        *time.Time  `json:"cancelledAt,omitempty"`
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
}

type SendEmailResponse interface {
	sendEmailResponse()
}

type SendEmailIDResponse struct {
	ID string `json:"id"`
}

func (SendEmailIDResponse) sendEmailResponse() {}

type SendEmailIDsResponse struct {
	IDs []string `json:"ids"`
}

func (SendEmailIDsResponse) sendEmailResponse() {}

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
