package main

import (
	"context"
	"fmt"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	connection, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672//")

	failOnError(err, "Error when get connection")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Error when get channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"test-consume", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Error when declare a queue")

	for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		failOnError(channel.Publish("", queue.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Body:         []byte(fmt.Sprint(i)),
		}), "fail on publishing")
	}

	if err := gwp.
		New().
		Stats().
		HealthCheck().
		DebugPprof().
		HandleError(func(w *worker.Worker, err error) {
			log.Printf("Worker [%s] error: %s", w.Name, err)
		}).
		Worker("w2", func(ctx context.Context) error {
			msgs, err := channel.Consume(queue.Name,
				"",
				true,
				false,
				false,
				false,
				nil)
			failOnError(err, "Error when create consumer")

			for  {
				select {
				case <-ctx.Done():
					return nil
				case msg := <-msgs :
					fmt.Println(string(msg.Body))
				}
			}
		}, worker.WithRestartAlways()).
		Run(); err != nil {
		panic(err)
	}
}
