package repository

import (
	"context"
	"errors"
)

type linkRepository struct{}

func NewLinkRepository() *linkRepository {
	return &linkRepository{}
}

func (repo *linkRepository) StoreLinks(ctx context.Context, links []string, linkNum int) error {
	return errors.New("not implemented")
}

func (repo *linkRepository) GetLinkNum(ctx context.Context, links []string) (linkNum int, isNew bool, err error) {
	return linkNum, isNew, errors.New("not implemented")
}
