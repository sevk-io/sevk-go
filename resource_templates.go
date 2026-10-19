package sevk

// TemplatesResource handles template operations
type TemplatesResource struct {
	client *Client
}

// List retrieves reusable email templates with pagination and search.
func (r *TemplatesResource) List(params *ListTemplatesParams) (*PaginatedResponse[Template], error) {
	path := "/templates"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "search", params.Search)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[Template]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a template by ID.
func (r *TemplatesResource) Get(id string) (*Template, error) {
	var result Template

	if err := r.client.get("/templates/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Create adds a reusable email template.
func (r *TemplatesResource) Create(params CreateTemplateParams) (*Template, error) {
	var result Template

	if err := r.client.post("/templates", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes a template while preserving omitted fields.
func (r *TemplatesResource) Update(id string, params UpdateTemplateParams) (*Template, error) {
	var result Template

	if err := r.client.put("/templates/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a template without affecting email already sent from it.
func (r *TemplatesResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/templates/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Duplicate creates an independent template copy for further editing.
func (r *TemplatesResource) Duplicate(id string) (*Template, error) {
	var result Template

	if err := r.client.post("/templates/"+id+"/duplicate", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
