package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// Worker is a struct that represents a concurrent worker
type Worker struct {
	id               int          // worker id
	url              string       // url to make requests to
	rpsForEachWorker int          // Requests per second for each worker
	client           *http.Client // HTTP client to use
}

// Result is a struct that holds the result of a request
type Result struct {
	workerID int           // worker id
	status   int           // status code
	latency  time.Duration // latency
	err      error         // error if any
}

// NewWorker creates a new worker with the given parameters
func NewWorker(id int, url string, rpsForEachWorker int, client *http.Client) *Worker {
	return &Worker{
		id:               id,
		url:              url,
		rpsForEachWorker: rpsForEachWorker,
		client:           client,
	}
}

// Run runs the worker and sends the results to the given channel
func (w *Worker) Run(results chan<- Result, duration time.Duration) {
	defer func() {
		// handle panic gracefully
		if r := recover(); r != nil {
			fmt.Println("Worker", w.id, "panicked:", r)
		}
	}()

	// Calculate request rate per second for each worker
	requestRatePerSecond := w.rpsForEachWorker / int(duration)
	// Calculate sleep duration between each request
	sleepDuration := time.Second / time.Duration(requestRatePerSecond)

	// Loop for the duration of d seconds
	for j := 0; j < int(duration); j++ {
		// Loop to make requests at the desired rate
		for i := 0; i < requestRatePerSecond; i++ {
			// Make a GET request and measure the latency
			start := time.Now()
			resp, err := w.client.Get(w.url)
			latency := time.Since(start)

			// Send the result to the channel
			result := Result{w.id, 0, latency, err}
			if err == nil {
				// Close the response body and get the status code
				defer resp.Body.Close()
				result.status = resp.StatusCode
			}
			results <- result

			// Sleep for the calculated duration before making the next request
			time.Sleep(sleepDuration)
		}
	}
}

// main function
func main() {
	startTime := time.Now()

	// parse command line arguments
	var reqPerSec, duration, totalRequests int
	var url string
	flag.IntVar(&reqPerSec, "rps", 10, "requests per second")
	flag.IntVar(&duration, "dur", 10, "duration in seconds")
	flag.StringVar(&url, "url", "https://example.com", "url to make requests to")
	flag.Parse()

	// Calculate the total number of requests needed
	totalRequests = reqPerSec * duration

	// determine the number of workers based on the number of CPUs
	workers := runtime.NumCPU()

	// Calculate the number of requests per worker
	requestsPerWorker := totalRequests / workers

	// Create a channel for results
	results := make(chan Result, totalRequests)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a wait group for workers
	wg := &sync.WaitGroup{}
	wg.Add(workers)

	// Create a map to store status code metrics
	statusMetrics := make(map[int]*StatusCodeMetrics)

	// Create and run workers
	for i := 0; i < workers; i++ {
		worker := NewWorker(i, url, requestsPerWorker, client)
		go func() {
			worker.Run(results, time.Duration(duration))
			wg.Done()
		}()
	}

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	// Collect and print metrics
	var totalErrors int
	var minLatency, maxLatency, sumLatency int64
	minLatency = 1<<63 - 1 // max int64 value

	// Iterate over the results
	for result := range results {
		if result.err != nil {
			totalErrors++
		}
		sumLatency += int64(result.latency)
		if result.latency < time.Duration(minLatency) {
			minLatency = int64(result.latency)
		}
		if result.latency > time.Duration(maxLatency) {
			maxLatency = int64(result.latency)
		}

		// Update status code metrics
		if _, ok := statusMetrics[result.status]; !ok {
			statusMetrics[result.status] = &StatusCodeMetrics{
				Count:      0,
				MinLatency: 1<<63 - 1,
				MaxLatency: 0,
				SumLatency: 0,
			}
		}
		statusMetrics[result.status].Count++
		statusMetrics[result.status].SumLatency += result.latency
		if result.latency < time.Duration(statusMetrics[result.status].MinLatency) {
			statusMetrics[result.status].MinLatency = result.latency
		}
		if result.latency > time.Duration(statusMetrics[result.status].MaxLatency) {
			statusMetrics[result.status].MaxLatency = result.latency
		}
	}

	// Calculate average latency
	avgLatency := float64(sumLatency) / float64(totalRequests)

	// Print metrics
	fmt.Println("Total Number of Requests:", totalRequests)
	fmt.Println("Average Latency:", time.Duration(avgLatency))
	fmt.Println("Requests Per Second:", reqPerSec)
	fmt.Println("Min Latency:", time.Duration(minLatency))
	fmt.Println("Max Latency:", time.Duration(maxLatency))
	fmt.Println("Error Rate:", float64(totalErrors)/float64(totalRequests)*100, "%")
	fmt.Println("Status Code      Counts      Min Latency      Max Latency      Avg Latency")
	for status, metrics := range statusMetrics {
		fmt.Printf("%-16d%-12d%-17s%-17s%-17s\n", status, metrics.Count, metrics.MinLatency, metrics.MaxLatency, time.Duration(metrics.SumLatency.Nanoseconds()/int64(metrics.Count)))
	}
	fmt.Println("Total execution time", time.Since(startTime))
}

// Define a struct to store the status code metrics
type StatusCodeMetrics struct {
	Count      int // number of requests with this status code
	MinLatency time.Duration
	MaxLatency time.Duration
	SumLatency time.Duration // sum of latencies for this status code
}
