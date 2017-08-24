package main

import (
	"flag"
	"log"

	"github.com/nabeken/aws-go-sqs/queue"
	"github.com/packer-daemon/worker"
	"github.com/stripe/aws-go/gen/sqs"
)

func main() {
	flag.Parse()

	q, err := NewSQSQueue(
		sqs.New(aws.DetectCreds("", "", ""), "ap-northeast-1", nil),
		*flagQueue,
	)
	if err != nil {
		log.Fatal(err)
	}

	worker.Start(q, worker.HandlerFunc(Print))
}

func NewSQSQueue(s *sqs.SQS, name string) (*queue.Queue, error) {
	return queue.New(s, name)
}
