package webcrawler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// MockFetcher is a mock implementation of ContentFetcher for testing
type MockFetcher struct {
	responses map[string][]byte
	errors    map[string]error
	callCount map[string]int
}

func NewMockFetcher() *MockFetcher {
	return &MockFetcher{
		responses: make(map[string][]byte),
		errors:    make(map[string]error),
		callCount: make(map[string]int),
	}
}

func (m *MockFetcher) SetResponse(url string, response []byte) {
	m.responses[url] = response
}

func (m *MockFetcher) SetError(url string, err error) {
	m.errors[url] = err
}

func (m *MockFetcher) GetCallCount(url string) int {
	return m.callCount[url]
}

func (m *MockFetcher) Fetch(ctx context.Context, url string) ([]byte, error) {
	m.callCount[url]++

	if err, ok := m.errors[url]; ok {
		return nil, err
	}

	if response, ok := m.responses[url]; ok {
		return response, nil
	}

	return nil, errors.New("no mock response for URL")
}

// MockProcessor is a mock implementation of ContentProcessor for testing
type MockProcessor struct {
	results   map[string]ProcessedData
	errors    map[string]error
	callCount int
}

func NewMockProcessor() *MockProcessor {
	return &MockProcessor{
		results:   make(map[string]ProcessedData),
		errors:    make(map[string]error),
		callCount: 0,
	}
}

func (m *MockProcessor) SetResult(content string, result ProcessedData) {
	m.results[content] = result
}

func (m *MockProcessor) SetError(content string, err error) {
	m.errors[content] = err
}

func (m *MockProcessor) GetCallCount() int {
	return m.callCount
}

func (m *MockProcessor) Process(ctx context.Context, content []byte) (ProcessedData, error) {
	m.callCount++
	contentStr := string(content)

	if err, ok := m.errors[contentStr]; ok {
		return ProcessedData{}, err
	}

	if result, ok := m.results[contentStr]; ok {
		return result, nil
	}

	return ProcessedData{}, errors.New("no mock result for content")
}

// TestNewContentAggregator tests the constructor function
func TestNewContentAggregator(t *testing.T) {
	fetcher := NewMockFetcher()
	processor := NewMockProcessor()

	tests := []struct {
		name              string
		fetcher           ContentFetcher
		processor         ContentProcessor
		workerCount       int
		requestsPerSecond int
		expectNil         bool
	}{
		{
			name:              "Valid configuration",
			fetcher:           fetcher,
			processor:         processor,
			workerCount:       5,
			requestsPerSecond: 10,
			expectNil:         false,
		},
		{
			name:              "No fetcher",
			fetcher:           nil,
			processor:         processor,
			workerCount:       5,
			requestsPerSecond: 10,
			expectNil:         true,
		},
		{
			name:              "No processor",
			fetcher:           fetcher,
			processor:         nil,
			workerCount:       5,
			requestsPerSecond: 10,
			expectNil:         true,
		},
		{
			name:              "Zero workers",
			fetcher:           fetcher,
			processor:         processor,
			workerCount:       0,
			requestsPerSecond: 10,
			expectNil:         true,
		},
		{
			name:              "Negative workers",
			fetcher:           fetcher,
			processor:         processor,
			workerCount:       -1,
			requestsPerSecond: 10,
			expectNil:         true,
		},
		{
			name:              "Zero rate limit",
			fetcher:           fetcher,
			processor:         processor,
			workerCount:       5,
			requestsPerSecond: 0,
			expectNil:         true,
		},
		{
			name:              "Negative rate limit",
			fetcher:           fetcher,
			processor:         processor,
			workerCount:       5,
			requestsPerSecond: -1,
			expectNil:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregator := NewContentAggregator(
				tt.fetcher,
				tt.processor,
				tt.workerCount,
				tt.requestsPerSecond,
			)

			if tt.expectNil && aggregator != nil {
				t.Errorf("Expected nil aggregator, but got non-nil")
			}

			if !tt.expectNil && aggregator == nil {
				t.Errorf("Expected non-nil aggregator, but got nil")
			}
		})
	}
}

