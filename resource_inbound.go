package sevk

// InboundResource handles inbound email operations and per-project inbound settings.
type InboundResource struct {
	client *Client
}

// List retrieves received email with optional pagination and filtering.
func (r *InboundResource) List(params ListInboundRequest) (*InboundListResponse, error) {
	path := "/inbound"
	query := newQueryParams()
	addPagination(query, params.Page, params.Limit, params.After, params.Before)
	addQueryValue(query, "search", params.Search)
	addQueryValue(query, "status", params.Status)
	addQueryValue(query, "domainId", params.DomainID)
	addQueryValue(query, "broadcastId", params.BroadcastID)
	addQueryValue(query, "isSpam", params.IsSpam)
	addQueryValue(query, "isVirus", params.IsVirus)
	path = query.appendTo(path)

	var result InboundListResponse

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a received email with parsed content and attachments.
func (r *InboundResource) Get(id string) (*InboundEmail, error) {
	var result InboundEmail

	if err := r.client.get("/inbound/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a received email and its stored attachments.
func (r *InboundResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/inbound/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Stats retrieves inbound processing statistics, optionally for one domain.
func (r *InboundResource) Stats(domainID string) (*InboundStats, error) {
	path := "/inbound/stats"

	if domainID != "" {
		query := newQueryParams()
		addQueryValue(query, "domainId", &domainID)
		path = query.appendTo(path)
	}

	var result InboundStats

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DownloadAttachment creates a short-lived URL for an inbound attachment.
func (r *InboundResource) DownloadAttachment(inboundEmailID, attachmentID string) (*InboundAttachmentDownload, error) {
	var result InboundAttachmentDownload

	if err := r.client.get("/inbound/"+inboundEmailID+"/attachments/"+attachmentID+"/download", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListRoutes retrieves rules used to store or forward matching inbound email.
func (r *InboundResource) ListRoutes() (*InboundRouteListResponse, error) {
	return r.ListRoutesWithParams(nil)
}

// ListRoutesWithParams retrieves routing rules using offset or cursor pagination.
func (r *InboundResource) ListRoutesWithParams(params *ListParams) (*InboundRouteListResponse, error) {
	path := "/inbound-settings/routes"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		path = query.appendTo(path)
	}

	var result InboundRouteListResponse

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetRoute retrieves an inbound routing rule by ID.
func (r *InboundResource) GetRoute(id string) (*InboundRoute, error) {
	var result InboundRoute

	if err := r.client.get("/inbound-settings/routes/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateRoute adds an inbound routing rule.
func (r *InboundResource) CreateRoute(params CreateInboundRouteRequest) (*InboundRoute, error) {
	var result InboundRoute

	if err := r.client.post("/inbound-settings/routes", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateRoute changes an inbound routing rule while preserving omitted fields.
func (r *InboundResource) UpdateRoute(id string, params UpdateInboundRouteRequest) (*InboundRoute, error) {
	var result InboundRoute

	if err := r.client.put("/inbound-settings/routes/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteRoute removes a routing rule without deleting stored email.
func (r *InboundResource) DeleteRoute(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/inbound-settings/routes/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListBlocklist retrieves addresses blocked from inbound processing.
func (r *InboundResource) ListBlocklist() (*InboundBlocklistListResponse, error) {
	return r.ListBlocklistWithParams(nil)
}

// ListBlocklistWithParams retrieves blocked addresses using offset or cursor pagination.
func (r *InboundResource) ListBlocklistWithParams(params *ListParams) (*InboundBlocklistListResponse, error) {
	path := "/inbound-settings/blocklist"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		path = query.appendTo(path)
	}

	var result InboundBlocklistListResponse

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddBlocklistEntry blocks an address from inbound processing.
func (r *InboundResource) AddBlocklistEntry(params AddInboundBlocklistRequest) (*InboundBlocklistEntry, error) {
	var result InboundBlocklistEntry

	if err := r.client.post("/inbound-settings/blocklist", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemoveBlocklistEntry removes an address from the inbound blocklist.
func (r *InboundResource) RemoveBlocklistEntry(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/inbound-settings/blocklist/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetQuota retrieves inbound quota settings and current usage counters.
func (r *InboundResource) GetQuota() (*InboundQuota, error) {
	var result InboundQuota

	if err := r.client.get("/inbound-settings/quota", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateQuota changes the inbound monthly limit while preserving omitted fields.
func (r *InboundResource) UpdateQuota(params UpdateInboundQuotaRequest) (*InboundQuota, error) {
	var result InboundQuota

	if err := r.client.put("/inbound-settings/quota", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResetQuota clears a quota pause without resetting current-period usage.
func (r *InboundResource) ResetQuota() (*InboundQuota, error) {
	var result InboundQuota

	if err := r.client.post("/inbound-settings/quota/reset", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetSecurityConfig retrieves effective inbound security settings for a domain.
func (r *InboundResource) GetSecurityConfig(domainID string) (*InboundSecurityConfig, error) {
	var result InboundSecurityConfig

	if err := r.client.get("/inbound-settings/security/"+domainID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateSecurityConfig changes inbound security settings for a domain.
func (r *InboundResource) UpdateSecurityConfig(domainID string, params UpdateInboundSecurityRequest) (*InboundSecurityConfig, error) {
	var result InboundSecurityConfig

	if err := r.client.put("/inbound-settings/security/"+domainID, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyDomain checks inbound MX records and returns diagnostic details.
func (r *InboundResource) VerifyDomain(id string) (*InboundVerifyResult, error) {
	var result InboundVerifyResult

	if err := r.client.post("/domains/"+id+"/inbound/verify", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
