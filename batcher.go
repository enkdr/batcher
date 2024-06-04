// allow the caller to submit a single Job, and it should return a JobResult
// process accepted Jobs in batches using a BatchProcessor
// Don't implement BatchProcessor. This should be a dependency of your library.
// provide a way to configure the batching behaviour i.e. size and frequency
// expose a shutdown method which returns after all previously accepted Jobs are processed

package main

import (
	"fmt"
	"sync"
)

// create custom types to hold Job and JobResults
type Job struct {
	ID    int
	Value string
}

type JobResult struct {
	JobID  int
	Status string
}

// create worker to pass in id and channels with waitgroup for syncing
// this is where we would use the batchProcessor
func worker(id int, jobs <-chan Job, results chan<- JobResult, wg *sync.WaitGroup) {
	// notify waitgroup when all jobs done
	defer wg.Done()

	// loop through jobs channel and add to results channel
	for job := range jobs {
		results <- JobResult{JobID: job.ID, Status: "processed"}
	}
}

// collect results
func collectJobResults(results <-chan JobResult, wg *sync.WaitGroup) {
	// notify waitgroup when all results collected
	defer wg.Done()
	for result := range results {
		fmt.Printf("Job ID %d, Status: %s\n", result.JobID, result.Status)
	}
}

func batcher(jobCount, workerCount int) {
	// create channels for jobs and results
	jobs := make(chan Job, jobCount)
	results := make(chan JobResult, jobCount)

	var wg sync.WaitGroup

	wg.Add(workerCount)
	// begin processing jobs and update JobResults column
	for wcount := 0; wcount < workerCount; wcount++ {
		go worker(wcount, jobs, results, &wg)
	}

	// begin to collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	var resultWaitGroup sync.WaitGroup
	resultWaitGroup.Add(1)
	go collectJobResults(results, &resultWaitGroup)

	for j := 0; j < jobCount; j++ {
		jobs <- Job{ID: j, Value: fmt.Sprintf("Job number %d", j)}
	}
	// close
	close(jobs)

	// ensure we have all results
	resultWaitGroup.Wait()
}

func main() {
	// specify how many jobs, workers
	const jobsTotal = 100
	const workersTotal = 3

	batcher(jobsTotal, workersTotal)
}
