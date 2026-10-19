package sevk

import "time"

// domainsummary contains the non-dns metadata returned by list operations.
type DomainSummary struct {
	ID              string          `json:"id"`
	Domain          string          `json:"domain"`
	Email           string          `json:"email"`
	SenderName      *string         `json:"senderName,omitempty"`
	Region          string          `json:"region"`
	Verified        bool            `json:"verified"`
	ClickTracking   bool            `json:"clickTracking"`
	OpenTracking    bool            `json:"openTracking"`
	TlsPolicy       DomainTLSPolicy `json:"tlsPolicy,omitempty"`
	InboundEnabled  bool            `json:"inboundEnabled,omitempty"`
	InboundVerified bool            `json:"inboundVerified,omitempty"`
	InboundPrefix   string          `json:"inboundPrefix"`
	ProjectID       string          `json:"projectId"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

// domainread contains domain metadata plus dns material when the caller has
// domain:dns-read.
type DomainRead struct {
	DomainSummary
	VerificationToken *string     `json:"verificationToken,omitempty"`
	DkimTokens        []string    `json:"dkimTokens,omitempty"`
	DnsRecords        []DnsRecord `json:"dnsRecords,omitempty"`
}

// InboundMxRecord represents a single MX record resolved during inbound verification
type InboundMxRecord struct {
	Exchange string `json:"exchange"`
	Priority int    `json:"priority"`
}

// inboundverifyresult is the capability-aware domain payload plus inbound mx diagnostics.
type InboundVerifyResult struct {
	DomainRead
	InboundMxResolved     []InboundMxRecord `json:"inboundMxResolved,omitempty"`
	InboundExpectedTarget string            `json:"inboundExpectedTarget,omitempty"`
}

// DomainsListResponse represents the response for listing domains
type DomainsListResponse = PaginatedResponse[DomainSummary]

// ListDomainsParams represents parameters for listing domains
type ListDomainsParams struct {
	Page     *int    `json:"page,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
	After    *string `json:"after,omitempty"`
	Before   *string `json:"before,omitempty"`
	Verified *bool   `json:"verified,omitempty"`
}

// CreateDomainParams represents parameters for creating a domain
type CreateDomainParams struct {
	Domain     string  `json:"domain"`
	Email      string  `json:"email"`
	SenderName *string `json:"senderName,omitempty"`
	Region     *string `json:"region,omitempty"`
}

// DnsRecord represents a DNS record for domain verification
type DnsRecord struct {
	Type  DNSRecordType `json:"type"`
	Name  string        `json:"name"`
	FQDN  string        `json:"fqdn"`
	Value string        `json:"value"`
}

// DnsRecordsResponse represents the response for getting DNS records
type DnsRecordsResponse struct {
	Items []DnsRecord `json:"items"`
}

// RegionsResponse represents the response for listing domain regions
type RegionsResponse struct {
	Items []string `json:"items"`
}

// UpdateDomainRequest represents parameters for updating a domain
type UpdateDomainRequest struct {
	Email          *string          `json:"email,omitempty"`
	SenderName     *string          `json:"senderName,omitempty"`
	ClickTracking  *bool            `json:"clickTracking,omitempty"`
	OpenTracking   *bool            `json:"openTracking,omitempty"`
	TlsPolicy      *DomainTLSPolicy `json:"tlsPolicy,omitempty"`
	InboundEnabled *bool            `json:"inboundEnabled,omitempty"`
	InboundPrefix  *string          `json:"inboundPrefix,omitempty"`
}
