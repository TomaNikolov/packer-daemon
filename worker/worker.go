package worker

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/tomanikolov/packer-daemon/logger"
	"github.com/tomanikolov/packer-daemon/printer"
	"github.com/tomanikolov/packer-daemon/services"
	"github.com/tomanikolov/packer-daemon/types"
)

// Start ...
func Start(config types.Config) {
	log.Println("worker: Start polling")
	queue := services.NewQueueService(config)
	for {
		result, err := queue.ReceiveMessage(config.QueueURL)
		if err != nil {
			log.Println(err)
			continue
		}

		if len(result.Messages) > 0 {
			run(queue, result.Messages, config)
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func run(q services.QueueService, messages []types.Message, c types.Config) {
	numMessages := len(messages)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for _, message := range messages {
		go func(m types.Message) {
			defer wg.Done()
			err := handleMessage(q, m, c)
			if err != nil {
				log.Println(err)
			}
		}(message)
	}

	wg.Wait()
}

func handleMessage(q services.QueueService, m types.Message, c types.Config) error {
	buildRequest, err := getBuildRequest(*m.Body)
	logger := getLogger(buildRequest.LogQURL, q)
	buildService := services.NewTemplateBuildService(buildRequest, c, logger)
	err = buildService.Start()
	if err != nil {
		logger.LogError(err.Error())
	}

	return q.DeleteMessage(c.QueueURL, *m.ReceiptHandle)
}

func getBuildRequest(buildMessage string) (types.BuildRequest, error) {
	buildRequest := types.BuildRequest{}
	err := json.Unmarshal([]byte(buildMessage), &buildRequest)
	return buildRequest, err
}

func getLogger(qLogURL string, q services.QueueService) logger.Logger {
	consolePrinter := printer.NewConsolePrinter()
	queuePrinter := printer.NewQueuePrinter(qLogURL, q)
	return logger.NewLogger([]types.Printer{consolePrinter, queuePrinter})
}
