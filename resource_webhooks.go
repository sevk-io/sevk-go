package sevk

// WebhooksResource handles webhook operations
type WebhooksResource struct {
	client *Client
}

func (r *WebhooksResource) Rotate(id string) (*RotateWebhookResponse, error) {
	var result RotateWebhookResponse

	if err := r.client.post("/webhooks/"+id+"/rotate", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// List retrieves configured webhooks and their recent delivery health.
func (r *WebhooksResource) List(params *ListWebhooksParams) (*WebhooksListResponse, error) {
	path := "/webhooks"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		path = query.appendTo(path)
	}

	var result WebhooksListResponse

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a webhook and its recent delivery state.
func (r *WebhooksResource) Get(id string) (*Webhook, error) {
	var result Webhook

	if err := r.client.get("/webhooks/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a webhook and returns its usable whsec_ signing secret exactly once.
func (r *WebhooksResource) Create(params CreateWebhookParams) (*CreatedWebhook, error) {
	var result CreatedWebhook

	if err := r.client.post("/webhooks", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a webhook while preserving omitted fields.
func (r *WebhooksResource) Update(id string, params UpdateWebhookParams) (*Webhook, error) {
	var result Webhook

	if err := r.client.put("/webhooks/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a webhook and stops future deliveries to its URL.
func (r *WebhooksResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/webhooks/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Test sends a test event and returns the delivery result.
func (r *WebhooksResource) Test(id string) (*TestWebhookResponse, error) {
	var result TestWebhookResponse

	if err := r.client.post("/webhooks/"+id+"/test", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Events retrieves event names that can be assigned to a webhook.
func (r *WebhooksResource) Events() (*WebhookEventsListResponse, error) {
	var result WebhookEventsListResponse

	if err := r.client.get("/webhooks/events", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
