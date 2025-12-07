package main

import (
	"block/application"
	"block/domain"
	"block/infrastructure"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

func main() {
	workers := 50 * runtime.NumCPU()
	bufferSize := 10 * workers
	wg := sync.WaitGroup{}
	wp := infrastructure.NewWorkerPool(workers, bufferSize, &wg)
	wp.Start()

	app := application.Config{
		State:      domain.NewState(),
		WorkerPool: wp,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", "8000"),
		Handler: app.Routes(),
	}
	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until a signal is received
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	log.Println("shutting down")

	// Close the channel and wait until all jobs are finished
	wp.Shutdown()
	wp.Wait()
}
