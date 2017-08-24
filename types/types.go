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
	QueueURL      string `json:"queueURL"`
}

// BuildRequest ...
type BuildRequest struct {
	Branch        string `json:"branch"`
	TemplateName  string `json:"templateName"`
	PackerOptions string `json:"packerOptions"`
	Stage         string `json:"stage"`
}
