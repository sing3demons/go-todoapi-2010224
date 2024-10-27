package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

type service struct {
	maxRetries     int
	retiesDelay    time.Duration
	max            int
	maxConcurrency int
}

func NewHttp(max int) *service {
	if max <= 0 {
		max = 1
	}
	return &service{
		maxRetries:  3,
		retiesDelay: 500 * time.Millisecond,
		max:         max,
	}
}

func (s *service) Call(urls string) ([]byte, error) {
	maxRetries := s.maxRetries
	retryDelay := s.retiesDelay

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(
			context.Background(),
			http.MethodGet,
			urls,
			nil,
		)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if attempt == maxRetries {
				errorMsg := err.Error()
				if err == nil {
					errorMsg = "Status " + resp.Status
				}
				return nil, fmt.Errorf("%s", errorMsg)

			}
			time.Sleep(retryDelay * time.Duration(attempt))
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed to get response")
}

func (s *service) SetMaxConcurrency(maxConcurrency int) {
	s.maxConcurrency = maxConcurrency
}

func (s *service) SetMaxRetries(maxRetries int) {
	s.maxRetries = maxRetries
}

func (s *service) SetRetiesDelay(retriesDelay time.Duration) {
	s.retiesDelay = retriesDelay
}

func (s *service) NewSemaphore() (semaphore chan struct{}) {
	if s.max <= 0 {
		s.max = 1
	}
	if s.maxConcurrency <= 0 {
		s.maxConcurrency = s.max
	}

	if s.maxConcurrency > s.max {
		s.maxConcurrency = s.max
	}
	maxConcurrency := s.maxConcurrency
	fmt.Println("maxConcurrency", maxConcurrency)
	return make(chan struct{}, maxConcurrency)
}

func Api(c *gin.Context) {
	const max = 10000
	const maxConcurrency = 1000
	var url = os.Getenv("API_URL")
	if url == "" {
		url = "http://localhost:8080/transfer/"
	}

	start := time.Now()
	var wg sync.WaitGroup
	// semaphore := make(chan struct{}, maxConcurrency)   // Semaphore channel for concurrency control
	h := NewHttp(max)
	h.maxConcurrency = maxConcurrency
	semaphore := h.NewSemaphore()

	successCh := make(chan responseMessage, max)
	failCh := make(chan failedRequest, max)

	// Launch goroutines for requests with concurrency control
	for id := 1; id <= max; id++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore

			result, err := h.Call(id)
			if err != nil {
				failCh <- failedRequest{ID: id, Error: err.Error()}
				return
			}

			var response responseMessage
			if err := json.Unmarshal(result, &response); err != nil {
				failCh <- failedRequest{ID: id, Error: err.Error()}
				return
			}

			successCh <- response
			// makeRequest(id, successCh, failCh)
			<-semaphore // Release semaphore
		}(url + strconv.Itoa(id))
	}

	// Close channels once all goroutines are done
	go func() {
		wg.Wait()
		close(successCh)
		close(failCh)
	}()

	// Collect responses
	var failCount uint
	var result []responseMessage
	var failures []failedRequest

	for response := range successCh {
		result = append(result, response)
	}

	for failure := range failCh {
		failures = append(failures, failure)
		failCount++
	}

	// Log completion and time taken
	slog.Info("finish",
		slog.String("elapsed", time.Since(start).String()),
		slog.Int("success_count", len(result)),
		slog.Any("fail_count", failCount))

	// sort result asc
	sort.Slice(result, func(i, j int) bool {
		a, _ := strconv.Atoi(result[i].ID)
		b, _ := strconv.Atoi(result[j].ID)
		return a > b

	})
	// Return the result as JSON
	c.JSON(http.StatusOK, gin.H{
		"success_count": len(result),
		"fail_count":    failCount,
		"responses":     result[:100],
		"errors":        failures,
		"elapsed":       time.Since(start).String(),
	})
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", Api)

	r.Run(":8081")
}

type responseMessage struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type failedRequest struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

// makeRequest sends an HTTP request and categorizes the result as success or failure
// func makeRequest(id int, successCh chan responseMessage, failCh chan failedRequest) {
// 	const maxRetries = 3
// 	const retryDelay = 500 * time.Millisecond

// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		req, err := http.NewRequestWithContext(
// 			context.Background(),
// 			http.MethodGet,
// 			"http://localhost:8080/transfer/"+strconv.Itoa(id),
// 			nil,
// 		)
// 		if err != nil {
// 			failCh <- failedRequest{ID: id, Error: err.Error()}
// 			return
// 		}

// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil || resp.StatusCode != http.StatusOK {
// 			if attempt == maxRetries {
// 				errorMsg := err.Error()
// 				if err == nil {
// 					errorMsg = "Status " + resp.Status
// 				}
// 				failCh <- failedRequest{ID: id, Error: errorMsg}
// 			}
// 			time.Sleep(retryDelay * time.Duration(attempt))
// 			continue
// 		}
// 		defer resp.Body.Close()

// 		var response responseMessage
// 		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 			if attempt == maxRetries {
// 				failCh <- failedRequest{ID: id, Error: err.Error()}
// 			}
// 			time.Sleep(retryDelay * time.Duration(attempt))
// 			continue
// 		}

// 		successCh <- response
// 		return
// 	}
// }
