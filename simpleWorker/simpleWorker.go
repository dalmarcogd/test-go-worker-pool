package main

import (
	"context"
	"errors"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/pkg/worker"
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
			func(ctx context.Context) error {
				select {
				case <-time.After(10 * time.Second):
					return errors.New("test")
				case <-ctx.Done():
					return errors.New("timeout")

				}
			},
			worker.WithRestartAlways(),
			worker.WithTimeout(11*time.Second)).
		Worker(
			"w2",
			func(ctx context.Context) error {
				select {
				case <-time.After(30 * time.Second):
				case <-ctx.Done():
					return ctx.Err()

				}
				return nil
			}).
		Worker(
			"w3",
			func(ctx context.Context) error {
				select {
				case <-time.After(1 * time.Minute):
				case <-ctx.Done():
					return errors.New("test")

				}
				return nil
			}).
		Run(); err != nil {
		panic(err)
	}
}
