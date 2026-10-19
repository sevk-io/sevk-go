package sevk

// OutboundQuota is the project's outbound quota state returned by the
// /outbound-settings/quota endpoints.
type OutboundQuota struct {
	ProjectID                  string  `json:"projectId"`
	MonthlyEmailLimit          *int    `json:"monthlyEmailLimit"`
	AdminMonthlyEmailLimit     *int    `json:"adminMonthlyEmailLimit"`
	EffectiveMonthlyEmailLimit *int    `json:"effectiveMonthlyEmailLimit"`
	MonthlyEmailsUsed          int     `json:"monthlyEmailsUsed"`
	LimitEmailsUsed            int     `json:"limitEmailsUsed"`
	MonthlyRemaining           *int    `json:"monthlyRemaining"`
	FreeMonthlyEmails          *int    `json:"freeMonthlyEmails"`
	FreeEmailsUsed             int     `json:"freeEmailsUsed"`
	RateLimit                  *int    `json:"rateLimit"`
	EmailPrice                 float64 `json:"emailPrice"`
	ResetAt                    string  `json:"resetAt"`
	PausedAt                   *string `json:"pausedAt"`
	PausedReason               *string `json:"pausedReason"`
	CreatedAt                  string  `json:"createdAt"`
	UpdatedAt                  string  `json:"updatedAt"`
}

// UpdateOutboundQuotaRequest sets the user-side monthly outbound cap.
// MonthlyEmailLimit is required (the endpoint rejects empty bodies).
// Pass nil to clear the user-set cap; the admin ceiling still applies.
type UpdateOutboundQuotaRequest struct {
	MonthlyEmailLimit *int `json:"monthlyEmailLimit"`
}
