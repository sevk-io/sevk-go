package sevk

import (
	"fmt"
	"net/url"
)

type queryValue interface {
	~bool | ~int | ~string
}

type queryParams struct {
	url.Values
}

func newQueryParams() queryParams {
	return queryParams{Values: make(url.Values)}
}

func addQueryValue[T queryValue](query queryParams, key string, value *T) {
	if value != nil {
		query.Set(key, fmt.Sprint(*value))
	}
}

func addPagination(
	query queryParams,
	page, limit *int,
	after, before *string,
) {
	addQueryValue(query, "page", page)
	addQueryValue(query, "limit", limit)
	addQueryValue(query, "after", after)
	addQueryValue(query, "before", before)
}

func (query queryParams) appendTo(path string) string {
	encoded := query.Encode()
	if encoded == "" {
		return path
	}
	return path + "?" + encoded
}
