package broker

import "github.com/jeremygprawira/go-link-generator/internal/config"

func New(config *config.Configuration) *Kafka {
	return &Kafka{
		config: config,
	}
}

type Kafka struct {
	config *config.Configuration
}
