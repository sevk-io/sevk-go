package sevk

import "time"

// BroadcastBase is the core projection returned by every broadcast endpoint.
type BroadcastBase struct {
	ID                           string                   `json:"id"`
	Name                         string                   `json:"name"`
	Subject                      string                   `json:"subject"`
	Body                         string                   `json:"body"`
	Style                        BroadcastStyle           `json:"style"`
	Status                       BroadcastLifecycleStatus `json:"status"`
	TargetType                   BroadcastTargetType      `json:"targetType"`
	AudienceID                   *string                  `json:"audienceId,omitempty"`
	TopicID                      *string                  `json:"topicId,omitempty"`
	SegmentID                    *string                  `json:"segmentId,omitempty"`
	Email                        *string                  `json:"email,omitempty"`
	SenderName                   *string                  `json:"senderName,omitempty"`
	ReplyTo                      *string                  `json:"replyTo,omitempty"`
	ReplyToInbound               *bool                    `json:"replyToInbound,omitempty"`
	DomainID                     *string                  `json:"domainId,omitempty"`
	ScheduledAt                  *time.Time               `json:"scheduledAt,omitempty"`
	SentAt                       *time.Time               `json:"sentAt,omitempty"`
	CompletedAt                  *time.Time               `json:"completedAt,omitempty"`
	Delivered                    *string                  `json:"delivered,omitempty"`
	IncludeListUnsubscribeHeader *bool                    `json:"includeListUnsubscribeHeader,omitempty"`
	OpenTracking                 *bool                    `json:"openTracking,omitempty"`
	ClickTracking                *bool                    `json:"clickTracking,omitempty"`
	ProjectID                    string                   `json:"projectId"`
	CreatedAt                    time.Time                `json:"createdAt"`
	UpdatedAt                    time.Time                `json:"updatedAt"`
	Audience                     *BroadcastAudienceRef    `json:"audience"`
	Topic                        *BroadcastTopicRef       `json:"topic"`
	Segment                      *BroadcastSegmentRef     `json:"segment"`
	Domain                       *BroadcastDomainRef      `json:"domain"`
}

// BroadcastSummary is returned by list/send/cancel and contains delivery counters.
type BroadcastSummary struct {
	BroadcastBase
	EstimatedCost     *string `json:"estimatedCost"`
	ActualCost        *string `json:"actualCost"`
	TotalContacts     int64   `json:"totalContacts"`
	SentCount         int64   `json:"sentCount"`
	FailedCount       int64   `json:"failedCount"`
	DeliveredCount    int64   `json:"deliveredCount"`
	OpenedCount       int64   `json:"openedCount"`
	ClickedCount      int64   `json:"clickedCount"`
	BouncedCount      int64   `json:"bouncedCount"`
	ComplainedCount   int64   `json:"complainedCount"`
	UnsubscribedCount int64   `json:"unsubscribedCount"`
}

// BroadcastWithVariants is returned by create and includes experiment state.
type BroadcastWithVariants struct {
	BroadcastBase
	Variants []BroadcastVariant `json:"variants"`
	AbTest   *BroadcastAbTest   `json:"abTest"`
}

// BroadcastUpdate preserves whether the update response omitted experiment state.
type BroadcastUpdate struct {
	BroadcastBase
	Variants *[]BroadcastVariant `json:"variants,omitempty"`
	AbTest   *BroadcastAbTest    `json:"abTest,omitempty"`
}

// Broadcast is the full get projection with counters and experiment state.
type Broadcast struct {
	BroadcastSummary
	Variants []BroadcastVariant `json:"variants"`
	AbTest   *BroadcastAbTest   `json:"abTest"`
}

// BroadcastAudienceRef is a snapshot of the target audience attached to a broadcast
type BroadcastAudienceRef struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ContactsCount           int    `json:"contactsCount"`
	SubscribedContactsCount int    `json:"subscribedContactsCount"`
	BroadcastsCount         int    `json:"broadcastsCount"`
	TopicsCount             int    `json:"topicsCount"`
	SegmentsCount           int    `json:"segmentsCount"`
}

// BroadcastTopicRef is a snapshot of the target topic attached to a broadcast
type BroadcastTopicRef struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ContactsCount           int    `json:"contactsCount"`
	SubscribedContactsCount int    `json:"subscribedContactsCount"`
	BroadcastsCount         int    `json:"broadcastsCount"`
}

// BroadcastSegmentRef is a snapshot of the target segment attached to a broadcast
type BroadcastSegmentRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BroadcastDomainRef is a snapshot of the sending domain attached to a broadcast
type BroadcastDomainRef struct {
	ID         string  `json:"id"`
	Domain     string  `json:"domain"`
	SenderName *string `json:"senderName,omitempty"`
	Verified   bool    `json:"verified"`
}

