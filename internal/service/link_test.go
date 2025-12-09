package service

import (
	"context"
	"errors"
	"testing"

	"Link-Status-Service/internal/entity"
	"Link-Status-Service/internal/mocks"
	"Link-Status-Service/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testCheckerOut struct {
	isAvailable bool
	err         error
}

func TestGetLinkStates(t *testing.T) {
	tests := []struct {
		name             string
		mockCheckerInput []string
		mockCheckerOut   []testCheckerOut
		expectedResult   []entity.LinkState
		expectErr        bool
	}{
		{
			name:             "Success case: zero links",
			mockCheckerInput: []string{},
			mockCheckerOut:   []testCheckerOut{},
			expectedResult:   []entity.LinkState{},
			expectErr:        false,
		},
		{
			name: "Success case: one link",
			mockCheckerInput: []string{
				"aaa.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "aaa.com",
					IsAvailable: true,
				},
			},
			expectErr: false,
		},
		{
			name: "Success case: many links(a,b,c)",
			mockCheckerInput: []string{
				"aaa.com",
				"bbb.com",
				"ccc.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: false,
					err:         nil,
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "aaa.com",
					IsAvailable: true,
				},
				{
					Link:        "bbb.com",
					IsAvailable: true,
				},
				{
					Link:        "ccc.com",
					IsAvailable: false,
				},
			},
			expectErr: false,
		},
		{
			name: "Success case: many links(b,c,a)",
			mockCheckerInput: []string{
				"bbb.com",
				"ccc.com",
				"aaa.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: false,
					err:         nil,
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "bbb.com",
					IsAvailable: true,
				},
				{
					Link:        "ccc.com",
					IsAvailable: true,
				},
				{
					Link:        "aaa.com",
					IsAvailable: false,
				},
			},
			expectErr: false,
		},
		{
			name: "Success case: many links(c,a,b)",
			mockCheckerInput: []string{
				"ccc.com",
				"aaa.com",
				"bbb.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: false,
					err:         nil,
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "ccc.com",
					IsAvailable: true,
				},
				{
					Link:        "aaa.com",
					IsAvailable: true,
				},
				{
					Link:        "bbb.com",
					IsAvailable: false,
				},
			},
			expectErr: false,
		},
		{
			name: "Failure case: one error",
			mockCheckerInput: []string{
				"good.com",
				"bad.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: false,
					err:         errors.New("connection failed"),
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "good.com",
					IsAvailable: true,
				},
				{
					Link:        "bad.com",
					IsAvailable: false,
				},
			},
			expectErr: false,
		},
		{
			name: "Failure case: several errors",
			mockCheckerInput: []string{
				"good_1.com",
				"bad_1.com",
				"bad_2.com",
				"good_2.com",
			},
			mockCheckerOut: []testCheckerOut{
				{
					isAvailable: true,
					err:         nil,
				},
				{
					isAvailable: false,
					err:         errors.New("connection failed"),
				},
				{
					isAvailable: false,
					err:         errors.New("connection failed"),
				},
				{
					isAvailable: true,
					err:         nil,
				},
			},
			expectedResult: []entity.LinkState{
				{
					Link:        "good_1.com",
					IsAvailable: true,
				},
				{
					Link:        "bad_1.com",
					IsAvailable: false,
				},
				{
					Link:        "bad_2.com",
					IsAvailable: false,
				},
				{
					Link:        "good_2.com",
					IsAvailable: true,
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockLinkRepository()
			mockChecker := mocks.NewMockLinkChecker()
			service := NewLinkService(mockRepo, mockChecker)

			require.Equal(t, len(tt.mockCheckerInput), len(tt.mockCheckerOut),
				"Input/Out checker slices must have equal length")
			for idx := range tt.mockCheckerInput {
				mockChecker.On("IsLinkAvailable", mock.Anything, tt.mockCheckerInput[idx]).
					Return(tt.mockCheckerOut[idx].isAvailable, tt.mockCheckerOut[idx].err)
			}

			linkStates, err := service.getLinkStates(context.Background(), tt.mockCheckerInput)

			if tt.expectErr {
				assert.NotNil(t, err)
				return
			}

			mockChecker.AssertExpectations(t)

			assert.Nil(t, err)
			assert.EqualValues(t, tt.expectedResult, linkStates)
		})
	}
}

func TestGetUniqueLinksFromLinkSets(t *testing.T) {
	type mockOutput struct {
		ReturnData []string
		ReturnErr  error
	}

	tests := []struct {
		name                 string
		inputLinkNums        []int
		mockOutputs          []mockOutput
		expectedResultSorted []string
		expectErr            bool
	}{
		{
			name:                 "Success case: no linkNums",
			inputLinkNums:        []int{},
			mockOutputs:          []mockOutput{},
			expectedResultSorted: []string{},
			expectErr:            false,
		},
		{
			name:          "Success case: one linkNum",
			inputLinkNums: []int{1},
			mockOutputs: []mockOutput{
				{
					ReturnData: []string{
						"aaa.com",
						"bbb.ru",
					},
					ReturnErr: nil,
				},
			},
			expectedResultSorted: []string{
				"aaa.com",
				"bbb.ru",
			},
			expectErr: false,
		},
		{
			name:          "Success case: union two not intersecting sets",
			inputLinkNums: []int{1, 2},
			mockOutputs: []mockOutput{
				{
					ReturnData: []string{
						"b_second.com",
						"c_third.com",
					},
					ReturnErr: nil,
				},
				{
					ReturnData: []string{
						"a_first.com",
						"d_fourth.com",
					},
					ReturnErr: nil,
				},
			},
			expectedResultSorted: []string{
				"a_first.com",
				"b_second.com",
				"c_third.com",
				"d_fourth.com",
			},
			expectErr: false,
		},
		{
			name:          "Success case: union two intersecting sets",
			inputLinkNums: []int{1, 2},
			mockOutputs: []mockOutput{
				{
					ReturnData: []string{
						"b_second.com",
						"a_first.com",
					},
					ReturnErr: nil,
				},
				{
					ReturnData: []string{
						"a_first.com",
						"c_third.com",
					},
					ReturnErr: nil,
				},
			},
			expectedResultSorted: []string{
				"a_first.com",
				"b_second.com",
				"c_third.com",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockLinkRepository()
			mockChecker := mocks.NewMockLinkChecker()
			service := NewLinkService(mockRepo, mockChecker)

			require.Len(t, tt.inputLinkNums, len(tt.mockOutputs))
			for idx, mockOutput := range tt.mockOutputs {
				mockRepo.On("GetLinksByLinkNum", mock.Anything, tt.inputLinkNums[idx]).
					Return(mockOutput.ReturnData, mockOutput.ReturnErr).Once()
			}

			uniqueLinks, err := service.getUniqueLinksFromLinkSets(context.Background(), tt.inputLinkNums)
			if tt.expectErr {
				assert.NotNil(t, err)
			}
			assert.Nil(t, err)

			sortedLinks := utils.SortStrings(uniqueLinks)
			assert.EqualValues(t, tt.expectedResultSorted, sortedLinks)

			mockRepo.AssertExpectations(t)
		})
	}
}
