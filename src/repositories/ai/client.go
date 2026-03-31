package ai

// Client wraps Azure OpenAI SDK.
type Client struct {
	apiKey   string
	endpoint string
	model    string
}

func NewClient(apiKey, endpoint, model string) *Client {
	return &Client{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    model,
	}
}
