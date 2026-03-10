package webcrawler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/time/rate"
)

// ContentFetcher defines an interface for fetching content from URLs
type ContentFetcher interface {
	Fetch(ctx context.Context, url string) ([]byte, error)
}

// ContentProcessor defines an interface for processing raw content
type ContentProcessor interface {
	Process(ctx context.Context, content []byte) (ProcessedData, error)
}

// ProcessedData represents structured data extracted from raw content
type ProcessedData struct {
	Title       string
	Description string
	Keywords    []string
	Timestamp   time.Time
	Source      string
}

// ContentAggregator manages the concurrent fetching and processing of content.
//
// It uses two WaitGroups with different scopes:
//   - wg: Tracks all workers across the aggregator's lifetime. Shutdown() waits on
//     this to ensure every worker from every batch has exited before returning.
//   - batchWg (per fanOut call): Tracks workers for a single batch. The coordinator
//     waits on this to know when to close result channels so the receive loop can
//     terminate. Both are needed because batchWg answers "when are this batch's
//     workers done?" while wg answers "when are all workers done?"
type ContentAggregator struct {
	fetcher      ContentFetcher
	processor    ContentProcessor
	workerCount  int
	limiter      *rate.Limiter
	wg           sync.WaitGroup
	shutdown     chan struct{}
	shutdownOnce sync.Once
}

// NewContentAggregator creates a new ContentAggregator with the specified configuration.
// Returns nil if fetcher or processor is nil, or if workerCount or requestsPerSecond
// is zero or negative. Callers must check for nil before using the returned value.
func NewContentAggregator(
	fetcher ContentFetcher,
	processor ContentProcessor,
	workerCount int,
	requestsPerSecond int,
) *ContentAggregator {
	if fetcher == nil || processor == nil || workerCount <= 0 || requestsPerSecond <= 0 {
		return nil
	}
	return &ContentAggregator{
		fetcher:     fetcher,
		processor:   processor,
		workerCount: workerCount,
		limiter:     rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond),
		shutdown:    make(chan struct{}),
	}
}

// FetchAndProcess concurrently fetches and processes content from multiple URLs
func (ca *ContentAggregator) FetchAndProcess(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, error) {
	if len(urls) == 0 {
		return []ProcessedData{}, nil
	}
	results, errs := ca.fanOut(ctx, urls)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return results, nil
}

// Shutdown performs cleanup and ensures all resources are properly released
func (ca *ContentAggregator) Shutdown() error {
	ca.shutdownOnce.Do(func() {
		close(ca.shutdown)
		ca.wg.Wait()
	})
	return nil
}

// workerPool spawns workerCount goroutines and returns immediately.
// Each worker signals both ca.wg (for Shutdown) and batchWg (for per-batch coordination).
func (ca *ContentAggregator) workerPool(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errs chan<- error,
	batchWg *sync.WaitGroup,
) {
	for range ca.workerCount {
		ca.wg.Add(1)
		batchWg.Add(1)
		go func() {
			defer ca.wg.Done()
			defer batchWg.Done()
			ca.runWorker(ctx, jobs, results, errs)
		}()
	}
}

// runWorker processes jobs from the jobs channel until it is closed or the context/shutdown is triggered.
func (ca *ContentAggregator) runWorker(
	ctx context.Context,
	jobs <-chan string,
	results chan<- ProcessedData,
	errs chan<- error,
) {
	for {
		select {
		case url, ok := <-jobs:
			if !ok {
				return
			}
			data, err := ca.processURL(ctx, url)
			if err != nil {
				errs <- err
				continue
			}
			results <- data
		case <-ctx.Done():
			return
		case <-ca.shutdown:
			return
		}
	}
}

// processURL rate-limits, fetches, and processes a single URL.
func (ca *ContentAggregator) processURL(ctx context.Context, url string) (ProcessedData, error) {
	if err := ca.limiter.Wait(ctx); err != nil {
		return ProcessedData{}, err
	}
	content, err := ca.fetcher.Fetch(ctx, url)
	if err != nil {
		return ProcessedData{}, err
	}
	data, err := ca.processor.Process(ctx, content)
	if err != nil {
		return ProcessedData{}, err
	}
	data.Source = url
	return data, nil
}

// fanOut implements a fan-out, fan-in pattern for processing multiple items concurrently.
// A coordinator goroutine waits for all workers to exit, then closes the result channels
// so the receive loop can terminate without deadlocking when the producer exits early.
func (ca *ContentAggregator) fanOut(
	ctx context.Context,
	urls []string,
) ([]ProcessedData, []error) {
	n := len(urls)
	jobs := make(chan string, n)
	results := make(chan ProcessedData, n)
	errs := make(chan error, n)

	// batchWg: per-batch scope; coordinator waits on it to close result channels.
	var batchWg sync.WaitGroup
	ca.workerPool(ctx, jobs, results, errs, &batchWg)

	// Producer: send URLs, then close jobs so workers know to exit
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			case <-ca.shutdown:
				return
			}
		}
	}()

	// Coordinator: wait for this batch's workers to exit, then close result channels.
	// Uses batchWg (not ca.wg) so we only wait for this batch, not workers from other calls.
	go func() {
		batchWg.Wait()
		close(results)
		close(errs)
	}()

	var allResults []ProcessedData
	var allErrors []error
	for results != nil || errs != nil {
		select {
		case r, ok := <-results:
			if !ok {
				results = nil
			} else {
				allResults = append(allResults, r)
			}
		case e, ok := <-errs:
			if !ok {
				errs = nil
			} else {
				allErrors = append(allErrors, e)
			}
		case <-ctx.Done():
			return nil, []error{ctx.Err()}
		case <-ca.shutdown:
			return nil, []error{errors.New("shut down")}
		}
	}
	return allResults, allErrors
}

// HTTPFetcher is a simple implementation of ContentFetcher that uses HTTP
type HTTPFetcher struct {
	Client *http.Client
}

// Fetch retrieves content from a URL via HTTP
func (hf *HTTPFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := hf.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// HTMLProcessor is a basic implementation of ContentProcessor for HTML content
type HTMLProcessor struct{}

// Process extracts structured data from HTML content
func (hp *HTMLProcessor) Process(_ context.Context, content []byte) (ProcessedData, error) {
	if len(content) == 0 {
		return ProcessedData{}, errors.New("empty content")
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		return ProcessedData{}, err
	}

	title := doc.Find("title").First().Text()
	if title == "" {
		return ProcessedData{}, errors.New("missing title")
	}

	description, _ := doc.Find(`meta[name="description"]`).Attr("content")
	keywords, _ := doc.Find(`meta[name="keywords"]`).Attr("content")

	var kw []string
	for _, k := range strings.Split(keywords, ",") {
		if s := strings.TrimSpace(k); s != "" {
			kw = append(kw, s)
		}
	}

	return ProcessedData{
		Title:       title,
		Description: description,
		Keywords:    kw,
		Timestamp:   time.Now(),
	}, nil
}
