package sevk

// BroadcastsResource handles broadcast operations
type BroadcastsResource struct {
	client *Client
}

// List retrieves broadcasts with pagination, search, and status filters.
func (r *BroadcastsResource) List(params *ListBroadcastsParams) (*PaginatedResponse[BroadcastSummary], error) {
	path := "/broadcasts"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		addQueryValue(query, "status", params.Status)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[BroadcastSummary]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a broadcast by ID with its current delivery state.
func (r *BroadcastsResource) Get(id string) (*Broadcast, error) {
	var result Broadcast

	if err := r.client.get("/broadcasts/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a draft or scheduled broadcast for a target group.
func (r *BroadcastsResource) Create(params CreateBroadcastParams) (*BroadcastWithVariants, error) {
	var result BroadcastWithVariants

	if err := r.client.post("/broadcasts", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a broadcast while preserving omitted fields.
func (r *BroadcastsResource) Update(id string, params UpdateBroadcastParams) (*BroadcastUpdate, error) {
	var result BroadcastUpdate

	if err := r.client.put("/broadcasts/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a broadcast that has not been sent.
func (r *BroadcastsResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/broadcasts/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Send queues a prepared broadcast for delivery.
func (r *BroadcastsResource) Send(id string) (*BroadcastSummary, error) {
	var result BroadcastSummary

	if err := r.client.post("/broadcasts/"+id+"/send", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Cancel stops a broadcast when its state still allows cancellation.
func (r *BroadcastsResource) Cancel(id string) (*BroadcastCancelResponse, error) {
	var result BroadcastCancelResponse

	if err := r.client.post("/broadcasts/"+id+"/cancel", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Test delivers a broadcast to test recipients without starting the campaign.
func (r *BroadcastsResource) Test(id string, params SendTestParams) (*SendTestResponse, error) {
	return r.TestWithOptions(id, params, nil)
}

// TestWithOptions delivers a broadcast to test recipients with optional retry idempotency.
func (r *BroadcastsResource) TestWithOptions(id string, params SendTestParams, options *BroadcastTestOptions) (*SendTestResponse, error) {
	var result SendTestResponse
	var headers map[string]string
	if options != nil {
		headers = map[string]string{"Idempotency-Key": options.IdempotencyKey}
	}

	if err := r.client.postWithHeaders("/broadcasts/"+id+"/test", params, headers, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Analytics retrieves aggregate delivery and engagement metrics.
func (r *BroadcastsResource) Analytics(id string) (*BroadcastAnalytics, error) {
	var result BroadcastAnalytics

	if err := r.client.get("/broadcasts/"+id+"/analytics", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Status retrieves live scheduling and delivery progress.
func (r *BroadcastsResource) Status(id string) (*BroadcastStatus, error) {
	var result BroadcastStatus

	if err := r.client.get("/broadcasts/"+id+"/status", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Emails retrieves email deliveries created for a broadcast.
func (r *BroadcastsResource) Emails(id string, params *ListBroadcastEmailsParams) (*PaginatedResponse[BroadcastEmail], error) {
	path := "/broadcasts/" + id + "/emails"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[BroadcastEmail]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EstimateCost calculates delivery cost without sending the broadcast.
func (r *BroadcastsResource) EstimateCost(id string) (*BroadcastCostEstimate, error) {
	var result BroadcastCostEstimate

	if err := r.client.get("/broadcasts/"+id+"/estimate-cost", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EstimateRecipients resolves the subscribed recipient count for a broadcast target.
func (r *BroadcastsResource) EstimateRecipients(params BroadcastRecipientEstimateParams) (*BroadcastRecipientEstimate, error) {
	var result BroadcastRecipientEstimate

	if err := r.client.post("/broadcasts/estimate-recipients", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Active retrieves broadcasts that are currently scheduled or processing.
func (r *BroadcastsResource) Active() (*ActiveBroadcastsResponse, error) {
	var result ActiveBroadcastsResponse

	if err := r.client.get("/broadcasts/active", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
