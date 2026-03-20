package config

type (
	Configuration struct {
		Application   Application   `mapstructure:"application"`
		Authorization Authorization `mapstructure:"authorization"`
		CORS          CORS          `mapstructure:"cors"`
		PostgreSQL    PostgreSQL    `mapstructure:"postgresql"`
		Tracing       Tracing       `mapstructure:"tracing"`
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
		Issuer  string             `mapstructure:"issuer"`
		APIKey  string             `mapstructure:"api_key"`
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
)
