package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config represents the application configuration
type Config struct {
	DB         DBConfig `env:",prefix=DB_,required"`
	JWTSecret  string   `env:"JWT_SECRET"`
	BcryptSalt int      `env:"BCRYPT_SALT"`
	S3Region   string   `env:"S3_REGION"`
	S3Bucket   string   `env:"S3_BUCKET_NAME"`
	S3Id       string   `env:"S3_ID"`
	S3Secret   string   `env:"S3_SECRET_KEY"`
}

type DBConfig struct {
	Name     string `env:"NAME"`
	Port     string `env:"PORT"`
	Host     string `env:"HOST"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Params   string `env:"PARAMS"`
}

func Load(ctx context.Context) (*Config, error) {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ConnectionURL returns the connection URL for the database
func (c DBConfig) ConnectionString() string {
	params := strings.ReplaceAll(c.Params, `"`, "")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Name, params)
}
