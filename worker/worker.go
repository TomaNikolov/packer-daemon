package worker

import (
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tomanikolov/packer-daemon/builder"
	"github.com/tomanikolov/packer-daemon/types"
)

// Start ...
func Start(q *sqs.SQS, config types.Config) {
	// URL to our queue
	log.Println("worker: Start polling")
	for {
		result, err := q.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: &config.QueueURL,
		})

		if err != nil {
			log.Println(err)
			continue
		}

		if len(result.Messages) > 0 {
			run(q, result.Messages, config)
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func run(q *sqs.SQS, messages []*sqs.Message, c types.Config) {
	numMessages := len(messages)
	log.Printf("worker: Received %d messages", numMessages)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			defer wg.Done()
			_, err := handleMessage(q, m, c)
			if err != nil {
				log.Println(err)
			}
		}(messages[i])
	}

	wg.Wait()
}

func handleMessage(q *sqs.SQS, m *sqs.Message, c types.Config) (*sqs.DeleteMessageOutput, error) {
	err := builder.Start(*m.Body, c)
	if err != nil {
		return nil, err
	}

	return q.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &c.QueueURL,
		ReceiptHandle: m.ReceiptHandle,
	})
}
