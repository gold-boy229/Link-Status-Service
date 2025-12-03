package dto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksGetStatus_Request_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		expectedResult LinksGetStatus_Request
		expectErr      bool
	}{
		{
			name:      "Success case: unmarshal empty string",
			inputJSON: `""`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal not empty string",
			inputJSON: `"aaa.com"`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{"aaa.com"},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal empty list of strings",
			inputJSON: `[]`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal list of one string",
			inputJSON: `["aaa.com"]`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{"aaa.com"},
			},
			expectErr: false,
		},
		{
			name:      "Success case: unmarshal list of two strings",
			inputJSON: `["aaa.com", "bbb.com"]`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{"aaa.com", "bbb.com"},
			},
			expectErr: false,
		},
		{
			name:      "Fail case: unmarshal wrong type (integer)",
			inputJSON: `123`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{},
			},
			expectErr: true,
		},
		{
			name:      "Fail case: unmarshal wrong type (bool)",
			inputJSON: `true`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{},
			},
			expectErr: true,
		},
		{
			name:      "Fail case: unmarshal wrong type (object)",
			inputJSON: `{"links":["aaa.com"]}`,
			expectedResult: LinksGetStatus_Request{
				Links: []string{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var reqData LinksGetStatus_Request
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
