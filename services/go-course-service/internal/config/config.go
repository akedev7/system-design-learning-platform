package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	S3       S3Config
}

type ServerConfig struct {
	Port  int
	Debug bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type AuthConfig struct {
	JWKSURL  string
	Audience string
	Issuer   string
}

type S3Config struct {
	Endpoint        string
	Region         string
	Bucket         string
	AccessKey      string
	SecretKey      string
	UseSSL         bool
	ForcePathStyle bool
	PublicURL      string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.debug", true)
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "courses")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("s3.region", "auto")
	v.SetDefault("s3.bucket", "uploads")
	v.SetDefault("s3.useSSL", true)
	v.SetDefault("s3.forcePathStyle", true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyEnvOverrides(&cfg)

	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = atoi(port)
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		cfg.Database.Port = atoi(port)
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		cfg.Database.DBName = dbname
	}
	if sslmode := os.Getenv("DB_SSLMODE"); sslmode != "" {
		cfg.Database.SSLMode = sslmode
	}
	if jwksURL := os.Getenv("AUTH_JWKS_URL"); jwksURL != "" {
		cfg.Auth.JWKSURL = jwksURL
	}
	if audience := os.Getenv("AUTH_AUDIENCE"); audience != "" {
		cfg.Auth.Audience = audience
	}
	if issuer := os.Getenv("AUTH_ISSUER"); issuer != "" {
		cfg.Auth.Issuer = issuer
	}
	if endpoint := os.Getenv("S3_ENDPOINT"); endpoint != "" {
		cfg.S3.Endpoint = endpoint
	}
	if region := os.Getenv("S3_REGION"); region != "" {
		cfg.S3.Region = region
	}
	if bucket := os.Getenv("S3_BUCKET"); bucket != "" {
		cfg.S3.Bucket = bucket
	}
	if accessKey := os.Getenv("S3_ACCESS_KEY"); accessKey != "" {
		cfg.S3.AccessKey = accessKey
	}
	if secretKey := os.Getenv("S3_SECRET_KEY"); secretKey != "" {
		cfg.S3.SecretKey = secretKey
	}
	if publicURL := os.Getenv("S3_PUBLIC_URL"); publicURL != "" {
		cfg.S3.PublicURL = publicURL
	}
}

func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
