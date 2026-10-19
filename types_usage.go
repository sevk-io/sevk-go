package sevk

// FreeEmailsBucket describes the free-email allowance within a usage bucket.
type FreeEmailsBucket struct {
	Current   int  `json:"current"`
	Limit     *int `json:"limit"`
	Remaining int  `json:"remaining"`
}

// InboundUsage contains current-period inbound email usage.
type InboundUsage struct {
	EmailsCount       int              `json:"emailsCount"`
	MonthlyEmailLimit *int             `json:"monthlyEmailLimit"`
	MonthlyEmailsUsed int              `json:"monthlyEmailsUsed"`
	MonthlyRemaining  *int             `json:"monthlyRemaining"`
	FreeEmails        FreeEmailsBucket `json:"freeEmails"`
	ResetAt           *string          `json:"resetAt"`
	PausedAt          *string          `json:"pausedAt"`
	PausedReason      *string          `json:"pausedReason"`
}

// OutboundUsage contains current-period outbound email usage.
type OutboundUsage struct {
	MonthlyEmailLimit          *int             `json:"monthlyEmailLimit"`
	AdminMonthlyEmailLimit     *int             `json:"adminMonthlyEmailLimit"`
	EffectiveMonthlyEmailLimit *int             `json:"effectiveMonthlyEmailLimit"`
	MonthlyEmailsUsed          int              `json:"monthlyEmailsUsed"`
	LimitEmailsUsed            int              `json:"limitEmailsUsed"`
	MonthlyRemaining           *int             `json:"monthlyRemaining"`
	RateLimit                  *int             `json:"rateLimit"`
	FreeEmails                 FreeEmailsBucket `json:"freeEmails"`
	ResetAt                    *string          `json:"resetAt"`
	PausedAt                   *string          `json:"pausedAt"`
	PausedReason               *string          `json:"pausedReason"`
}

// ProjectUsage contains project limits, counts, pricing, and email usage.
type ProjectUsage struct {
	ID                        string        `json:"id"`
	Name                      string        `json:"name"`
	Balance                   *string       `json:"balance,omitempty"`
	OutboundEmailPrice        *string       `json:"outboundEmailPrice,omitempty"`
	InboundEmailPrice         *string       `json:"inboundEmailPrice,omitempty"`
	MonthlyOutboundEmailLimit *int          `json:"monthlyOutboundEmailLimit"`
	OutboundRateLimit         *int          `json:"outboundRateLimit"`
	FreeMonthlyOutboundEmails *int          `json:"freeMonthlyOutboundEmails"`
	FreeMonthlyInboundEmails  *int          `json:"freeMonthlyInboundEmails"`
	AudienceLimit             *int          `json:"audienceLimit"`
	ContactLimit              *int          `json:"contactLimit"`
	BroadcastLimit            *int          `json:"broadcastLimit"`
	DomainLimit               *int          `json:"domainLimit"`
	WebhookLimit              *int          `json:"webhookLimit"`
	StorageLimit              *string       `json:"storageLimit"`
	StorageUsed               string        `json:"storageUsed"`
	AudiencesCount            int64         `json:"audiencesCount"`
	ContactsCount             int64         `json:"contactsCount"`
	BroadcastsCount           int64         `json:"broadcastsCount"`
	MonthlyBroadcastsUsed     int64         `json:"monthlyBroadcastsUsed"`
	DomainsCount              int64         `json:"domainsCount"`
	WebhooksCount             int64         `json:"webhooksCount"`
	InboundEmailsCount        int           `json:"inboundEmailsCount"`
	Inbound                   InboundUsage  `json:"inbound"`
	Outbound                  OutboundUsage `json:"outbound"`
}
