package main

import (
	"context"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"time"
)

func main() {

	numberOfConcurrency := 9
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
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(2 * time.Second):
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					ch <- true
					log.Printf("Produced %t", true)
				}
				return nil
			},

			worker.WithRestartAlways()).
		Worker(
			"w2",
			func(ctx context.Context) error {
				for {
					select {
					case <-ctx.Done():
						return nil
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
