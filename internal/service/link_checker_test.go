package service

import (
	"Link-Status-Service/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsLinkAvailableTableDriven(t *testing.T) {
	type mockExpectation struct {
		URL         string
		ReturnAvail bool
		ReturnErr   error
		CalledTimes int
	}

	tests := []struct {
		name          string
		inputLink     string
		expectations  []mockExpectation
		wantAvailable bool
		wantErr       bool
		errContains   string
	}{
		{
			name:      "Absolute URL returns true",
			inputLink: "https://example.com/page",
			expectations: []mockExpectation{
				{URL: "https://example.com/page", ReturnAvail: true, ReturnErr: nil, CalledTimes: 1},
			},
			wantAvailable: true,
			wantErr:       false,
		},
		{
			name:      "Absolute URL returns false (not available)",
			inputLink: "https://example.com/page",
			expectations: []mockExpectation{
				{URL: "https://example.com/page", ReturnAvail: false, ReturnErr: nil, CalledTimes: 1},
			},
			wantAvailable: false,
			wantErr:       false,
		},
		{
			name:      "Relative URL, HTTP succeeds first (true)",
			inputLink: "/mypage",
			expectations: []mockExpectation{
				{URL: "https:///mypage", ReturnAvail: true, ReturnErr: nil, CalledTimes: 1},
				{URL: "http:///mypage", ReturnAvail: false, ReturnErr: nil, CalledTimes: 1},
			},
			wantAvailable: true,
			wantErr:       false,
		},
		{
			name:      "Relative URL, Both calls return false (no errors)",
			inputLink: "/mypage",
			expectations: []mockExpectation{
				{URL: "http:///mypage", ReturnAvail: false, ReturnErr: nil, CalledTimes: 1},
				{URL: "https:///mypage", ReturnAvail: false, ReturnErr: nil, CalledTimes: 1},
			},
			wantAvailable: false,
			wantErr:       false,
		},
		{
			name:      "Relative URL, Both calls error out",
			inputLink: "/broken",
			expectations: []mockExpectation{
				{URL: "http:///broken", ReturnAvail: false, ReturnErr: errors.New("http conn error"), CalledTimes: 1},
				{URL: "https:///broken", ReturnAvail: false, ReturnErr: errors.New("https conn error"), CalledTimes: 1},
			},
			wantAvailable: false,
			wantErr:       true,
			errContains:   "got two errors",
		},
		{
			name:          "Bad input link parsing fails early",
			inputLink:     ":invalid-url",
			expectations:  nil, // Моки не требуются, функция завершается рано
			wantAvailable: false,
			wantErr:       true,
			errContains:   "cannot be parsed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockLinkChecker()
			checker := NewHTTPLinkChecker(mockClient)

			for _, exp := range tt.expectations {
				mockClient.On("IsLinkAvailable", mock.Anything, exp.URL).
					Return(exp.ReturnAvail, exp.ReturnErr).
					Times(exp.CalledTimes)
			}

			gotAvailable, gotErr := checker.IsLinkAvailable(context.Background(), tt.inputLink)

			assert.Equal(t, tt.wantAvailable, gotAvailable, "availability mismatch")

			if tt.wantErr {
				assert.NotNil(t, gotErr, "expected an error but got nil")
				if tt.errContains != "" && gotErr != nil {
					assert.Contains(t, gotErr.Error(), tt.errContains)
				}
			} else {
				assert.Nil(t, gotErr, "expected no error but got one")
			}

			mockClient.AssertExpectations(t)
		})
	}
}
