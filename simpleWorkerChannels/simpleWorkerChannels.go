package main

import (
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/worker"
	"log"
	"time"
)

func main() {

	ch := make(chan bool)

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
				log.Printf("Produced %t", true)
				return nil
			},
			1,
			true).
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
			1,
			false).
		Run(); err != nil {
		panic(err)
	}
}
