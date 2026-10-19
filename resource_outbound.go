package sevk

// OutboundResource handles the project's outbound sending quota.
type OutboundResource struct {
	client *Client
}

// GetQuota retrieves outbound quota settings and current usage counters.
func (r *OutboundResource) GetQuota() (*OutboundQuota, error) {
	var result OutboundQuota

	if err := r.client.get("/outbound-settings/quota", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateQuota changes the cap and clears a cap-exhausted pause when possible.
func (r *OutboundResource) UpdateQuota(params UpdateOutboundQuotaRequest) (*OutboundQuota, error) {
	var result OutboundQuota

	if err := r.client.put("/outbound-settings/quota", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResetQuota clears an outbound pause while preserving usage and the cap.
func (r *OutboundResource) ResetQuota() (*OutboundQuota, error) {
	var result OutboundQuota

	if err := r.client.post("/outbound-settings/quota/reset", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ResetLimit resets custom-limit usage while preserving the cap and monthly usage.
func (r *OutboundResource) ResetLimit() (*OutboundQuota, error) {
	var result OutboundQuota

	if err := r.client.post("/outbound-settings/quota/reset-limit", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// PauseQuota stops transactional sending until ResetQuota is called.
func (r *OutboundResource) PauseQuota() (*OutboundQuota, error) {
	var result OutboundQuota

	if err := r.client.post("/outbound-settings/quota/pause", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
