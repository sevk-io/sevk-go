package sevk

import "time"

// Template represents an email template in sevk
type Template struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Style     TemplateStyle `json:"style"`
	ProjectID string        `json:"projectId"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

// CreateTemplateParams represents parameters for creating a template
type CreateTemplateParams struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Style   *TemplateStyle `json:"style,omitempty"`
}

// UpdateTemplateParams represents parameters for updating a template
type UpdateTemplateParams struct {
	Title   *string        `json:"title,omitempty"`
	Content *string        `json:"content,omitempty"`
	Style   *TemplateStyle `json:"style,omitempty"`
}

// ListTemplatesParams represents parameters for listing templates.
// Before enables backward cursor pagination; supply the cursor of the
// first item in the current page to fetch the previous page.
type ListTemplatesParams struct {
	Page   *int    `json:"page,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
	Search *string `json:"search,omitempty"`
}
