package dto

import (
	"Link-Status-Service/internal/consts"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksStatus_Response_MarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		inputData    LinksStatus_Response
		expectedJSON string
		expectErr    bool
	}{
		{
			name:         "Success case: marshal empty list of links",
			inputData:    LinksStatus_Response{},
			expectedJSON: `{}`,
			expectErr:    false,
		},
		{
			name: "Success case: marshal not empty list of links",
			inputData: LinksStatus_Response{
				LinkStatus_Response{Address: "aaa.com", Status: consts.AVAILABLE},
				LinkStatus_Response{Address: "bbb.com", Status: consts.NOT_AVAILABLE},
			},
			expectedJSON: fmt.Sprintf(`{"aaa.com":%q,"bbb.com":%q}`, consts.AVAILABLE, consts.NOT_AVAILABLE),
			expectErr:    false,
		},
		{
			name: "Success case: marshal same link addresses",
			inputData: LinksStatus_Response{
				LinkStatus_Response{Address: "aaa.com", Status: "first"},
				LinkStatus_Response{Address: "aaa.com", Status: "second"},
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
