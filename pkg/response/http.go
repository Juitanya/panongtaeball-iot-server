package response

import (
	"encoding/json"
)

type HTTPResponse struct {
	Data any `json:"data"`

	Error error `json:"error"`
}

type marshalResponse struct {
	Data  any     `json:"data"`
	Error *string `json:"error"`
}

func newMarshalResponse(hr *HTTPResponse) *marshalResponse {
	if hr.Error != nil {
		errMsg := hr.Error.Error()
		return &marshalResponse{
			Data:  hr.Data,
			Error: &errMsg,
		}
	}
	return &marshalResponse{
		Data:  hr.Data,
		Error: nil,
	}
}

func (h HTTPResponse) MarshalJSON() ([]byte, error) {
	r := newMarshalResponse(&h)
	return json.Marshal(r)
}
