package printer

import (
	"encoding/base64"
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/tomanikolov/packer-daemon/services"
)

// QueuePrinter ...
type QueuePrinter struct {
	qURL           string
	queue          services.QueueService
	messageGroupID string
}

// NewQueuePrinter ...
func NewQueuePrinter(qURL string, q services.QueueService) QueuePrinter {
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
	data := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s\n", message)))
	queuePrinter.queue.SendMessage(queuePrinter.qURL, data, queuePrinter.messageGroupID)
}
