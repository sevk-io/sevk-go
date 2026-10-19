package sevk

import "encoding/json"

func (params CreateBroadcastParams) MarshalJSON() ([]byte, error) {
	type wire CreateBroadcastParams
	var variants *[]VariantInput
	if params.Variants != nil {
		variants = &params.Variants
	}
	return json.Marshal(struct {
		wire
		Variants *[]VariantInput `json:"variants,omitempty"`
	}{wire(params), variants})
}
