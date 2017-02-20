package config

type OutputConfig struct {
	Service    string
	WebhookURL string
	FilterPath string
}

type Config struct {
	OutputWebhookURL           string
	NamesJSONURL               string
	MovesJSONURL               string
	NormalThumbnailURLTemplate string
	ShinyThumbnailURLTemplate  string
	Port                       int
	Outputs                    []OutputConfig
}
