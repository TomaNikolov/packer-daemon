package printer

// QueuePrinter ...
type QueuePrinter struct {
}

// NewQueuePrinter ...
func NewQueuePrinter() QueuePrinter {
	return QueuePrinter{}
}

// Print ...
func (queuePrinter QueuePrinter) Print(message string) {
	// TODO: implement log queue and push stdout to it.
}
