package sevk

// AudiencesResource handles audience operations
type AudiencesResource struct {
	client *Client
}

// List retrieves audiences with optional pagination and search.
func (r *AudiencesResource) List(params *ListAudiencesParams) (*PaginatedResponse[Audience], error) {
	path := "/audiences"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Audience]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves an audience by ID.
func (r *AudiencesResource) Get(id string) (*Audience, error) {
	var result Audience

	if err := r.client.get("/audiences/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds an audience for organizing contacts, topics, and segments.
func (r *AudiencesResource) Create(params CreateAudienceParams) (*AudienceBase, error) {
	var result AudienceBase

	if err := r.client.post("/audiences", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes an audience while preserving omitted fields.
func (r *AudiencesResource) Update(id string, params UpdateAudienceParams) (*AudienceBase, error) {
	var result AudienceBase

	if err := r.client.put("/audiences/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes an audience without deleting its contacts.
func (r *AudiencesResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/audiences/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddContacts attaches contacts to an audience, ignoring existing memberships.
func (r *AudiencesResource) AddContacts(audienceID string, contactIDs []string) (*AddContactsToAudienceResponse, error) {
	body := map[string]any{
		"contactIds": contactIDs,
	}

	var result AddContactsToAudienceResponse

	if err := r.client.post("/audiences/"+audienceID+"/contacts", body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListContacts retrieves contacts attached to an audience with optional filters.
func (r *AudiencesResource) ListContacts(audienceID string, params *ListAudienceContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/audiences/" + audienceID + "/contacts"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		addQueryValue(query, "subscribed", params.Subscribed)
		addQueryValue(query, "excludeTopicId", params.ExcludeTopicID)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Contact]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemoveContact detaches a contact from an audience without deleting it.
func (r *AudiencesResource) RemoveContact(audienceID, contactID string) (*AudienceRemoveContactResponse, error) {
	var result AudienceRemoveContactResponse

	if err := r.client.delete("/audiences/"+audienceID+"/contacts/"+contactID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
