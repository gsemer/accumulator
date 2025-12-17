package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"signal/application"
	"signal/persistence"
	"syscall"

	"github.com/robfig/cron/v3"
)

func main() {
	accumulatorAPI := "http://localhost:8000"

	repository := persistence.NewBackgroundRepository(http.Client{}, accumulatorAPI)
	service := application.NewBackgroundService(repository)

	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 * * * * *", func() {
		log.Println("Crontab job triggered.")

		_, err := service.Run()
		if err != nil {
			log.Printf("Run failed: %v", err)
		}
	})

	c.Start()
	log.Println("Crontab job started.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
