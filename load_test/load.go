package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(url string, results chan<- time.Duration, errors chan<- error, RequestPerWorker int, wg *sync.WaitGroup) {
	defer wg.Done()

	startTime := time.Now()
	for i := 0; i < RequestPerWorker; i++ {
		resp, err := http.Get(url)
		if err != nil {
			errors <- err
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}

	elapsed := time.Since(startTime)
	results <- elapsed
}

func main() {
	// URL to make GET requests to
	url := "http://localhost:9011/books?bookId=1"

	// Number of concurrent workers
	numWorkers := 5
	//Number of Request per workers
	RequestPerWorker := numWorkers * 10

	// Duration for which workers should run
	// duration := 5 * time.Second

	// Channels to collect results and errors
	results := make(chan time.Duration, RequestPerWorker)
	errors := make(chan error, RequestPerWorker)

	// Channel to signal when all workers have completed
	// done := make(chan struct{}, numWorkers)

	// WaitGroup to wait for all workers to finish
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start workers
	startTime := time.Now()
	for i := 0; i < numWorkers; i++ {
		go worker(url, results, errors, RequestPerWorker, &wg)
	}

	// Aggregate metrics
	var totalResponseTime time.Duration
	var numRequests int
	var minLatency, maxLatency time.Duration = time.Hour, 0
	var errorCount int

	// for time.Since(startTime) < duration {
	// 	select {
	// 	case result := <-results:
	// 		// Handle result (response time) received from results channel
	// 		totalResponseTime += result
	// 		if result < minLatency {
	// 			minLatency = result
	// 		}
	// 		if result > maxLatency {
	// 			maxLatency = result
	// 		}
	// 		numRequests++
	// 		fmt.Printf("Received result: %v\n", result)
	// 	case err := <-errors:
	// 		// Handle error received from errors channel
	// 		errorCount++
	// 		fmt.Printf("Received error: %v\n", err)
	// 	}
	// }

	for result := range results {
		totalResponseTime += result
		if result < minLatency {
			minLatency = result
		}
		if result > maxLatency {
			maxLatency = result
		}
		numRequests++
	}
	for range errors {
		errorCount++
	}

	// Wait for all workers to finish
	wg.Wait()

	// Close the channels
	close(results)
	close(errors)

	avgLatency := totalResponseTime / time.Duration(numRequests)
	requestsPerSecond := float64(numRequests) / time.Since(startTime).Seconds()

	// Print metrics
	fmt.Printf("Total Number of Requests %v\n", numRequests)
	fmt.Printf("Average Latency: %v\n", avgLatency)
	fmt.Printf("Min Latency: %v\n", minLatency)
	fmt.Printf("Max Latency: %v\n", maxLatency)
	fmt.Printf("Requests Per Second: %f\n", requestsPerSecond)
	fmt.Printf("Error Rate: %.2f%%\n", float64(errorCount)/float64(numRequests)*100)
}
