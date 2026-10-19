package sevk

import "github.com/sevk-io/sevk-go/markup"

// EmailsResource handles email operations
type EmailsResource struct {
	client *Client
}

// Get retrieves a sent email and its current delivery state.
func (r *EmailsResource) Get(id string) (*Email, error) {
	var result Email

	if err := r.client.get("/emails/"+id, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// update moves an unattempted scheduled email to a new delivery time.
func (r *EmailsResource) Update(id string, params UpdateScheduledEmailParams) (*ScheduledEmailMutationResponse, error) {
	var result ScheduledEmailMutationResponse

	if err := r.client.patch("/emails/"+id, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// cancel cancels an unattempted scheduled email.
func (r *EmailsResource) Cancel(id string) (*ScheduledEmailMutationResponse, error) {
	var result ScheduledEmailMutationResponse

	if err := r.client.post("/emails/"+id+"/cancel", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// renderMarkup renders Markup to HTML in-place; the Markup field has json:"-" so it
// never reaches the wire either way, but we surface the rendered content via HTML.
func renderMarkup(p *SendEmailParams) error {
	if p.Markup != "" {
		html, err := markup.Render(p.Markup)
		if err != nil {
			return err
		}
		p.HTML = html
		p.Markup = ""
	}

	return nil
}

func renderBulkMarkup(p *SendBulkEmailParams) error {
	if p.Markup != "" {
		html, err := markup.Render(p.Markup)
		if err != nil {
			return err
		}
		p.HTML = html
		p.Markup = ""
	}

	return nil
}

// SendOptions carries per-call options like an idempotency key.
type SendOptions struct {
	IdempotencyKey string
}

func idempotencyHeaders(opts *SendOptions) map[string]string {
	if opts == nil {
		return nil
	}

	return map[string]string{"Idempotency-Key": opts.IdempotencyKey}
}

// Send queues one transactional email.
func (r *EmailsResource) Send(params SendEmailParams) (SendEmailResponse, error) {
	return r.SendWithOptions(params, nil)
}

// SendWithOptions queues one email with optional retry idempotency.
func (r *EmailsResource) SendWithOptions(params SendEmailParams, opts *SendOptions) (SendEmailResponse, error) {
	if err := renderMarkup(&params); err != nil {
		return nil, err
	}

	var result sendEmailResponse

	if err := r.client.postWithHeaders("/emails", params, idempotencyHeaders(opts), &result); err != nil {
		return nil, err
	}

	return result.value, nil
}

// SendBulk queues up to 100 emails and returns per-item results.
func (r *EmailsResource) SendBulk(emails []SendBulkEmailParams) (*BulkEmailResponse, error) {
	return r.SendBulkWithOptions(emails, nil)
}

// SendBulkWithOptions queues an email batch with optional retry idempotency.
func (r *EmailsResource) SendBulkWithOptions(emails []SendBulkEmailParams, opts *SendOptions) (*BulkEmailResponse, error) {
	wireEmails := make([]SendBulkEmailParams, len(emails))
	copy(wireEmails, emails)

	for i := range wireEmails {
		if err := renderBulkMarkup(&wireEmails[i]); err != nil {
			return nil, err
		}
	}

	body := map[string]any{
		"emails": wireEmails,
	}

	var result BulkEmailResponse

	if err := r.client.postWithHeaders("/emails/bulk", body, idempotencyHeaders(opts), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
