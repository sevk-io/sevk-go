package sevk

import (
	"strings"
)

// ActivityResource handles activity operations
type ActivityResource struct {
	client *Client
}

// List retrieves project activity with pagination and optional filters.
func (r *ActivityResource) List(params *ListActivityParams) (*PaginatedResponse[ActivityItem], error) {
	path := "/activity"

	if params != nil {
		query := newQueryParams()
		addPagination(query, params.Page, params.Limit, params.After, params.Before)
		addQueryValue(query, "category", params.Category)
		if params.Type != nil {
			values := make([]string, len(params.Type))
			for index, action := range params.Type {
				values[index] = string(action)
			}
			query.Set("type", strings.Join(values, ","))
		}
		addQueryValue(query, "days", params.Days)
		addQueryValue(query, "search", params.Search)
		path = query.appendTo(path)
	}

	var result PaginatedResponse[ActivityItem]

	if err := r.client.get(path, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Stats retrieves aggregate counts and rates for recent project activity.
func (r *ActivityResource) Stats() (*ActivityStats, error) {
	var result ActivityStats

	if err := r.client.get("/activity/stats", &result); err != nil {
		return nil, err
	}

	return &result, nil
}
