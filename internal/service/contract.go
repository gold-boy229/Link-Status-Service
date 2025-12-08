package service

import (
	"context"
)

type linkRepository interface {
	GetLinkNum(ctx context.Context, links []string) (linkNum int, isNew bool, err error)
	GetLinksByLinkNum(ctx context.Context, linkNum int) (links []string, err error)
	StoreLinks(ctx context.Context, links []string, linkNum int) error
}

type linkChecker interface {
	IsLinkAvailable(ctx context.Context, link string) (bool, error)
}
