package setups

import (
	"github.com/caarlos0/env"
)

// Application contains data related to application configuration parameters.
type Application struct {
	DryRun          bool   `env:"KBS_DRY_RUN" envDefault:"false"`
	ApplicationPort string `env:"KBS_APPLICATION_PORT" envDefault:":8080"`
	LogLevel        string `env:"KBS_LOG_ENVIRONMENT" envDefault:"production"`
	Repository      RepositoryParameters
}

// RepositoryParameters contains data related to a repository.
type RepositoryParameters struct {
	Region   string `env:"KBS_AWS_REGION" envDefault:"us-east-1"`
	Endpoint string `env:"KBS_AWS_ENDPOINT" envDefault:"5432"`
}

const (
	ProductionLog  = "production"
	DevelopmentLog = "development"
)

var (
	Version    string
	BuildDate  string
	CommitHash string
)

// Load load application configuration
func Load() (Application, error) {
	cfg := Application{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	repository := RepositoryParameters{}
	if err := env.Parse(&repository); err != nil {
		return cfg, err
	}
	cfg.Repository = repository
	return cfg, nil
}
