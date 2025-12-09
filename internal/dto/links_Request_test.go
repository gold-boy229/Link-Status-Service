package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksGetStatus_Request_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputJSON      string
		expectedResult LinksGetStatusRequest
		expectErr      bool
	}{
		{
			name:      "Success case: unmarshal empty string",
			inputJSON: `{"links":""}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal not empty string",
			inputJSON: `{"links":"aaa.com"}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{"aaa.com"},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal empty list of strings",
			inputJSON: `{"links":[]}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal list of one string",
			inputJSON: `{"links":["aaa.com"]}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{"aaa.com"},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal list of two strings",
			inputJSON: `{"links":["aaa.com", "bbb.com"]}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{"aaa.com", "bbb.com"},
			},
			expectErr: false,
		},
		{
			name:      "Fail case: unmarshal wrong type (integer)",
			inputJSON: `{"links":123}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{},
			},
			expectErr: true,
		},
		{
			name:      "Fail case: unmarshal wrong type (bool)",
			inputJSON: `{"links":true}`,
			expectedResult: LinksGetStatusRequest{
				Links: []string{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var reqData LinksGetStatusRequest
			err := json.Unmarshal([]byte(tt.inputJSON), &reqData)

			if tt.expectErr {
				assert.NotNil(t, err, "Expected an error but got nil")
			} else {
				assert.Nil(t, err, "Expected no error but got one: %v", err)
				assert.EqualValues(t, tt.expectedResult, reqData)
			}
		})
	}
}
