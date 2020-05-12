package main

import (
	"errors"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"time"
)

func main() {
	if err := gwp.
		New().
		Stats().
		HealthCheck().
		DebugPprof().
		HandleError(func(w *worker.Worker, err error) {
			log.Printf("Worker [%s] error: %s", w.Name, err)
		}).
		Worker(
			"w1",
			func() error {
				<-time.After(10 * time.Second)
				return errors.New("test")
			},
			worker.WithRestartAlways()).
		Worker(
			"w2",
			func() error {
				<-time.After(30 * time.Second)
				return nil
			}).
		Worker(
			"w3",
			func() error {
				<-time.After(1 * time.Minute)
				return errors.New("test")
			}).
		Run(); err != nil {
		panic(err)
	}
}