// BroadcastVariant represents one A/B test variant on a broadcast.
type BroadcastVariant struct {
	ID                string     `json:"id"`
	Label             string     `json:"label"`
	Subject           *string    `json:"subject,omitempty"`
	Body              *string    `json:"body,omitempty"`
	FromName          *string    `json:"fromName,omitempty"`
	TotalContacts     int64      `json:"totalContacts"`
	SentCount         int64      `json:"sentCount"`
	FailedCount       int64      `json:"failedCount"`
	DeliveredCount    int64      `json:"deliveredCount"`
	OpenedCount       int64      `json:"openedCount"`
	ClickedCount      int64      `json:"clickedCount"`
	BouncedCount      int64      `json:"bouncedCount"`
	ComplainedCount   int64      `json:"complainedCount"`
	UnsubscribedCount int64      `json:"unsubscribedCount"`
	IsWinner          bool       `json:"isWinner"`
	BroadcastID       string     `json:"broadcastId,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

// BroadcastAbTest represents the A/B test configuration and runtime state.
type BroadcastAbTest struct {
	ID               string                `json:"id"`
	TestPercent      int                   `json:"testPercent"`
	EvalAfterMinutes int                   `json:"evalAfterMinutes"`
	WinnerMetric     BroadcastWinnerMetric `json:"winnerMetric"`
	Status           BroadcastABTestStatus `json:"status"`
	ScheduledEvalAt  *time.Time            `json:"scheduledEvalAt,omitempty"`
	WinnerVariantID  *string               `json:"winnerVariantId,omitempty"`
	EvaluatedAt      *time.Time            `json:"evaluatedAt,omitempty"`
	BroadcastID      string                `json:"broadcastId,omitempty"`
	CreatedAt        *time.Time            `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time            `json:"updatedAt,omitempty"`
}

// VariantInput is the per-variant override sent on create/update.
type VariantInput struct {
	Label    string  `json:"label"`
	Subject  *string `json:"subject,omitempty"`
	Body     *string `json:"body,omitempty"`
	FromName *string `json:"fromName,omitempty"`
}

// AbTestConfigInput is the A/B test config sent on create/update.
type AbTestConfigInput struct {
	TestPercent      int                    `json:"testPercent"`
	EvalAfterMinutes int                    `json:"evalAfterMinutes"`
	WinnerMetric     *BroadcastWinnerMetric `json:"winnerMetric,omitempty"`
}

// ListBroadcastsParams represents parameters for listing broadcasts.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListBroadcastsParams struct {
	Page   *int                      `json:"page,omitempty"`
	Limit  *int                      `json:"limit,omitempty"`
	After  *string                   `json:"after,omitempty"`
	Before *string                   `json:"before,omitempty"`
	Search *string                   `json:"search,omitempty"`
	Status *BroadcastLifecycleStatus `json:"status,omitempty"`
}

// CreateBroadcastParams represents parameters for creating a broadcast.
// Body is optional; the API defaults it to an empty string.
type CreateBroadcastParams struct {
	Name                         string               `json:"name"`
	Subject                      string               `json:"subject"`
	Body                         string               `json:"body,omitempty"`
	Style                        *BroadcastStyle      `json:"style,omitempty"`
	TargetType                   *BroadcastTargetType `json:"targetType,omitempty"`
	AudienceID                   *string              `json:"audienceId,omitempty"`
	TopicID                      *string              `json:"topicId,omitempty"`
	SegmentID                    *string              `json:"segmentId,omitempty"`
	SenderName                   string               `json:"senderName"`
	Email                        *string              `json:"email,omitempty"`
	ReplyTo                      *string              `json:"replyTo,omitempty"`
	ReplyToInbound               *bool                `json:"replyToInbound,omitempty"`
	DomainID                     string               `json:"domainId"`
	ScheduledAt                  *string              `json:"scheduledAt,omitempty"`
	IncludeListUnsubscribeHeader *bool                `json:"includeListUnsubscribeHeader,omitempty"`
	OpenTracking                 *bool                `json:"openTracking,omitempty"`
	ClickTracking                *bool                `json:"clickTracking,omitempty"`
	Variants                     []VariantInput       `json:"variants,omitempty"`
	AbTest                       *AbTestConfigInput   `json:"abTest,omitempty"`
}

