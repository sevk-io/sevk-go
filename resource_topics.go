package sevk

// TopicsResource handles topic operations
type TopicsResource struct {
	client *Client
}

// List retrieves topics belonging to an audience with optional filters.
func (r *TopicsResource) List(audienceID string, params *ListTopicsParams) (*PaginatedResponse[Topic], error) {
	path := "/audiences/" + audienceID + "/topics"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Topic]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a subscription topic under an audience.
func (r *TopicsResource) Create(audienceID string, params CreateTopicParams) (*TopicBase, error) {
	var result TopicBase

	if err := r.client.post("/audiences/"+audienceID+"/topics", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a topic by ID.
func (r *TopicsResource) Get(topicID string) (*Topic, error) {
	var result Topic

	if err := r.client.get("/topics/"+topicID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a topic while preserving omitted fields.
func (r *TopicsResource) Update(topicID string, params UpdateTopicParams) (*TopicBase, error) {
	var result TopicBase

	if err := r.client.put("/topics/"+topicID, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a topic without deleting contacts or its parent audience.
func (r *TopicsResource) Delete(topicID string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/topics/"+topicID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddContacts attaches contacts to a topic, ignoring existing memberships.
func (r *TopicsResource) AddContacts(topicID string, params AddTopicContactsParams) (*AddContactsToTopicResponse, error) {
	var result AddContactsToTopicResponse

	if err := r.client.post("/topics/"+topicID+"/contacts", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RemoveContact detaches a contact from a topic without deleting it.
func (r *TopicsResource) RemoveContact(topicID, contactID string) (*TopicRemoveContactResponse, error) {
	var result TopicRemoveContactResponse

	if err := r.client.delete("/topics/"+topicID+"/contacts/"+contactID, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListContacts retrieves contacts attached to a topic with optional filters.
func (r *TopicsResource) ListContacts(topicID string, params *ListTopicContactsParams) (*PaginatedResponse[Contact], error) {
	path := "/topics/" + topicID + "/contacts"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		addQueryValue(query, "subscribed", params.Subscribed)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Contact]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
