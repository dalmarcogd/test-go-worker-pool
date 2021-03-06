package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/dalmarcogd/gwp"
	"github.com/dalmarcogd/gwp/pkg/worker"
	"log"
	"strconv"
)

func main() {

	params := &sqs.CreateQueueInput{
		QueueName: aws.String("test-consume-sqs"), // Required
	}
	ss, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String("http://localhost:9324"),
		Region:   aws.String("us-east-1"),
	})
	svc := sqs.New(ss)

	var resp, err = svc.CreateQueue(params)

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println(resp)

	queueURL := aws.String("http://localhost:9324/queue/test-consume-sqs")

	for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		paramsSend := &sqs.SendMessageInput{
			MessageBody: aws.String("Testing " + strconv.Itoa(i)), // Required
			QueueUrl:    queueURL,                                 // Required
		}
		respSend, err := svc.SendMessage(paramsSend)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		fmt.Println(respSend)
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
			params := &sqs.ReceiveMessageInput{
				QueueUrl:            queueURL, // Required
				MaxNumberOfMessages: aws.Int64(10),
				VisibilityTimeout:   aws.Int64(20),
			}

			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					resp, err := svc.ReceiveMessage(params)
					if err != nil {
						fmt.Println(err.Error())
						return err
					}
					fmt.Println(resp.Messages)
					for _, msg := range resp.Messages {
						fmt.Println(aws.StringValue(msg.Body))
					}
				}
			}
		}, worker.WithRestartAlways()).
		Run(); err != nil {
		panic(err)
	}
}
