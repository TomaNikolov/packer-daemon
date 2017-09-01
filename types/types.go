package types

// Config ...
type Config struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	GitUsername   string `json:"gitUsername"`
	GitPassword   string `json:"gitPassword"`
	AwsPublicKey  string `json:"awsPublicKey"`
	AwsPriveteKey string `json:"awsPriveteKey"`
	AwsRegion     string `json:"awsRegion"`
	StoragePath   string `json:"storagePath"`
	Repository    string `json:"repository"`
	QueueURL      string `json:"queueUrl"`
	GovcPassword  string `json:"govcPassword"`
	GovcUsername  string `json:"govcUsername"`
	GovcURL       string `json:"govcUrl"`
	GovcInsecure  string `json:"govcInsecure"`
}

// BuildRequest ...
type BuildRequest struct {
	Branch        string `json:"branch"`
	TemplateName  string `json:"templateName"`
	PackerOptions string `json:"packerOptions"`
	Stage         string `json:"stage"`
	LogQURL       string `json:"logQURL"`
}

// Printer ...
type Printer interface {
	Print(message string)
}

// SendMessageOutput ...
type SendMessageOutput struct {
	MessageID      *string
	SequenceNumber *string
}

// ReciveMessageOutput ...
type ReciveMessageOutput struct {
	Messages []Message
}

// Message ...
type Message struct {
	Attributes    map[string]*string
	Body          *string
	MessageID     *string
	ReceiptHandle *string
}
