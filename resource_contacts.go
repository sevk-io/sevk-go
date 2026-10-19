package sevk

// ContactsResource handles contact operations
type ContactsResource struct {
	client *Client
}

// List retrieves contacts with optional pagination, search, and filters.
func (r *ContactsResource) List(params *ListContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/contacts"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		addQueryValue(query, "excludeAudience", params.ExcludeAudience)
		addQueryValue(query, "subscribed", params.Subscribed)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Contact]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a contact by ID.
func (r *ContactsResource) Get(id string) (*Contact, error) {
	var result Contact

	if err := r.client.get("/contacts/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a contact and optionally stores custom metadata.
func (r *ContactsResource) Create(params CreateContactParams) (*Contact, error) {
	var result Contact

	if err := r.client.post("/contacts", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a contact while preserving omitted fields.
func (r *ContactsResource) Update(id string, params UpdateContactParams) (*Contact, error) {
	var result Contact

	if err := r.client.put("/contacts/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a contact and its audience and topic memberships.
func (r *ContactsResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/contacts/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// BulkUpdate changes contacts by email and reports entries that were not found.
func (r *ContactsResource) BulkUpdate(params BulkUpdateContactsParams) (*BulkUpdateContactsResponse, error) {
	var result BulkUpdateContactsResponse

	if err := r.client.put("/contacts/bulk-update", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// BulkImport adds contacts and can attach them to an audience or topic.
func (r *ContactsResource) BulkImport(params ImportContactsParams) (*ImportContactsResponse, error) {
	var result ImportContactsResponse

	if err := r.client.post("/contacts/import", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Activity retrieves delivery and engagement history for a contact.
func (r *ContactsResource) Activity(id string, params *ListContactActivityParams) (*PaginatedResponse[ContactActivityItem], error) {
	path := "/contacts/" + id + "/activity"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[ContactActivityItem]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
