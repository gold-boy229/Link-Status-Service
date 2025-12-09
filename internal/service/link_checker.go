package service

import (
	"context"
	"fmt"
	"net/url"
	"sync"
)

type hTTPLinkChecker struct {
	checker linkChecker
}

func NewHTTPLinkChecker(clientChecker linkChecker) *hTTPLinkChecker {
	return &hTTPLinkChecker{checker: clientChecker}
}

func (h *hTTPLinkChecker) IsLinkAvailable(ctx context.Context, link string) (bool, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return false, fmt.Errorf("link %q cannot be parsed: %w", link, err)
	}

	if parsedURL.IsAbs() {
		return h.checker.IsLinkAvailable(ctx, parsedURL.String())
	}

	schemes := []string{"http", "https"}
	callsNumber := len(schemes)
	urls := make([]string, callsNumber)
	for idx := range schemes {
		u := *parsedURL
		u.Scheme = schemes[idx]
		urls[idx] = u.String()
	}

	wg := &sync.WaitGroup{}
	goodChan := make(chan bool, callsNumber)
	errChan := make(chan error, callsNumber)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			var isAvailable bool
			isAvailable, err = h.checker.IsLinkAvailable(ctx, url)
			if err != nil {
				select {
				case errChan <- err:
				case <-ctx.Done():
				}
				return
			}
			select {
			case goodChan <- isAvailable:
			case <-ctx.Done():
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(goodChan)
		close(errChan)
	}()

	goodResultCnt := 0
	errs := make([]error, 0, callsNumber)
	for {
		select {
		case isAvailable := <-goodChan:
			goodResultCnt++
			if isAvailable {
				return true, nil
			}
			if goodResultCnt+len(errs) == callsNumber {
				return false, nil
			}

		case err = <-errChan:
			errs = append(errs, err)
			if len(errs) == callsNumber {
				return false, fmt.Errorf("got two errors for http and https calls: \n1)%w \n2)%w", errs[0], errs[1])
			}
			if goodResultCnt > 0 {
				return false, nil
			}

		case <-ctx.Done():
			return false, ctx.Err()
		}
	}
}