// UpdateBroadcastParams represents parameters for updating a broadcast.
// Nullable fields support explicit JSON null. Null[string]() clears
// ReplyTo/ScheduledAt. Any target-ID field, including an explicit null,
// requires TargetType. Null[bool]() resets OpenTracking/ClickTracking to the
// domain default. Variants and AbTest can be updated independently, but
// clearing Variants cannot be combined with a non-null AbTest configuration.
type UpdateBroadcastParams struct {
	Name                         *string                      `json:"name,omitempty"`
	Subject                      *string                      `json:"subject,omitempty"`
	Body                         *string                      `json:"body,omitempty"`
	Style                        *BroadcastStyle              `json:"style,omitempty"`
	TargetType                   *BroadcastTargetType         `json:"targetType,omitempty"`
	AudienceID                   *Nullable[string]            `json:"audienceId,omitempty"`
	TopicID                      *Nullable[string]            `json:"topicId,omitempty"`
	SegmentID                    *Nullable[string]            `json:"segmentId,omitempty"`
	SenderName                   *string                      `json:"senderName,omitempty"`
	Email                        *string                      `json:"email,omitempty"`
	ReplyTo                      *Nullable[string]            `json:"replyTo,omitempty"`
	ReplyToInbound               *bool                        `json:"replyToInbound,omitempty"`
	DomainID                     *string                      `json:"domainId,omitempty"`
	ScheduledAt                  *Nullable[string]            `json:"scheduledAt,omitempty"`
	IncludeListUnsubscribeHeader *bool                        `json:"includeListUnsubscribeHeader,omitempty"`
	OpenTracking                 *Nullable[bool]              `json:"openTracking,omitempty"`
	ClickTracking                *Nullable[bool]              `json:"clickTracking,omitempty"`
	Variants                     *Nullable[[]VariantInput]    `json:"variants,omitempty"`
	AbTest                       *Nullable[AbTestConfigInput] `json:"abTest,omitempty"`
}

// SendTestParams represents parameters for sending a test broadcast email
type SendTestParams struct {
	Emails []string `json:"emails"`
}

// BroadcastTestOptions carries retry-safety options for a test send.
type BroadcastTestOptions struct {
	IdempotencyKey string
}

// SendTestResponse represents the response from sending test broadcast emails
type SendTestResponse struct {
	TestEmails int    `json:"testEmails"`
	Message    string `json:"message"`
}

// BroadcastStatusWarning represents the warning sub-object in a broadcast status response
type BroadcastStatusWarning struct {
	Type        BroadcastStatusWarningType `json:"type"`
	Message     string                     `json:"message"`
	FailedCount int                        `json:"failedCount"`
}

// BroadcastStatusProgress represents the progress sub-object (only present while SENDING)
type BroadcastStatusProgress struct {
	Percentage             int     `json:"percentage"`
	EmailsSent             int     `json:"emailsSent"`
	TotalContacts          int     `json:"totalContacts"`
	PendingEmails          int     `json:"pendingEmails"`
	EstimatedTimeRemaining *string `json:"estimatedTimeRemaining,omitempty"`
}

// BroadcastStatusSchedule represents the schedule sub-object (only present while SCHEDULED).
// TimeUntilSend is milliseconds until the scheduled send.
type BroadcastStatusSchedule struct {
	ScheduledAt   *string `json:"scheduledAt,omitempty"`
	TimeUntilSend int64   `json:"timeUntilSend"`
	WillSendIn    string  `json:"willSendIn"`
	IsPastDue     bool    `json:"isPastDue"`
}

// BroadcastStatusCompletion represents the completion sub-object (only present while SENT)
type BroadcastStatusCompletion struct {
	CompletedAt     *string `json:"completedAt,omitempty"`
	TotalEmailsSent int     `json:"totalEmailsSent"`
	SuccessRate     int     `json:"successRate"`
}

// BroadcastStatus represents the broadcast status snapshot
type BroadcastStatus struct {
	ID               string                     `json:"id"`
	Name             string                     `json:"name"`
	Status           BroadcastLifecycleStatus   `json:"status"`
	TotalContacts    int                        `json:"totalContacts"`
	EmailsSent       int                        `json:"emailsSent"`
	FailedEmails     int                        `json:"failedEmails"`
	PendingEmails    int                        `json:"pendingEmails"`
	CreatedAt        time.Time                  `json:"createdAt"`
	SentAt           *time.Time                 `json:"sentAt,omitempty"`
	CompletedAt      *time.Time                 `json:"completedAt,omitempty"`
	ScheduledAt      *time.Time                 `json:"scheduledAt,omitempty"`
	SendingStartedAt *time.Time                 `json:"sendingStartedAt,omitempty"`
	Step             *string                    `json:"step,omitempty"`
	StepLabel        *string                    `json:"stepLabel,omitempty"`
	StepIndex        *int                       `json:"stepIndex,omitempty"`
	StepTotal        *int                       `json:"stepTotal,omitempty"`
	StepProgress     *float64                   `json:"stepProgress,omitempty"`
	Warning          *BroadcastStatusWarning    `json:"warning,omitempty"`
	Progress         *BroadcastStatusProgress   `json:"progress,omitempty"`
	Schedule         *BroadcastStatusSchedule   `json:"schedule,omitempty"`
	Completion       *BroadcastStatusCompletion `json:"completion,omitempty"`
}

