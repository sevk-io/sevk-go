package sevk

import "time"

// InboundAttachment represents an attachment on an inbound email
type InboundAttachment struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	MimeType  string    `json:"mimeType"`
	ContentID *string   `json:"contentId"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

// InboundReplyEmailRef is the email reference returned when an inbound is linked
// to an outbound email. MessageID, From, Status, CreatedAt, BroadcastID, and
// ContactID are only populated on the Get variant.
type InboundReplyEmailRef struct {
	ID          string       `json:"id"`
	Subject     *string      `json:"subject,omitempty"`
	To          *string      `json:"to,omitempty"`
	MessageID   *string      `json:"messageId,omitempty"`
	From        *string      `json:"from,omitempty"`
	Status      *EmailStatus `json:"status,omitempty"`
	CreatedAt   *time.Time   `json:"createdAt,omitempty"`
	BroadcastID *string      `json:"broadcastId,omitempty"`
	ContactID   *string      `json:"contactId,omitempty"`
}

// InboundReplyBroadcastRef is the broadcast reference returned when an inbound is
// linked to a broadcast. Subject, Status, and CreatedAt are only populated on
// the Get variant.
type InboundReplyBroadcastRef struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	Subject   *string                   `json:"subject,omitempty"`
	Status    *BroadcastLifecycleStatus `json:"status,omitempty"`
	CreatedAt *time.Time                `json:"createdAt,omitempty"`
}

// InboundReplyEmailSummary is the lean reply reference returned by list.
type InboundReplyEmailSummary struct {
	ID      string  `json:"id"`
	Subject *string `json:"subject"`
	To      *string `json:"to"`
}

// InboundReplyBroadcastSummary is the lean broadcast reference returned by list.
type InboundReplyBroadcastSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// InboundEmailSummary is the exact projection returned by list operations.
type InboundEmailSummary struct {
	ID                 string                        `json:"id"`
	MessageID          *string                       `json:"messageId"`
	FromAddress        string                        `json:"fromAddress"`
	ToAddress          string                        `json:"toAddress"`
	Subject            *string                       `json:"subject"`
	Status             InboundEmailStatus            `json:"status"`
	IsSpam             bool                          `json:"isSpam"`
	IsVirus            bool                          `json:"isVirus"`
	SpamScore          *float64                      `json:"spamScore"`
	RawSize            int64                         `json:"rawSize"`
	DomainID           *string                       `json:"domainId"`
	ReplyToEmailID     *string                       `json:"replyToEmailId"`
	ReplyToBroadcastID *string                       `json:"replyToBroadcastId"`
	ReceivedAt         time.Time                     `json:"receivedAt"`
	CreatedAt          time.Time                     `json:"createdAt"`
	ReplyToEmail       *InboundReplyEmailSummary     `json:"replyToEmail"`
	ReplyToBroadcast   *InboundReplyBroadcastSummary `json:"replyToBroadcast"`
}

// InboundEmail represents an inbound email received by the project.
// BodyText, BodyHtml, InReplyTo, References, and Attachments are only
// populated on the Get variant.
type InboundEmail struct {
	ID                 string                    `json:"id"`
	MessageID          *string                   `json:"messageId,omitempty"`
	FromAddress        string                    `json:"fromAddress"`
	ToAddress          string                    `json:"toAddress"`
	CcAddresses        []string                  `json:"ccAddresses,omitempty"`
	Subject            *string                   `json:"subject,omitempty"`
	Status             InboundEmailStatus        `json:"status"`
	IsSpam             bool                      `json:"isSpam"`
	IsVirus            bool                      `json:"isVirus"`
	SpamScore          *float64                  `json:"spamScore,omitempty"`
	RawSize            int64                     `json:"rawSize"`
	DomainID           *string                   `json:"domainId"`
	ReceivedAt         time.Time                 `json:"receivedAt"`
	CreatedAt          time.Time                 `json:"createdAt"`
	InReplyTo          *string                   `json:"inReplyTo,omitempty"`
	References         []string                  `json:"references,omitempty"`
	BodyText           *string                   `json:"bodyText,omitempty"`
	BodyHtml           *string                   `json:"bodyHtml,omitempty"`
	Headers            any                       `json:"headers"`
	ProjectID          string                    `json:"projectId,omitempty"`
	ProcessedAt        *time.Time                `json:"processedAt,omitempty"`
	UpdatedAt          *time.Time                `json:"updatedAt,omitempty"`
	Attachments        []InboundAttachment       `json:"attachments,omitempty"`
	ReplyToEmailID     *string                   `json:"replyToEmailId,omitempty"`
	ReplyToBroadcastID *string                   `json:"replyToBroadcastId,omitempty"`
	ReplyToEmail       *InboundReplyEmailRef     `json:"replyToEmail,omitempty"`
	ReplyToBroadcast   *InboundReplyBroadcastRef `json:"replyToBroadcast,omitempty"`
}

// InboundStats represents per-status counts for inbound mail in the project
type InboundStats struct {
	Total             int                `json:"total"`
	Counts            InboundStatsCounts `json:"counts"`
	Received          int                `json:"received"`
	Processed         int                `json:"processed"`
	Bounced           int                `json:"bounced"`
	Blocked           int                `json:"blocked"`
	BlocklistRejected int                `json:"blocklistRejected"`
	Spam              int                `json:"spam"`
	Failed            int                `json:"failed"`
}

type InboundStatsCounts struct {
	Received  int `json:"RECEIVED"`
	Processed int `json:"PROCESSED"`
	Bounced   int `json:"BOUNCED"`
	Blocked   int `json:"BLOCKED"`
	Failed    int `json:"FAILED"`
}

// InboundRoute represents a forward rule for inbound mail
type InboundRoute struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Name      string    `json:"name"`
	Pattern   string    `json:"pattern"`
	Target    string    `json:"target"`
	Priority  int       `json:"priority"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// InboundBlocklistEntry represents a single blocklist entry
