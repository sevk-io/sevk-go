package sevk

import (
	"encoding/json"
	"errors"
)

type sendEmailResponse struct {
	value SendEmailResponse
}

func (response *sendEmailResponse) UnmarshalJSON(data []byte) error {
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(data, &envelope); err != nil {
		return err
	}

	id, hasID := envelope["id"]
	ids, hasIDs := envelope["ids"]
	if hasID == hasIDs {
		return errors.New("expected exactly one of id or ids")
	}

	if hasID {
		if string(id) == "null" {
			return errors.New("expected id to be a string")
		}

		var single SendEmailIDResponse
		if err := json.Unmarshal(id, &single.ID); err != nil {
			return err
		}
		response.value = single
		return nil
	}
	if string(ids) == "null" {
		return errors.New("expected ids to be a string array")
	}

	var rawIDs []json.RawMessage
	if err := json.Unmarshal(ids, &rawIDs); err != nil {
		return err
	}
	multiple := SendEmailIDsResponse{IDs: make([]string, len(rawIDs))}
	for index, rawID := range rawIDs {
		if string(rawID) == "null" {
			return errors.New("expected ids to contain strings")
		}
		if err := json.Unmarshal(rawID, &multiple.IDs[index]); err != nil {
			return err
		}
	}
	response.value = multiple
	return nil
}
