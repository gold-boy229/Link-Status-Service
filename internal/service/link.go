package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"Link-Status-Service/internal/entity"
	"Link-Status-Service/internal/utils"
)

type linkService struct {
	repo    linkRepository
	checker linkChecker
}

func NewLinkService(repo linkRepository, checker linkChecker) *linkService {
	return &linkService{repo: repo, checker: checker}
}

func (s *linkService) GetStatus(
	ctx context.Context,
	params entity.LinkGetStatusParams,
) (entity.LinkGetStatusResult, error) {
	linkNum, isLinkNumNew, err := s.repo.GetLinkNum(ctx, params.Links)
	if err != nil {
		return entity.LinkGetStatusResult{}, fmt.Errorf("error during getting LinkNum: %w", err)
	}
	if isLinkNumNew {
		if err = s.repo.StoreLinks(ctx, params.Links, linkNum); err != nil {
			return entity.LinkGetStatusResult{}, fmt.Errorf("error during storing set of links: %w", err)
		}
	}

	linkStates, err := s.getLinkStates(ctx, params.Links)
	if err != nil {
		return entity.LinkGetStatusResult{}, fmt.Errorf("error during getting Link states: %w", err)
	}

	return entity.LinkGetStatusResult{
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

	wg := &sync.WaitGroup{}
	wg.Add(numLinks)

	for idx, link := range links {
		go func(idx int, link string) {
			defer wg.Done()

			isAvailable, _ := s.checker.IsLinkAvailable(ctx, link)
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
	}()

	linkStates := make([]entity.LinkState, numLinks)
	for i := 0; i < numLinks; i++ {
		select {
		case result := <-resultChan:
			linkStates[result.Index] = result.LinkState
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return linkStates, nil
}

func (s *linkService) GetStatusesOfLinkSets(
	ctx context.Context,
	params entity.LinkBuildPDSParams,
) (entity.LinkBuildPDSResult, error) {
	uniqueLinks, err := s.getUniqueLinksFromLinkSets(ctx, params.LinkNums)
	if err != nil {
		return entity.LinkBuildPDSResult{}, fmt.Errorf("cannot get uniqueLinks: %w", err)
	}

	sortedLinks := utils.SortStrings(uniqueLinks)
	linkStates, err := s.getLinkStates(ctx, sortedLinks)
	if err != nil {
		return entity.LinkBuildPDSResult{}, fmt.Errorf("cannot get linkStates: %w", err)
	}
	return entity.LinkBuildPDSResult{LinkStates: linkStates}, nil
}

func (s *linkService) getUniqueLinksFromLinkSets(ctx context.Context, linkNums []int) ([]string, error) {
	var (
		mp      = new(sync.Map)
		mpSize  = new(atomic.Int64)
		wg      = new(sync.WaitGroup)
		errChan = make(chan error, 1)
	)

	wg.Add(len(linkNums))
	for _, linkNum := range linkNums {
		go func(linkNum int) {
			defer wg.Done()

			links, err := s.repo.GetLinksByLinkNum(ctx, linkNum)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("cannot get links by linkNum[%v]: %w", linkNum, err):
				default:
				}
				return
			}

			for _, link := range links {
				_, alreadyExisted := mp.LoadOrStore(link, struct{}{})
				if !alreadyExisted {
					mpSize.Add(1)
				}
			}
		}(linkNum)
	}
	wg.Wait()

	select {
	case err := <-errChan:
		return []string{}, err
	default:
	}

	uniqueLinks := make([]string, 0, mpSize.Load())
	mp.Range(func(key, value any) bool {
		link, _ := key.(string)
		uniqueLinks = append(uniqueLinks, link)
		return true
	})
	return uniqueLinks, nil
}
