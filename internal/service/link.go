package service

import (
	"Link-Status-Service/internal/entity"
	"Link-Status-Service/internal/utils"
	"context"
	"fmt"
	"sync"
)

type linkService struct {
	repo    linkRepository
	checker linkChecker
}

func NewLinkService(repo linkRepository, checker linkChecker) *linkService {
	return &linkService{repo: repo, checker: checker}
}

func (s *linkService) GetStatus(ctx context.Context, params entity.LinkGetStatus_Params) (entity.LinkGetStatus_Result, error) {
	links := utils.SortStrings(params.Links)

	linkNum, isLinkNumNew, err := s.repo.GetLinkNum(ctx, links)
	if err != nil {
		return entity.LinkGetStatus_Result{}, fmt.Errorf("error during getting LinkNum: %w", err)
	}
	if isLinkNumNew {
		if err := s.repo.StoreLinks(ctx, links, linkNum); err != nil {
			return entity.LinkGetStatus_Result{}, fmt.Errorf("error during storing set of links: %w", err)
		}
	}

	linkStates, err := s.getLinkStates(ctx, links)
	if err != nil {
		return entity.LinkGetStatus_Result{}, fmt.Errorf("error during getting Link states: %w", err)
	}

	return entity.LinkGetStatus_Result{
		LinkStates: linkStates,
		LinkNum:    linkNum,
	}, nil
}

type linkStateResult struct {
	entity.LinkState
	Index int
}

func (s *linkService) getLinkStates(ctx context.Context, links []string) ([]entity.LinkState, error) {
	numLinks := len(links)

	resultChan := make(chan linkStateResult, numLinks)
	errorChan := make(chan error, 1)

	wg := &sync.WaitGroup{}
	wg.Add(numLinks)

	for idx, link := range links {
		go func(idx int, link string) {
			defer wg.Done()

			isAvailable, err := s.checker.IsLinkAvailable(ctx, link)
			if err != nil {
				// Send error to error channel and return early
				select {
				case errorChan <- fmt.Errorf("link %s check failed: %w", link, err):
				default:
				}
				return
			}

			// Send successful result to result channel
			select {
			case resultChan <- linkStateResult{
				LinkState: entity.LinkState{
					Link:        link,
					IsAvailable: isAvailable,
				},
				Index: idx,
			}:
			default:
			}
		}(idx, link)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	linkStates := make([]entity.LinkState, numLinks)
	for i := 0; i < numLinks; i++ {
		select {
		case result := <-resultChan:
			linkStates[result.Index] = result.LinkState
		case err := <-errorChan:
			return nil, err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return linkStates, nil
}
