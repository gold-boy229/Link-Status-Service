package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreLinks(t *testing.T) {
	tests := []struct {
		name         string
		inputLinks   []string
		inputLinkNum int
		expectErr    bool
	}{
		{
			name:         "Success case: empty array",
			inputLinks:   []string{},
			inputLinkNum: 1,
			expectErr:    false,
		},
		{
			name: "Success case: not empty array",
			inputLinks: []string{
				"aaa.com",
				"bbb.ru",
			},
			inputLinkNum: 1,
			expectErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := NewLinkRepository()

			// store
			err := repo.StoreLinks(context.Background(), tt.inputLinks, tt.inputLinkNum)
			if tt.expectErr {
				assert.NotNil(t, err)
			}
			assert.Nil(t, err)

			// read
			resultLinks, err := repo.GetLinksByLinkNum(context.Background(), tt.inputLinkNum)
			if tt.expectErr {
				assert.NotNil(t, err)
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.inputLinks, resultLinks)
		})
	}
}
