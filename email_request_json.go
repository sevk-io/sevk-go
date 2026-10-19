package sevk

import "encoding/json"

func (params SendEmailParams) MarshalJSON() ([]byte, error) {
	type wire SendEmailParams
	return json.Marshal(struct {
		wire
		ReplyTo *StringOrList `json:"reply_to,omitempty"`
		CC      *StringOrList `json:"cc,omitempty"`
		BCC     *StringOrList `json:"bcc,omitempty"`
	}{wire(params), optionalRecipients(params.ReplyTo), optionalRecipients(params.CC), optionalRecipients(params.BCC)})
}

func (params SendBulkEmailParams) MarshalJSON() ([]byte, error) {
	type wire SendBulkEmailParams
	return json.Marshal(struct {
		wire
		ReplyTo *StringOrList `json:"reply_to,omitempty"`
		CC      *StringOrList `json:"cc,omitempty"`
		BCC     *StringOrList `json:"bcc,omitempty"`
	}{wire(params), optionalRecipients(params.ReplyTo), optionalRecipients(params.CC), optionalRecipients(params.BCC)})
}

func optionalRecipients(value StringOrList) *StringOrList {
	if value == nil {
		return nil
	}
	return &value
}
