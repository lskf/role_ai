package llm

type Options struct {
	BaseURL string
	ApiKey  string
}

type Option func(o *Options)

func WithBaseURL(baseURL string) Option {
	return func(o *Options) {
		o.BaseURL = baseURL
	}
}

func WithApiKey(apiKey string) Option {
	return func(o *Options) {
		o.ApiKey = apiKey
	}
}
