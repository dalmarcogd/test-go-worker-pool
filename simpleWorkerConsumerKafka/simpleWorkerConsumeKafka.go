package main

import (
	"context"
	"fmt"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	topic := "teste"
	partition := 1

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	failOnError(err, "Fail when create connection")

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, _ = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)

	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	if err := gwp.
		New().
		Stats().
		HealthCheck().
		DebugPprof().
		HandleError(func(w *worker.Worker, err error) {
			log.Printf("Worker [%s] error: %s", w.Name, err)
		}).
		Worker("w2", func(ctx context.Context) error {
			batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
			b := make([]byte, 10e3)            // 10KB max per message
			for {
				select {
				case <-ctx.Done():
					_ = batch.Close()
					return nil
				default:
					_, err := batch.Read(b)
					if err != nil {
						break
					}
					fmt.Println(string(b))
				}
			}
		},  worker.WithRestartAlways()).
		Run(); err != nil {
		panic(err)
	}
}