// TestFetchAndProcess tests the concurrent fetching and processing functionality
func TestFetchAndProcess(t *testing.T) {
	ctx := context.Background()

	// Setup mock data
	htmlContent1 := []byte("<html><title>Test Page 1</title></html>")
	htmlContent2 := []byte("<html><title>Test Page 2</title></html>")

	expectedData1 := ProcessedData{
		Title:       "Test Page 1",
		Description: "Description 1",
		Keywords:    []string{"test", "page1"},
		Timestamp:   time.Now(),
		Source:      "https://example.com/1",
	}

	expectedData2 := ProcessedData{
		Title:       "Test Page 2",
		Description: "Description 2",
		Keywords:    []string{"test", "page2"},
		Timestamp:   time.Now(),
		Source:      "https://example.com/2",
	}

	// Setup mocks
	fetcher := NewMockFetcher()
	fetcher.SetResponse("https://example.com/1", htmlContent1)
	fetcher.SetResponse("https://example.com/2", htmlContent2)
	fetcher.SetError("https://example.com/error", errors.New("fetch error"))

	processor := NewMockProcessor()
	processor.SetResult(string(htmlContent1), expectedData1)
	processor.SetResult(string(htmlContent2), expectedData2)
	processor.SetError(string([]byte("error content")), errors.New("processing error"))

	aggregator := NewContentAggregator(fetcher, processor, 3, 10)
	if aggregator == nil {
		t.Fatal("Failed to create ContentAggregator")
	}

	tests := []struct {
		name     string
		urls     []string
		expected []ProcessedData
		wantErr  bool
	}{
		{
			name:     "Single URL",
			urls:     []string{"https://example.com/1"},
			expected: []ProcessedData{expectedData1},
			wantErr:  false,
		},
		{
			name:     "Multiple URLs",
			urls:     []string{"https://example.com/1", "https://example.com/2"},
			expected: []ProcessedData{expectedData1, expectedData2},
			wantErr:  false,
		},
		{
			name:     "Error URL",
			urls:     []string{"https://example.com/error"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Mixed success and error",
			urls:     []string{"https://example.com/1", "https://example.com/error"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Empty URL list",
			urls:     []string{},
			expected: []ProcessedData{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := aggregator.FetchAndProcess(ctx, tt.urls)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.wantErr {
				if len(results) != len(tt.expected) {
					t.Errorf("Expected %d results, got %d", len(tt.expected), len(results))
				}

				// Check that all expected results are present (order may vary due to concurrency)
				if len(tt.expected) > 0 {
					for _, expected := range tt.expected {
						found := false
						for _, result := range results {
							if result.Title == expected.Title && result.Source == expected.Source {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Expected result with title '%s' not found", expected.Title)
						}
					}
				}
			}
		})
	}
}

// TestShutdown tests proper resource cleanup
func TestShutdown(t *testing.T) {
	fetcher := NewMockFetcher()
	processor := NewMockProcessor()

	aggregator := NewContentAggregator(fetcher, processor, 3, 10)
	if aggregator == nil {
		t.Fatal("Failed to create ContentAggregator")
	}

	err := aggregator.Shutdown()
	if err != nil {
		t.Errorf("Shutdown returned error: %v", err)
	}

	// Test double shutdown should not cause errors
	err = aggregator.Shutdown()
	if err != nil {
		t.Errorf("Second shutdown call returned error: %v", err)
	}
}

// TestHTTPFetcher tests the HTTP implementation of ContentFetcher
func TestHTTPFetcher(t *testing.T) {
	// Create a mock HTTP client
	mockClient := &http.Client{
		Transport: &mockTransport{
			responses: map[string]*http.Response{
				"https://example.com": {
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString("test content")),
				},
				"https://example.com/404": {
					StatusCode: 404,
					Body:       io.NopCloser(bytes.NewBufferString("not found")),
				},
			},
		},
	}

	fetcher := &HTTPFetcher{
		Client: mockClient,
	}

	tests := []struct {
		name    string
		url     string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Successful fetch",
			url:     "https://example.com",
			want:    []byte("test content"),
			wantErr: false,
		},
		{
			name:    "404 response",
			url:     "https://example.com/404",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			url:     "https://",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := fetcher.Fetch(ctx, tt.url)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("Expected %s, got %s", string(tt.want), string(got))
			}
		})
	}
}

// TestHTMLProcessor tests the HTML implementation of ContentProcessor
func TestHTMLProcessor(t *testing.T) {
	processor := &HTMLProcessor{}

	tests := []struct {
		name    string
		content []byte
		want    ProcessedData
		wantErr bool
	}{
		{
			name: "Valid HTML",
			content: []byte(
				"<html><head><title>Test Page</title><meta name=\"description\" content=\"Test Description\"><meta name=\"keywords\" content=\"test,keywords\"></head><body>Content</body></html>",
			),
			want: ProcessedData{
				Title:       "Test Page",
				Description: "Test Description",
				Keywords:    []string{"test", "keywords"},
				Source:      "",
			},
			wantErr: false,
		},
		{
			name:    "Invalid HTML",
			content: []byte("<malformed>"),
			want:    ProcessedData{},
			wantErr: true,
		},
		{
			name:    "Empty HTML",
			content: []byte(""),
			want:    ProcessedData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := processor.Process(ctx, tt.content)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.wantErr {
				if got.Title != tt.want.Title {
					t.Errorf("Expected title %s, got %s", tt.want.Title, got.Title)
				}

				if got.Description != tt.want.Description {
					t.Errorf(
						"Expected description %s, got %s",
						tt.want.Description,
						got.Description,
					)
				}

				if !reflect.DeepEqual(got.Keywords, tt.want.Keywords) {
					t.Errorf("Expected keywords %v, got %v", tt.want.Keywords, got.Keywords)
				}
			}
		})
	}
}

// mockTransport implements the RoundTripper interface for testing HTTP requests
type mockTransport struct {
	responses map[string]*http.Response
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response, ok := m.responses[req.URL.String()]
	if !ok {
		return nil, errors.New("no mock response for URL")
	}
	return response, nil
}
