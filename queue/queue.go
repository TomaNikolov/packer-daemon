package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/tomanikolov/packer-daemon/types"
)

// Queue ...
type Queue struct {
	sqs *sqs.SQS
}

// NewQueue ...
func NewQueue(c types.Config) Queue {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sqs := sqs.New(sess, &aws.Config{
		Region:      aws.String(c.AwsRegion),
		Credentials: credentials.NewStaticCredentials(c.AwsPublicKey, c.AwsPriveteKey, ""),
	})

	return Queue{
		sqs: sqs,
	}
}

// SendMessage ...
func (q Queue) SendMessage(qURL, message, messageGroupID string) (*types.SendMessageOutput, error) {
	//u, _ := uuid.NewV4()
	//messageGroupId := u.String()
	u, _ := uuid.NewV4()
	messageDeduplicationID := u.String()
	s, err := q.sqs.SendMessage(&sqs.SendMessageInput{
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
func (q Queue) ReceiveMessage(qURL string) (*types.ReciveMessageOutput, error) {
	result, err := q.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &qURL,
	})

	if err != nil {
		return nil, err
	}

	messages := q.mapMessages(result.Messages)
	reciveMessageOutput := &types.ReciveMessageOutput{
		Messages: messages,
	}

	return reciveMessageOutput, nil
}

// DeleteMessage ...
func (q Queue) DeleteMessage(qURL string, receiptHandle string) error {
	_, err := q.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &qURL,
		ReceiptHandle: &receiptHandle,
	})

	return err
}

func (q Queue) mapMessages(sqsMessages []*sqs.Message) []types.Message {
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
