package config

type Config struct {
	HTTP struct {
		IP   string `env:"HTTP-IP"`
		Port int    `env:"HTTP-PORT"`
	}
	GRPC struct {
		IP   string `env:"HTTP-IP"`
		Port int    `env:"HTTP-PORT"`
	}

	PostgresSQL struct {
		Username string `env:"PG_USER"`
		Password string `env:"PG_PWD"`
		Host     string `env:"PG_HOST"`
		Port     string `env:"PG_PORT"`
		Database string `env:"PG_DATABASE"`
	}

	CustomerGRPC struct {
		IP   string `env:"CUSTOMER_IP"`
		Port int    `env:"CUSTOMER_PORT"`
	}

	RestaurantGRPC struct {
		IP   string `env:"RESTAURANT_IP"`
		Port int    `env:"RESTAURANT_PORT"`
	}

	Kafka []string `env:"KAFKA"`
	Topic string   `env:"TOPIC"`
}
