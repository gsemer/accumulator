package infrastructure

import (
	"accumulator/domain"
	"log"
	"sync"
)

type WorkerPool struct {
	workers int
	jobs    chan domain.Job
	wg      *sync.WaitGroup
}

func NewWorkerPool(workers, size int, wg *sync.WaitGroup) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan domain.Job, size),
		wg:      wg,
	}
}

// Worker processes jobs from the channel
func (wp *WorkerPool) worker(i int) {
	for job := range wp.jobs {
		log.Printf("Worker %d processes job", i)
		err := job.Process()
		if err != nil {
			log.Println("Error while processing the job:", err)
		}
		wp.wg.Done()
	}
}

// Start launches worker goroutines
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go wp.worker(i)
	}
}

// Adds a job to the queue
func (wp *WorkerPool) Submit(job domain.Job) {
	wp.wg.Add(1)
	wp.jobs <- job
}

// Closes the job channel
func (wp *WorkerPool) Shutdown() {
	close(wp.jobs)
}

// Waits until all submitted jobs to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
