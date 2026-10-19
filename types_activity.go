package sevk

import "time"

// ActivityItem represents an activity entry in sevk
type ActivityItem struct {
	ID            string            `json:"id"`
	EventType     ActivityEventType `json:"eventType"`
	Action        ActivityAction    `json:"action"`
	Description   *string           `json:"description,omitempty"`
	EmailID       *string           `json:"emailId,omitempty"`
	EmailSubject  *string           `json:"emailSubject,omitempty"`
	ContactID     *string           `json:"contactId,omitempty"`
	ContactEmail  *string           `json:"contactEmail,omitempty"`
	BroadcastID   *string           `json:"broadcastId,omitempty"`
	BroadcastName *string           `json:"broadcastName,omitempty"`
	Link          *string           `json:"link,omitempty"`
	UserAgent     *string           `json:"userAgent,omitempty"`
	Metadata      *string           `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
}

// ListActivityParams represents parameters for listing activity.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListActivityParams struct {
	Page     *int              `json:"page,omitempty"`
	Limit    *int              `json:"limit,omitempty"`
	After    *string           `json:"after,omitempty"`
	Before   *string           `json:"before,omitempty"`
	Category *ActivityCategory `json:"category,omitempty"`
	Type     []ActivityAction  `json:"type,omitempty"`
	Days     *int              `json:"days,omitempty"`
	Search   *string           `json:"search,omitempty"`
}

// ActivityRates represents activity rate statistics
type ActivityRates struct {
	Open     float64 `json:"open"`
	Click    float64 `json:"click"`
	Bounce   float64 `json:"bounce"`
	Delivery float64 `json:"delivery"`
}

// ActivityStats represents activity statistics
type ActivityStats struct {
	TotalEmails         int           `json:"totalEmails"`
	SentEmails          int           `json:"sentEmails"`
	DeliveredEmails     int           `json:"deliveredEmails"`
	OpenedEmails        int           `json:"openedEmails"`
	ClickedEmails       int           `json:"clickedEmails"`
	BouncedEmails       int           `json:"bouncedEmails"`
	MarketingEmails     int           `json:"marketingEmails"`
	TransactionalEmails int           `json:"transactionalEmails"`
	Rates               ActivityRates `json:"rates"`
}
