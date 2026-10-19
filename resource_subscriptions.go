package sevk

// SubscriptionsResource handles subscription operations
type SubscriptionsResource struct {
	client *Client
}

// Subscribe opts in a contact and can attach it to an audience.
func (r *SubscriptionsResource) Subscribe(params SubscribeParams) (*SubscribeResponse, error) {
	var result SubscribeResponse

	if err := r.client.post("/subscriptions/subscribe", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Unsubscribe globally opts out a contact and can attribute the action.
func (r *SubscriptionsResource) Unsubscribe(params UnsubscribeParams) (*UnsubscribeResponse, error) {
	var result UnsubscribeResponse

	if err := r.client.post("/subscriptions/unsubscribe", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
