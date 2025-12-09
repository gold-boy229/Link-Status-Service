package dto

import (
	"bytes"
	"fmt"
)

type LinksGetStatusResponse struct {
	Links    LinksStatusResponse `json:"links"`
	LinksNum int                 `json:"links_num"`
}

type LinksStatusResponse []LinkStatusResponse

type LinkStatusResponse struct {
	Address string
	Status  string
}

func (links LinksStatusResponse) MarshalJSON() ([]byte, error) {
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
