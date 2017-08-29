package printer

import (
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/tomanikolov/packer-daemon/queue"
)

// QueuePrinter ...
type QueuePrinter struct {
	qURL           string
	queue          queue.Queue
	messageGroupID string
}

// NewQueuePrinter ...
func NewQueuePrinter(qURL string, q queue.Queue) QueuePrinter {
	u, _ := uuid.NewV4()
	messageGroupID := u.String()
	return QueuePrinter{
		qURL:           qURL,
		queue:          q,
		messageGroupID: messageGroupID,
	}
}

// Print ...
func (queuePrinter QueuePrinter) Print(message string) {
	queuePrinter.queue.SendMessage(queuePrinter.qURL, fmt.Sprintf("%s\n", message), queuePrinter.messageGroupID)
}
