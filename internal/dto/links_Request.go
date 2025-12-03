package dto

import "encoding/json"

type LinksGetStatus_Request struct {
	Links []string `json:"links" validate:"required,min=1"`
}

func (req *LinksGetStatus_Request) UnmarshalJSON(data []byte) error {
	var singleLinkStr string
	err := json.Unmarshal(data, &singleLinkStr)
	if err == nil {
		if singleLinkStr == "" {
			req.Links = make([]string, 0)
		} else {
			req.Links = []string{singleLinkStr}
		}
		return nil
	}
	return json.Unmarshal(data, &req.Links)
}
