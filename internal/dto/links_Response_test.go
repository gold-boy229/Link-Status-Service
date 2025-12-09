package dto

import (
	"encoding/json"
	"fmt"
	"testing"

	"Link-Status-Service/internal/consts"

	"github.com/stretchr/testify/assert"
)

func TestLinksStatus_Response_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		inputData    LinksStatusResponse
		expectedJSON string
		expectErr    bool
	}{
		{
			name:         "Success case: marshal empty list of links",
			inputData:    LinksStatusResponse{},
			expectedJSON: `{}`,
			expectErr:    false,
		},
		{
			name: "Success case: marshal not empty list of links",
			inputData: LinksStatusResponse{
				LinkStatusResponse{Address: "aaa.com", Status: consts.Available},
				LinkStatusResponse{Address: "bbb.com", Status: consts.NotAvailable},
			},
			expectedJSON: fmt.Sprintf(`{"aaa.com":%q,"bbb.com":%q}`, consts.Available, consts.NotAvailable),
			expectErr:    false,
		},
		{
			name: "Success case: marshal same link addresses",
			inputData: LinksStatusResponse{
				LinkStatusResponse{Address: "aaa.com", Status: "first"},
				LinkStatusResponse{Address: "aaa.com", Status: "second"},
			},
			expectedJSON: `{"aaa.com":"first","aaa.com":"second"}`,
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resJSON, err := json.Marshal(tt.inputData)
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedJSON, string(resJSON))
			}
		})
	}
}
