package sevk

// DomainsResource handles domain operations
type DomainsResource struct {
	client *Client
}

// List retrieves sending domains, optionally filtered by verification state.
func (r *DomainsResource) List(params *ListDomainsParams) (*DomainsListResponse, error) {
	path := "/domains"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "verified", params.Verified)
		path = query.appendTo(path)
	}

	var result DomainsListResponse

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get retrieves a domain and its current verification state.
func (r *DomainsResource) Get(id string) (*DomainRead, error) {
	var result DomainRead

	if err := r.client.get("/domains/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// create registers a domain; dns material requires domain:dns-read.
func (r *DomainsResource) Create(params CreateDomainParams) (*DomainRead, error) {
	var result DomainRead

	if err := r.client.post("/domains", params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Delete removes a domain after the API confirms it is safe to remove.
func (r *DomainsResource) Delete(id string) (*DeleteResponse, error) {
	var result DeleteResponse

	if err := r.client.delete("/domains/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Update changes domain sending, tracking, TLS, or inbound settings.
func (r *DomainsResource) Update(id string, data UpdateDomainRequest) (*DomainRead, error) {
	var result DomainRead

	if err := r.client.put("/domains/"+id, data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Verify rechecks sending DNS records and returns the updated state.
func (r *DomainsResource) Verify(id string) (*DomainRead, error) {
	var result DomainRead

	if err := r.client.post("/domains/"+id+"/verify", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyInbound checks inbound MX records and returns diagnostic details.
func (r *DomainsResource) VerifyInbound(id string) (*InboundVerifyResult, error) {
	var result InboundVerifyResult

	if err := r.client.post("/domains/"+id+"/inbound/verify", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DnsRecords retrieves records required for sending and inbound mail.
func (r *DomainsResource) DnsRecords(id string) (*DnsRecordsResponse, error) {
	var result DnsRecordsResponse

	if err := r.client.get("/domains/"+id+"/dns-records", &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Regions retrieves region identifiers accepted when registering a domain.
func (r *DomainsResource) Regions() (*RegionsResponse, error) {
	var result RegionsResponse

	if err := r.client.get("/domains/regions", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
