package main

import (
	"sync"
	"testing"
)

func TestBatcher(t *testing.T) {
	const jobsTotal = 10
	const workersTotal = 3

	jobs := make(chan Job, jobsTotal)
	results := make(chan JobResult, jobsTotal)

	var wg sync.WaitGroup

	wg.Add(workersTotal)
	// start workers
	for w := 0; w < workersTotal; w++ {
		go worker(w, jobs, results, &wg)
	}

	// collect our results in separate goroutine
	var resultWaitGroup sync.WaitGroup
	resultWaitGroup.Add(1)
	collectedResults := []JobResult{}
	go func() {
		defer resultWaitGroup.Done()
		for result := range results {
			collectedResults = append(collectedResults, result)
		}
	}()

	// send jobs to the jobs channel
	for j := 0; j < jobsTotal; j++ {
		jobs <- Job{ID: j, Value: "test"}
	}
	close(jobs)

	// wait for all workers to finish and close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// wait for all results to be collected
	resultWaitGroup.Wait()

	// check the results
	if len(collectedResults) != jobsTotal {
		t.Fatalf("Expected %d results, but got %d", jobsTotal, len(collectedResults))
	}

	expectedStatus := "processed"
	for _, result := range collectedResults {
		if result.Status != expectedStatus {
			t.Errorf("Expected Job ID %d to have status '%s', but got '%s'", result.JobID, expectedStatus, result.Status)
		}
	}
}
