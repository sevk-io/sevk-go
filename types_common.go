package sevk

// JSONValue represents any value that can occur in a JSON document.
// Go has no recursive scalar/container union, so this remains an alias to any.
type JSONValue = any

type UnsubscribeSource string

type AudienceVisibility string

type TemplateStyle string

type BroadcastStyle string

type BroadcastLifecycleStatus string

type BroadcastTargetType string

type BroadcastWinnerMetric string

type BroadcastABTestStatus string

type DomainTLSPolicy string

type DNSRecordType string

type SegmentOperator string

type SegmentRuleOperator string

type EmailStatus string

type EmailType string

type BroadcastStatusWarningType string

type WebhookEventType string

type WebhookCategory string

type ActivityEventType string

type ActivityAction string

type ActivityCategory string

type InboundEmailStatus string

type InboundBlocklistType string

const (
	UnsubscribeSourceUser      UnsubscribeSource = "USER"
	UnsubscribeSourceAPI       UnsubscribeSource = "API"
	UnsubscribeSourceAdmin     UnsubscribeSource = "ADMIN"
	UnsubscribeSourceBounce    UnsubscribeSource = "BOUNCE"
	UnsubscribeSourceComplaint UnsubscribeSource = "COMPLAINT"

	AudienceVisibilityPublic  AudienceVisibility = "PUBLIC"
	AudienceVisibilityPrivate AudienceVisibility = "PRIVATE"

	TemplateStyleHTML   TemplateStyle = "HTML"
	TemplateStyleMarkup TemplateStyle = "MARKUP"

	BroadcastStyleHTML   BroadcastStyle = "HTML"
	BroadcastStyleMarkup BroadcastStyle = "MARKUP"

	BroadcastStatusDraft      BroadcastLifecycleStatus = "DRAFT"
	BroadcastStatusScheduled  BroadcastLifecycleStatus = "SCHEDULED"
	BroadcastStatusStarting   BroadcastLifecycleStatus = "STARTING"
	BroadcastStatusSending    BroadcastLifecycleStatus = "SENDING"
	BroadcastStatusSent       BroadcastLifecycleStatus = "SENT"
	BroadcastStatusFailed     BroadcastLifecycleStatus = "FAILED"
	BroadcastStatusCancelling BroadcastLifecycleStatus = "CANCELLING"
	BroadcastStatusCancelled  BroadcastLifecycleStatus = "CANCELLED"

	BroadcastTargetAll      BroadcastTargetType = "ALL"
	BroadcastTargetAudience BroadcastTargetType = "AUDIENCE"
	BroadcastTargetTopic    BroadcastTargetType = "TOPIC"
	BroadcastTargetSegment  BroadcastTargetType = "SEGMENT"

	BroadcastWinnerOpenRate  BroadcastWinnerMetric = "OPEN_RATE"
	BroadcastWinnerClickRate BroadcastWinnerMetric = "CLICK_RATE"

	BroadcastABTestScheduled  BroadcastABTestStatus = "SCHEDULED"
	BroadcastABTestTesting    BroadcastABTestStatus = "TESTING"
	BroadcastABTestEvaluating BroadcastABTestStatus = "EVALUATING"
	BroadcastABTestRollingOut BroadcastABTestStatus = "ROLLING_OUT"
	BroadcastABTestCompleted  BroadcastABTestStatus = "COMPLETED"
	BroadcastABTestFailed     BroadcastABTestStatus = "FAILED"

	DomainTLSPolicyOptional DomainTLSPolicy = "OPTIONAL"
	DomainTLSPolicyRequire  DomainTLSPolicy = "REQUIRE"

	DNSRecordTypeCNAME DNSRecordType = "CNAME"
	DNSRecordTypeTXT   DNSRecordType = "TXT"
	DNSRecordTypeMX    DNSRecordType = "MX"

	SegmentOperatorAnd SegmentOperator = "AND"
	SegmentOperatorOr  SegmentOperator = "OR"

	SegmentRuleEquals         SegmentRuleOperator = "equals"
	SegmentRuleNotEquals      SegmentRuleOperator = "notEquals"
	SegmentRuleContains       SegmentRuleOperator = "contains"
	SegmentRuleNotContains    SegmentRuleOperator = "notContains"
	SegmentRuleStartsWith     SegmentRuleOperator = "startsWith"
	SegmentRuleEndsWith       SegmentRuleOperator = "endsWith"
	SegmentRuleGreaterThan    SegmentRuleOperator = "greaterThan"
	SegmentRuleLessThan       SegmentRuleOperator = "lessThan"
	SegmentRuleGreaterOrEqual SegmentRuleOperator = "greaterOrEqual"
	SegmentRuleLessOrEqual    SegmentRuleOperator = "lessOrEqual"
	SegmentRuleIsEmpty        SegmentRuleOperator = "isEmpty"
	SegmentRuleIsNotEmpty     SegmentRuleOperator = "isNotEmpty"
	SegmentRuleIn             SegmentRuleOperator = "in"
	SegmentRuleNotIn          SegmentRuleOperator = "notIn"

	EmailStatusPending    EmailStatus = "PENDING"
	EmailStatusScheduled  EmailStatus = "SCHEDULED"
	EmailStatusSent       EmailStatus = "SENT"
	EmailStatusDelivered  EmailStatus = "DELIVERED"
	EmailStatusBounced    EmailStatus = "BOUNCED"
	EmailStatusOpened     EmailStatus = "OPENED"
	EmailStatusComplaint  EmailStatus = "COMPLAINT"
	EmailStatusFailed     EmailStatus = "FAILED"
	EmailStatusSuppressed EmailStatus = "SUPPRESSED"
	EmailStatusCancelled  EmailStatus = "CANCELLED"

	EmailTypeMarketing     EmailType = "MARKETING"
	EmailTypeTransactional EmailType = "TRANSACTIONAL"

	BroadcastWarningFailedEmails BroadcastStatusWarningType = "FAILED_EMAILS"

	WebhookCategoryEmail        WebhookCategory = "email"
	WebhookCategorySubscription WebhookCategory = "subscription"
	WebhookCategoryBroadcast    WebhookCategory = "broadcast"

	ActivityEventEmail        ActivityEventType = "EMAIL"
	ActivityEventSubscription ActivityEventType = "SUBSCRIPTION"

	ActivityActionRequested    ActivityAction = "REQUESTED"
	ActivityActionSent         ActivityAction = "SENT"
	ActivityActionDelivered    ActivityAction = "DELIVERED"
	ActivityActionOpened       ActivityAction = "OPENED"
	ActivityActionClicked      ActivityAction = "CLICKED"
	ActivityActionBounced      ActivityAction = "BOUNCED"
	ActivityActionComplaint    ActivityAction = "COMPLAINT"
	ActivityActionFailed       ActivityAction = "FAILED"
	ActivityActionSubscribed   ActivityAction = "SUBSCRIBED"
	ActivityActionUnsubscribed ActivityAction = "UNSUBSCRIBED"
	ActivityActionDelayed      ActivityAction = "DELAYED"
	ActivityActionRejected     ActivityAction = "REJECTED"

	ActivityCategoryEmail        ActivityCategory = "EMAIL"
	ActivityCategorySubscription ActivityCategory = "SUBSCRIPTION"
	ActivityCategoryAll          ActivityCategory = "all"

	InboundEmailStatusReceived  InboundEmailStatus = "RECEIVED"
	InboundEmailStatusProcessed InboundEmailStatus = "PROCESSED"
	InboundEmailStatusBounced   InboundEmailStatus = "BOUNCED"
	InboundEmailStatusBlocked   InboundEmailStatus = "BLOCKED"
	InboundEmailStatusFailed    InboundEmailStatus = "FAILED"

	InboundBlocklistEmail  InboundBlocklistType = "EMAIL"
	InboundBlocklistDomain InboundBlocklistType = "DOMAIN"
	InboundBlocklistIP     InboundBlocklistType = "IP"
)

// PaginatedResponse represents a paginated API response.
//
// The shape varies by mode:
//   - Offset (?page=): Total, Page, TotalPages are set.
//   - Cursor (default, ?after=, or ?before=): HasMore, NextCursor, and PrevCursor are set.
//
// Items is always present.
type PaginatedResponse[T any] struct {
	Items      []T     `json:"items"`
	Total      *int    `json:"total,omitempty"`
	Page       *int    `json:"page,omitempty"`
	TotalPages *int    `json:"totalPages,omitempty"`
	HasMore    *bool   `json:"hasMore,omitempty"`
	NextCursor *string `json:"nextCursor,omitempty"`
	PrevCursor *string `json:"prevCursor,omitempty"`
}

type ListParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// DeleteResponse represents the standardized response from a DELETE operation
type DeleteResponse struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}
