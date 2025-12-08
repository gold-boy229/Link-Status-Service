package mocks

import (
	"Link-Status-Service/internal/entity"

	"codeberg.org/go-pdf/fpdf"
	"github.com/stretchr/testify/mock"
)

type mockPDFBuilder struct {
	mock.Mock
}

func NewMockPDFBuilder() *mockPDFBuilder {
	return &mockPDFBuilder{}
}

func (m *mockPDFBuilder) BuildPDF(linkStatuses []entity.LinkStatus) *fpdf.Fpdf {
	args := m.Called(linkStatuses)
	return args.Get(0).(*fpdf.Fpdf)
}
