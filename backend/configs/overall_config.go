package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseConfig(cfgPath string) (*Config, error) {
	// Read config file
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	// Unmarshal config
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Config mirrors the config.yaml file
type Config struct {
	Postgres            PostgresConfig            `yaml:"postgres"`
	Redis               RedisConfig               `yaml:"redis"`
	AuthenticateAndPost AuthenticateAndPostConfig `yaml:"authenticate_and_post_config"`
	Newsfeed            NewsfeedConfig            `yaml:"newsfeed_config"`
	WebConfig           WebConfig                 `yaml:"web_config"`
}

// PostgresConfig holds the shared PostgreSQL settings.
type PostgresConfig struct {
	DSN         string `yaml:"dsn"`
	MaxPoolSize int    `yaml:"max_pool_size"`
	MinPoolSize int    `yaml:"min_pool_size"`
	SearchPath  string `yaml:"search_path"`
}

// RedisConfig holds Redis connection info.
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

// AuthenticateAndPostConfig config for the auth+post service.
type AuthenticateAndPostConfig struct {
	Port     int            `yaml:"port"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

// NewsfeedConfig config for the newsfeed service.
type NewsfeedConfig struct {
	Port     int            `yaml:"port"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

// WebConfig config for the BFF/web-app service.
type WebConfig struct {
	Port                int         `yaml:"port"`
	APIVersions         []string    `yaml:"api_version"`
	AuthenticateAndPost HostConfig  `yaml:"authenticate_and_post"`
	Newsfeed            HostConfig  `yaml:"newsfeed"`
	Redis               RedisConfig `yaml:"redis"`
}

// HostConfig holds a list of host:port strings.
type HostConfig struct {
	Hosts []string `yaml:"hosts"`
}
