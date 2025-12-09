package handlers

import (
	"Link-Status-Service/internal/consts"
	"Link-Status-Service/internal/dto"
	"Link-Status-Service/internal/entity"
	"Link-Status-Service/internal/mocks"
	"Link-Status-Service/internal/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStatus(t *testing.T) {
	// setup service mock
	mockService := mocks.NewMockLinkService()
	mockPDFBuilder := mocks.NewMockPDFBuilder()
	handler := NewLinkHandler(mockService, mockPDFBuilder)
	validator := utils.NewCustomValidator()

	// table tests declaration
	tests := []struct {
		name              string
		requestBody       string
		responseBody      string
		expectedStatus    int
		expectedErrorCode string

		mockInputData  entity.LinkGetStatus_Params
		mockReturnData entity.LinkGetStatus_Result
		mockReturnErr  error
	}{
		{
			name:           "Success case: not empty array",
			requestBody:    `{"links":["aaa.com","bbb.com"]}`,
			responseBody:   `{"links":{"aaa.com":"available","bbb.com":"not available"},"links_num":1}`,
			expectedStatus: http.StatusOK,

			mockInputData: entity.LinkGetStatus_Params{Links: []string{"aaa.com", "bbb.com"}},
			mockReturnData: entity.LinkGetStatus_Result{
				LinkStates: []entity.LinkState{
					{
						Link:        "aaa.com",
						IsAvailable: true,
					},
					{
						Link:        "bbb.com",
						IsAvailable: false,
					},
				},
				LinkNum: 1,
			},
			mockReturnErr: nil,
		},
		{
			name:           "Success case: one link on input as string",
			requestBody:    `{"links":"aaa.com"}`,
			responseBody:   `{"links":{"aaa.com":"available"},"links_num":1}`,
			expectedStatus: http.StatusOK,

			mockInputData: entity.LinkGetStatus_Params{Links: []string{"aaa.com"}},
			mockReturnData: entity.LinkGetStatus_Result{
				LinkStates: []entity.LinkState{
					{
						Link:        "aaa.com",
						IsAvailable: true,
					},
				},
				LinkNum: 1,
			},
			mockReturnErr: nil,
		},
		{
			name:              "Bad Request: empty list input",
			requestBody:       `{"links":[]}`,
			responseBody:      ``,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorCode: consts.ERROR_CODE_BAD_REQUEST,

			mockInputData:  entity.LinkGetStatus_Params{Links: []string{}},
			mockReturnData: entity.LinkGetStatus_Result{},
			mockReturnErr:  nil,
		},
		{
			name:              "Bad Request: empty single input string",
			requestBody:       `{"links":""}`,
			responseBody:      ``,
			expectedStatus:    http.StatusBadRequest,
			expectedErrorCode: consts.ERROR_CODE_BAD_REQUEST,

			mockInputData:  entity.LinkGetStatus_Params{Links: []string{}},
			mockReturnData: entity.LinkGetStatus_Result{},
			mockReturnErr:  nil,
		},
		{
			name:              "Internal Server error: bad sequence of operations",
			requestBody:       `{"links":["aaa.com"]}`,
			responseBody:      ``,
			expectedStatus:    http.StatusInternalServerError,
			expectedErrorCode: consts.ERROR_CODE_INTERNAL_SERVER_ERROR,

			mockInputData: entity.LinkGetStatus_Params{Links: []string{"aaa.com"}},
			mockReturnData: entity.LinkGetStatus_Result{
				LinkStates: []entity.LinkState{
					{
						Link:        "aaa.com",
						IsAvailable: true,
					},
				},
				LinkNum: 1,
			},
			mockReturnErr: ErrBadOperationsSequence,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test request and responseReader
			reqBody := strings.NewReader(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", reqBody)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Setup echo context
			echo_inst := echo.New()
			echo_inst.Validator = validator
			echo_context := echo_inst.NewContext(req, rr)

			// If we expect DB call
			if tt.expectedStatus != http.StatusBadRequest {
				mockService.On("GetStatus", mock.Anything, tt.mockInputData).
					Return(tt.mockReturnData, tt.mockReturnErr).Once()
			}

			err := handler.GetStatus(echo_context)
			assert.Nil(t, err)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Handler returned wrong status code")

			mockService.AssertExpectations(t)

			// check response Body
			if tt.expectedStatus == http.StatusOK {
				var rawJSON string
				rawJSON, err = getRawJSON(rr.Body)
				assert.Nil(t, err)

				rawJSON = strings.TrimSpace(rawJSON)
				assert.Equal(t, tt.responseBody, rawJSON)
			} else {
				var errorRes dto.ErrorResponse
				err = json.NewDecoder(rr.Body).Decode(&errorRes)
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedErrorCode, errorRes.Code)
			}
		})
	}
}

func getRawJSON(body io.Reader) (string, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	rawJSON := string(bodyBytes)
	return rawJSON, nil
}
