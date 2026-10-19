package sevk

// SegmentsResource handles segment operations
type SegmentsResource struct {
	client *Client
}

// List retrieves segments belonging to an audience with optional filters.
func (r *SegmentsResource) List(audienceID string, params *ListSegmentsParams) (*PaginatedResponse[Segment], error) {
	path := "/audiences/" + audienceID + "/segments"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Segment]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a saved contact segment under an audience.
func (r *SegmentsResource) Create(audienceID string, params CreateSegmentParams) (*Segment, error) {
	var result Segment

	if err := r.client.post("/audiences/"+audienceID+"/segments", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Preview evaluates segment rules and returns matches without saving them.
func (r *SegmentsResource) Preview(audienceID string, params PreviewSegmentParams) (*PreviewSegmentResponse, error) {
	var result PreviewSegmentResponse

	if err := r.client.post("/audiences/"+audienceID+"/segments/preview", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a saved segment by ID.
func (r *SegmentsResource) Get(segmentID string) (*Segment, error) {
	var result Segment

	if err := r.client.get("/segments/"+segmentID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a segment while preserving omitted fields.
func (r *SegmentsResource) Update(segmentID string, params UpdateSegmentParams) (*Segment, error) {
	var result Segment

	if err := r.client.put("/segments/"+segmentID, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a segment without deleting any contacts.
func (r *SegmentsResource) Delete(segmentID string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/segments/"+segmentID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Calculate refreshes and returns the current matching contact count.
func (r *SegmentsResource) Calculate(segmentID string) (*CalculateSegmentResponse, error) {
	var result CalculateSegmentResponse

	if err := r.client.get("/segments/"+segmentID+"/calculate", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
