package dto

import (
	"encoding/json"
)

type LinksGetStatusRequest struct {
	Links LinkList `json:"links" validate:"required,min=1"`
}

type LinkList []string

func (linkList *LinkList) UnmarshalJSON(data []byte) error {
	var singleLinkStr string
	err := json.Unmarshal(data, &singleLinkStr)
	if err == nil {
		if singleLinkStr == "" {
			*linkList = LinkList([]string{})
		} else {
			*linkList = LinkList([]string{singleLinkStr})
		}
		return nil
	}

	var arr []string
	err = json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	*linkList = LinkList(arr)
	return nil
}

type LinkBuildPDFRequest struct {
	LinkNums []int `json:"links_list" validate:"required,min=1"`
}
