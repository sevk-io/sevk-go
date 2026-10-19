package sevk

import "time"

// SegmentRule represents a single rule in a segment
type SegmentRule struct {
	Field         string              `json:"field"`
	Operator      SegmentRuleOperator `json:"operator"`
	Value         any                 `json:"value"`
	GroupID       *string             `json:"groupId,omitempty"`
	GroupOperator *SegmentOperator    `json:"groupOperator,omitempty"`
}

// SegmentGroup represents a group of segment rules
type SegmentGroup struct {
	ID       string          `json:"id"`
	Operator SegmentOperator `json:"operator"`
}

// Segment represents a segment in sevk
type Segment struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	Rules         []SegmentRule   `json:"rules"`
	Operator      SegmentOperator `json:"operator"`
	LimitToTopics []string        `json:"limitToTopics"`
	Groups        []SegmentGroup  `json:"groups,omitempty"`
	AudienceID    string          `json:"audienceId"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	LastRunAt     *time.Time      `json:"lastRunAt,omitempty"`
	LastRunCount  *int            `json:"lastRunCount,omitempty"`
}

// CreateSegmentParams represents parameters for creating a segment
type CreateSegmentParams struct {
	Name          string          `json:"name"`
	Rules         []SegmentRule   `json:"rules"`
	Operator      SegmentOperator `json:"operator"`
	LimitToTopics []string        `json:"limitToTopics,omitempty"`
	Groups        []SegmentGroup  `json:"groups,omitempty"`
}

// UpdateSegmentParams represents parameters for updating a segment.
// Rules, LimitToTopics, and Groups are pointers so a nil value omits the field
// (leaving it unchanged) while a pointer to an empty slice serializes as [] and
// clears it server-side.
type UpdateSegmentParams struct {
	Name          *string          `json:"name,omitempty"`
	Rules         *[]SegmentRule   `json:"rules,omitempty"`
	Operator      *SegmentOperator `json:"operator,omitempty"`
	LimitToTopics *[]string        `json:"limitToTopics,omitempty"`
	Groups        *[]SegmentGroup  `json:"groups,omitempty"`
}

// ListSegmentsParams represents parameters for listing segments.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListSegmentsParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
	Search *string `json:"search,omitempty"`
}

// PreviewSegmentParams represents parameters for previewing a segment
type PreviewSegmentParams struct {
	Rules         []SegmentRule   `json:"rules"`
	Operator      SegmentOperator `json:"operator"`
	LimitToTopics []string        `json:"limitToTopics,omitempty"`
	Groups        []SegmentGroup  `json:"groups,omitempty"`
}

// PreviewSegmentContact represents a contact in a segment preview response
type PreviewSegmentContact struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Subscribed bool           `json:"subscribed"`
	Data       map[string]any `json:"data,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
}

// PreviewSegmentResponse represents the response from previewing a segment
type PreviewSegmentResponse struct {
	Count    int                     `json:"count"`
	Contacts []PreviewSegmentContact `json:"contacts"`
}

// CalculateSegmentResponse represents the response from calculating a segment
type CalculateSegmentResponse struct {
	Count       int       `json:"count"`
	SegmentID   string    `json:"segmentId"`
	SegmentName string    `json:"segmentName"`
	LastRunAt   time.Time `json:"lastRunAt"`
}