type InboundBlocklistEntry struct {
	ID        string               `json:"id"`
	ProjectID string               `json:"projectId"`
	Type      InboundBlocklistType `json:"type"`
	Value     string               `json:"value"`
	Reason    *string              `json:"reason,omitempty"`
	CreatedAt time.Time            `json:"createdAt"`
}

// InboundQuota represents the inbound quota state for the project
type InboundQuota struct {
	ProjectID                string    `json:"projectId"`
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
	MonthlyEmailLimit        *int      `json:"monthlyEmailLimit"`
	MonthlyEmailsUsed        int       `json:"monthlyEmailsUsed"`
	FreeMonthlyInboundEmails *int      `json:"freeMonthlyInboundEmails"`
	FreeEmailsUsed           int       `json:"freeEmailsUsed"`
	InboundEmailPrice        float64   `json:"inboundEmailPrice"`
	ResetAt                  string    `json:"resetAt"`
	PausedAt                 *string   `json:"pausedAt,omitempty"`
	PausedReason             *string   `json:"pausedReason,omitempty"`
}

// InboundSecurityConfig represents per-domain inbound security configuration.
// ID, CreatedAt, and UpdatedAt are absent when the domain has no saved
// config yet (the defaults are returned).
type InboundSecurityConfig struct {
	ID                 string                  `json:"id,omitempty"`
	DomainID           string                  `json:"domainId"`
	DmarcStrict        bool                    `json:"dmarcStrict"`
	SpfRequired        bool                    `json:"spfRequired"`
	SpamThreshold      *float64                `json:"spamThreshold"`
	CatchAllEnabled    bool                    `json:"catchAllEnabled"`
	AcceptedRecipients []string                `json:"acceptedRecipients"`
	MaxAttachmentSize  int64                   `json:"maxAttachmentSize"`
	AllowedMimeTypes   []string                `json:"allowedMimeTypes"`
	CreatedAt          *time.Time              `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time              `json:"updatedAt,omitempty"`
	Platform           InboundSecurityPlatform `json:"platform"`
}

// InboundSecurityPlatform carries instance-level inbound limits returned
// alongside every security config; MaxAttachmentSizeCap is the hard ceiling
// MaxAttachmentSize cannot exceed.
type InboundSecurityPlatform struct {
	MaxAttachmentSizeCap int64 `json:"maxAttachmentSizeCap"`
}

// CreateInboundRouteRequest represents parameters for creating an inbound route
type CreateInboundRouteRequest struct {
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Target   string `json:"target"`
	Priority *int   `json:"priority,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

// UpdateInboundRouteRequest represents parameters for updating an inbound route
type UpdateInboundRouteRequest struct {
	Name     *string `json:"name,omitempty"`
	Pattern  *string `json:"pattern,omitempty"`
	Target   *string `json:"target,omitempty"`
	Priority *int    `json:"priority,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

// AddInboundBlocklistRequest represents parameters for adding a blocklist entry
type AddInboundBlocklistRequest struct {
	Type   InboundBlocklistType `json:"type"`
	Value  string               `json:"value"`
	Reason *string              `json:"reason,omitempty"`
}

// UpdateInboundQuotaRequest represents parameters for updating the inbound quota.
// MonthlyEmailLimit is tri-state: nil omits the field, NewNullable sets it,
// and Null[int]() clears the cap with an explicit JSON null.
type UpdateInboundQuotaRequest struct {
	MonthlyEmailLimit *Nullable[int] `json:"monthlyEmailLimit,omitempty"`
}

// UpdateInboundSecurityRequest represents parameters for updating the per-domain
// inbound security config. AllowedMimeTypes is a pointer so a nil value omits the
// field while a pointer to an empty slice serializes as [] and clears the
// whitelist server-side.
type UpdateInboundSecurityRequest struct {
	DmarcStrict        *bool              `json:"dmarcStrict,omitempty"`
	SpfRequired        *bool              `json:"spfRequired,omitempty"`
	SpamThreshold      *Nullable[float64] `json:"spamThreshold,omitempty"`
	CatchAllEnabled    *bool              `json:"catchAllEnabled,omitempty"`
	AcceptedRecipients *[]string          `json:"acceptedRecipients,omitempty"`
	MaxAttachmentSize  *int64             `json:"maxAttachmentSize,omitempty"`
	AllowedMimeTypes   *[]string          `json:"allowedMimeTypes,omitempty"`
}

// ListInboundRequest represents query parameters for listing inbound mail.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListInboundRequest struct {
	Page        *int                `json:"page,omitempty"`
	Limit       *int                `json:"limit,omitempty"`
	After       *string             `json:"after,omitempty"`
	Before      *string             `json:"before,omitempty"`
	Search      *string             `json:"search,omitempty"`
	Status      *InboundEmailStatus `json:"status,omitempty"`
	DomainID    *string             `json:"domainId,omitempty"`
	BroadcastID *string             `json:"broadcastId,omitempty"`
	IsSpam      *bool               `json:"isSpam,omitempty"`
	IsVirus     *bool               `json:"isVirus,omitempty"`
}

// InboundListResponse represents the paginated list response for inbound mail
type InboundListResponse struct {
	Items      []InboundEmailSummary `json:"items"`
	Total      *int                  `json:"total,omitempty"`
	Page       *int                  `json:"page,omitempty"`
	TotalPages *int                  `json:"totalPages,omitempty"`
	HasMore    *bool                 `json:"hasMore,omitempty"`
	NextCursor *string               `json:"nextCursor,omitempty"`
	PrevCursor *string               `json:"prevCursor,omitempty"`
}

// InboundAttachmentDownload represents a short-lived download URL for an inbound attachment
type InboundAttachmentDownload struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
}

// InboundRouteListResponse preserves the route list response envelope.
type InboundRouteListResponse struct {
	Items      []InboundRoute `json:"items"`
	Total      *int           `json:"total,omitempty"`
	Page       *int           `json:"page,omitempty"`
	TotalPages *int           `json:"totalPages,omitempty"`
	HasMore    *bool          `json:"hasMore,omitempty"`
	NextCursor *string        `json:"nextCursor,omitempty"`
	PrevCursor *string        `json:"prevCursor,omitempty"`
}

// InboundBlocklistListResponse preserves the blocklist response envelope.
type InboundBlocklistListResponse struct {
	Items      []InboundBlocklistEntry `json:"items"`
	Total      *int                    `json:"total,omitempty"`
	Page       *int                    `json:"page,omitempty"`
	TotalPages *int                    `json:"totalPages,omitempty"`
	HasMore    *bool                   `json:"hasMore,omitempty"`
	NextCursor *string                 `json:"nextCursor,omitempty"`
	PrevCursor *string                 `json:"prevCursor,omitempty"`
}