// BroadcastAnalytics represents the response from getting broadcast analytics
type BroadcastAnalytics struct {
	TotalEmails        int     `json:"totalEmails"`
	SentEmails         int     `json:"sentEmails"`
	DeliveredEmails    int     `json:"deliveredEmails"`
	OpenedEmails       int     `json:"openedEmails"`
	ClickedEmails      int     `json:"clickedEmails"`
	TotalClicks        int     `json:"totalClicks"`
	BouncedEmails      int     `json:"bouncedEmails"`
	FailedEmails       int     `json:"failedEmails"`
	UnsubscribedEmails int     `json:"unsubscribedEmails"`
	DeliveryRate       float64 `json:"deliveryRate"`
	OpenRate           float64 `json:"openRate"`
	ClickRate          float64 `json:"clickRate"`
	BounceRate         float64 `json:"bounceRate"`
	UnsubscribeRate    float64 `json:"unsubscribeRate"`
}

// ListBroadcastEmailsParams represents parameters for listing broadcast emails.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListBroadcastEmailsParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// BroadcastCostEstimate represents the response from estimating a broadcast cost
type BroadcastCostEstimate struct {
	EstimatedCost        float64 `json:"estimatedCost"`
	EmailPrice           float64 `json:"emailPrice"`
	ContactCount         int     `json:"contactCount"`
	CurrentBalance       float64 `json:"currentBalance"`
	HasSufficientBalance bool    `json:"hasSufficientBalance"`
	Shortfall            float64 `json:"shortfall"`
	FreeEmailsAvailable  int     `json:"freeEmailsAvailable"`
	FreeEmailsToUse      int     `json:"freeEmailsToUse"`
	PaidEmailsCount      int     `json:"paidEmailsCount"`
}

// BroadcastRecipientEstimateParams selects the target whose subscribed
// recipients should be counted.
type BroadcastRecipientEstimateParams struct {
	TargetType BroadcastTargetType `json:"targetType"`
	AudienceID *string             `json:"audienceId,omitempty"`
	TopicID    *string             `json:"topicId,omitempty"`
	SegmentID  *string             `json:"segmentId,omitempty"`
}

// BroadcastRecipientEstimate is the current subscribed recipient count.
type BroadcastRecipientEstimate struct {
	Count int `json:"count"`
}

// BroadcastEmailContact represents the compact contact reference embedded in a
// BroadcastEmail. ID is nil when the contact was deleted after the send.
type BroadcastEmailContact struct {
	ID    *string `json:"id"`
	Email string  `json:"email"`
}

// BroadcastEmail represents a per-recipient email row for a broadcast
type BroadcastEmail struct {
	ID             string                `json:"id"`
	MessageID      *string               `json:"messageId"`
	To             string                `json:"to"`
	Subject        string                `json:"subject"`
	Status         EmailStatus           `json:"status"`
	Type           EmailType             `json:"type"`
	OpenedAt       *time.Time            `json:"openedAt,omitempty"`
	ClickedAt      *time.Time            `json:"clickedAt,omitempty"`
	UnsubscribedAt *time.Time            `json:"unsubscribedAt,omitempty"`
	CreatedAt      time.Time             `json:"createdAt"`
	Contact        BroadcastEmailContact `json:"contact"`
}

// ActiveBroadcastItem represents a single row in the Active Broadcasts response.
// The active endpoint returns a lean projection of the Broadcast row, only the
// fields needed to render dashboard tiles.
type ActiveBroadcastItem struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Status        BroadcastLifecycleStatus `json:"status"`
	SentCount     int                      `json:"sentCount"`
	FailedCount   int                      `json:"failedCount"`
	TotalContacts int                      `json:"totalContacts"`
	ScheduledAt   *time.Time               `json:"scheduledAt,omitempty"`
	AbTest        *ActiveBroadcastAbTest   `json:"abTest,omitempty"`
}

// ActiveBroadcastAbTest is the lean A/B test summary embedded in an
// ActiveBroadcastItem; the active endpoint selects only these three fields.
type ActiveBroadcastAbTest struct {
	Status          BroadcastABTestStatus `json:"status"`
	ScheduledEvalAt *time.Time            `json:"scheduledEvalAt,omitempty"`
	TestPercent     int                   `json:"testPercent"`
}

// ActiveBroadcastsResponse represents the response for listing active broadcasts
type ActiveBroadcastsResponse struct {
	Items []ActiveBroadcastItem `json:"items"`
	Total int                   `json:"total"`
}

// BroadcastCancelResponse represents the response from cancelling a broadcast
type BroadcastCancelResponse struct {
	Broadcast     BroadcastSummary `json:"broadcast"`
	CancelledJobs int              `json:"cancelledJobs"`
}
