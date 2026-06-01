package config

type Config struct {
	HTTPPort   int    `envconfig:"HTTP_PORT" default:"8080"`
	DBHost     string `envconfig:"DB_HOST"`
	DBPort     int    `envconfig:"DB_PORT" default:"5432"`
	DBUser     string `envconfig:"DB_USER"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

func New() *Config {
	return &Config{}
}
