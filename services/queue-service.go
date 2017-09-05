package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/tomanikolov/packer-daemon/types"
)

// QueueService ...
type QueueService struct {
	sqs *sqs.SQS
}

// NewQueueService ...
func NewQueueService(c types.Config) QueueService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqs := sqs.New(sess, &aws.Config{
		Region:      aws.String(c.AwsRegion),
		Credentials: credentials.NewStaticCredentials(c.AwsPublicKey, c.AwsPriveteKey, ""),
	})

	return QueueService{
		sqs: sqs,
	}
}

// SendMessage ...
func (service QueueService) SendMessage(qURL, message, messageGroupID string) (*types.SendMessageOutput, error) {
	//u, _ := uuid.NewV4()
	//messageGroupId := u.String()
	u, _ := uuid.NewV4()
	messageDeduplicationID := u.String()
	s, err := service.sqs.SendMessage(&sqs.SendMessageInput{
		QueueUrl:               &qURL,
		MessageBody:            &message,
		MessageGroupId:         &messageGroupID,
		MessageDeduplicationId: &messageDeduplicationID,
	})

	if err != nil {
		return nil, err
	}

	sendMessageOutput := &types.SendMessageOutput{
		MessageID:      s.MessageId,
		SequenceNumber: s.SequenceNumber,
	}

	return sendMessageOutput, nil
}

// ReceiveMessage ...
func (service QueueService) ReceiveMessage(qURL string) (*types.ReciveMessageOutput, error) {
	result, err := service.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &qURL,
	})

	if err != nil {
		return nil, err
	}

	messages := service.mapMessages(result.Messages)
	reciveMessageOutput := &types.ReciveMessageOutput{
		Messages: messages,
	}

	return reciveMessageOutput, nil
}

// DeleteMessage ...
func (service QueueService) DeleteMessage(qURL string, receiptHandle string) error {
	_, err := service.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: &receiptHandle,
	})

	return err
}

func (service QueueService) mapMessages(sqsMessages []*sqs.Message) []types.Message {
	messages := []types.Message{}
	for _, m := range sqsMessages {
		message := types.Message{
			Attributes:    m.Attributes,
			Body:          m.Body,
			MessageID:     m.MessageId,
			ReceiptHandle: m.ReceiptHandle,
		}

		messages = append(messages, message)
	}

	return messages
}
