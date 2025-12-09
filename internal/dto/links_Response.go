package dto

import (
	"bytes"
	"fmt"
)

type LinksGetStatus_Response struct {
	Links    LinksStatus_Response `json:"links"`
	LinksNum int                  `json:"links_num"`
}

type LinksStatus_Response []LinkStatus_Response

type LinkStatus_Response struct {
	Address string
	Status  string
}

func (links LinksStatus_Response) MarshalJSON() ([]byte, error) {
	buff := bytes.NewBufferString(`{`)
	for i, link := range links {
		if i > 0 {
			buff.WriteString(", ")
		}

		fmt.Fprintf(buff, "%q:%q", link.Address, link.Status)
	}
	buff.WriteString("}")

	return buff.Bytes(), nil
}
