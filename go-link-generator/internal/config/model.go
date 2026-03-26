package config

import "time"

type (
	Configuration struct {
		Application   Application   `mapstructure:"application"`
		Authorization Authorization `mapstructure:"authorization"`
		CORS          CORS          `mapstructure:"cors"`
		PostgreSQL    PostgreSQL    `mapstructure:"postgresql"`
		Tracing       Tracing       `mapstructure:"tracing"`
		Url           Url           `mapstructure:"url"`
		Kafka         Kafka         `mapstructure:"kafka"`
	}

	Application struct {
		Name        string `mapstructure:"name"`
		Version     string `mapstructure:"version"`
		Port        int    `mapstructure:"port"`
		Environment string `mapstructure:"environment"`
		Host        string `mapstructure:"host"`
		Timeout     int    `mapstructure:"timeout"`
		Timezone    string `mapstructure:"timezone"`
	}

	CORS struct {
		HeadersAllowed []string `mapstructure:"headers_allowed"`
	}

	Authorization struct {
		Issuer string `mapstructure:"issuer"`
		APIKey string `mapstructure:"api_key"`
	}

	PostgreSQL struct {
		Name            string `mapstructure:"name"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		SSLMode         string `mapstructure:"ssl_mode"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
	}

	Tracing struct {
		Endpoint    string `mapstructure:"endpoint"`
		ServiceName string `mapstructure:"service_name"`
	}

	Kafka struct {
		Brokers  []string      `mapstructure:"brokers"`
		Topics   KafkaTopics   `mapstructure:"topics"`
		Producer KafkaProducer `mapstructure:"producer"`
		Consumer KafkaConsumer `mapstructure:"consumer"`
	}

	KafkaTopics struct {
		Link    string `mapstructure:"link"`
		LinkDLQ string `mapstructure:"link_dlq"`
	}

	// ProducerConfig holds producer-specific configuration.
	KafkaProducer struct {
		RequiredAcks    string               `mapstructure:"required_acks"`
		Retry           KafkaRetry           `mapstructure:"retry"`
		Partitioner     string               `mapstructure:"partitioner"`
		MaxMessageBytes int                  `mapstructure:"max_message_bytes"`
		Timeout         KafkaProducerTimeout `mapstructure:"timeout"`
	}

	KafkaRetry struct {
		Max     int           `mapstructure:"max"`
		Backoff time.Duration `mapstructure:"backoff"`
	}

	KafkaProducerTimeout struct {
		Dial  time.Duration `mapstructure:"dial"`
		Read  time.Duration `mapstructure:"read"`
		Write time.Duration `mapstructure:"write"`
	}

	KafkaConsumer struct {
		GroupID           string               `mapstructure:"group_id"`
		InitialOffset     string               `mapstructure:"initial_offset"`
		Retry             KafkaRetry           `mapstructure:"retry"`
		Timeout           KafkaConsumerTimeout `mapstructure:"timeout"`
		HeartbeatInterval time.Duration        `mapstructure:"heartbeat_interval"`
	}

	KafkaConsumerTimeout struct {
		Dial      time.Duration `mapstructure:"dial"`
		Read      time.Duration `mapstructure:"read"`
		Write     time.Duration `mapstructure:"write"`
		Session   time.Duration `mapstructure:"session"`
		Rebalance time.Duration `mapstructure:"rebalance"`
	}

	Url struct {
		Length                int    `mapstructure:"length"`
		CodeGenerationRetries int    `mapstructure:"code_generation_retries"`
		CodeGenerationBackoff int    `mapstructure:"code_generation_backoff"`
		Secret                string `mapstructure:"secret"`
		SnowflakeMachineID    int64  `mapstructure:"snowflake_machine_id"`
		SecureLength          int    `mapstructure:"secure_length"`
	}
)
