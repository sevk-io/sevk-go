package sevk

// UsageResource handles project usage and limit queries.
type UsageResource struct {
	client *Client
}

// Get retrieves current limits, resource counts, and email usage.
func (r *UsageResource) Get() (*ProjectUsage, error) {
	var result ProjectUsage

	if err := r.client.get("/limits", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
