package main

import (
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"time"
)

func main() {

	numberOfConcurrency := 10
	ch := make(chan bool, numberOfConcurrency)

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
				ch <- true
				ch <- true
				ch <- true
				ch <- true
				ch <- true
				ch <- true
				ch <- true
				log.Printf("Produced %t", true)
				return nil
			},

			worker.WithRestartAlways()).
		Worker(
			"w2",
			func() error {
				for {
					select {
					case r := <-ch:
						log.Printf("Received %t", r)
					}
				}
			},
			worker.WithConcurrency(numberOfConcurrency)).
		Run(); err != nil {
		panic(err)
	}
}
