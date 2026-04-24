// Паттерн: worker pool с rate limit и retry. Задача показывает, как ограничивать частоту обработки, повторять неудачные запросы и собирать результаты нескольких воркеров.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type URLCheckResult struct {
	URL        string
	StatusCode int
	Error      error
	CheckedAt  time.Time
	Latency    time.Duration
}

func main() {
	requestTicker := time.NewTicker(100 * time.Millisecond)
	defer requestTicker.Stop()

	checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	urls := []string{
		"https://google.com",
		"https://nonexistent-site-12345.com",
		"https://github.com",
		"https://stackoverflow.com",
	}

	urlCh := make(chan string)
	resultCh := make(chan URLCheckResult)

	var workerGroup sync.WaitGroup
	for workerID := 1; workerID <= 5; workerID++ {
		workerGroup.Add(1)
		go checkURLWorker(checkCtx, workerID, requestTicker, urlCh, resultCh, &workerGroup)
	}

	go func() {
		for _, url := range urls {
			urlCh <- url
		}
		close(urlCh)
	}()

	go func() {
		workerGroup.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		fmt.Println(result)
	}
}

func checkURLWorker(
	ctx context.Context,
	workerID int,
	requestTicker *time.Ticker,
	urlCh <-chan string,
	resultCh chan<- URLCheckResult,
	workerGroup *sync.WaitGroup,
) {
	defer workerGroup.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Воркер", workerID, "остановлен")
			return
		case url, isOpen := <-urlCh:
			if !isOpen {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-requestTicker.C:
			}

			resultCh <- checkURLWithRetry(ctx, url)
		}
	}
}

func checkURLWithRetry(ctx context.Context, url string) URLCheckResult {
	retryCtx, cancelRetry := context.WithTimeout(ctx, 2*time.Second)
	defer cancelRetry()

	checkStartedAt := time.Now()
	statusCode, requestErr := fakeRequest(retryCtx, url)
	if requestErr == nil {
		return URLCheckResult{
			URL:        url,
			StatusCode: statusCode,
			Error:      nil,
			CheckedAt:  time.Now(),
			Latency:    time.Since(checkStartedAt),
		}
	}

	retryTicker := time.NewTicker(100 * time.Millisecond)
	defer retryTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return buildCheckResult(url, statusCode, ctx.Err(), checkStartedAt)
		case <-retryCtx.Done():
			return buildCheckResult(url, statusCode, requestErr, checkStartedAt)
		case <-retryTicker.C:
			statusCode, requestErr = fakeRequest(retryCtx, url)
			if requestErr == nil {
				return buildCheckResult(url, statusCode, nil, checkStartedAt)
			}
		}
	}
}

func buildCheckResult(url string, statusCode int, requestErr error, checkStartedAt time.Time) URLCheckResult {
	return URLCheckResult{
		URL:        url,
		StatusCode: statusCode,
		Error:      requestErr,
		CheckedAt:  time.Now(),
		Latency:    time.Since(checkStartedAt),
	}
}

func fakeRequest(ctx context.Context, url string) (int, error) {
	switch {
	case url == "https://nonexistent-site-12345.com":
		time.Sleep(100 * time.Millisecond)
		return 0, fmt.Errorf("connection refused")
	case url == "https://google.com":
		select {
		case <-time.After(1500 * time.Millisecond):
			return 200, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	default:
		select {
		case <-time.After(200 * time.Millisecond):
			return 200, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
}
