package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Worker is a struct that represents a concurrent worker
type Worker struct {
	id     int          // worker id
	url    string       // url to make requests to
	n      int          // number of requests to make
	client *http.Client // http client to use
}

// Result is a struct that holds the result of a request
type Result struct {
	workerID int           // worker id
	status   int           // status code
	latency  time.Duration // latency
	err      error         // error if any
}

// NewWorker creates a new worker with the given parameters
func NewWorker(id int, url string, n int, client *http.Client) *Worker {
	return &Worker{id, url, n, client}
}

// Run runs the worker and sends the results to the given channel
func (w *Worker) Run(results chan<- Result) {
	defer func() {
		// handle panic gracefully
		if r := recover(); r != nil {
			fmt.Println("Worker", w.id, "panicked:", r)
		}
	}()

	for i := 0; i < w.n; i++ {
		// make a GET request and measure the latency
		start := time.Now()
		resp, err := w.client.Get(w.url)
		
		latency := time.Since(start)

		// send the result to the channel
		result := Result{w.id, 0, latency, err}
		if err == nil {
			// close the response body and get the status code
			defer resp.Body.Close()
			result.status = resp.StatusCode
		}
		results <- result
	}
}

// main function
func main() {
	// constants
	const workers = 10                                 // number of workers
	const requests = 100                               // number of requests per worker
	const url = "http://localhost:9011/books?bookId=1" // url to make requests to

	// create a channel for results
	results := make(chan Result, workers*requests)

	// create a http client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// create a wait group for workers
	wg := &sync.WaitGroup{}
	wg.Add(workers)

	// create and run workers
	for i := 0; i < workers; i++ {
		worker := NewWorker(i, url, requests, client)
		go func() {
			worker.Run(results)
			wg.Done()
		}()
	}

	// wait for all workers to finish
	wg.Wait()
	close(results)

	// collect metrics
	var totalRequests, totalErrors, minLatency, maxLatency, sumLatency int64
	var avgLatency, reqPerSec, errorRate float64
	minLatency = 1<<63 - 1 // max int64 value
	startTime := time.Now()

	// iterate over the results
	for result := range results {
		totalRequests++
		sumLatency += int64(result.latency)
		if result.err != nil {
			totalErrors++
		}
		if result.latency < time.Duration(minLatency) {
			minLatency = int64(result.latency)
		}
		if result.latency > time.Duration(maxLatency) {
			maxLatency = int64(result.latency)
		}
	}

	// calculate metrics
	elapsedTime := time.Since(startTime)
	avgLatency = float64(sumLatency) / float64(totalRequests)
	reqPerSec = float64(totalRequests) / elapsedTime.Seconds()
	errorRate = float64(totalErrors) / float64(totalRequests) * 100

	// print metrics
	fmt.Println("Total Number of Requests:", totalRequests)
	fmt.Println("Average Latency:", time.Duration(avgLatency))
	fmt.Println("Requests Per Second:", reqPerSec)
	fmt.Println("Min Latency:", time.Duration(minLatency))
	fmt.Println("Max Latency:", time.Duration(maxLatency))
	fmt.Println("Error Rate:", errorRate, "%")
}
