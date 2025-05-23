package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// PostgresConfig represents the configuration for PostgreSQL database
type PostgresConfig struct {
	DSN                       string `yaml:"dsn"`
	DefaultStringSize         uint   `yaml:"defaultStringSize"`
	DisableDatetimePrecision  bool   `yaml:"disableDatetimePrecision"`
	DontSupportRenameIndex    bool   `yaml:"dontSupportRenameIndex"`
	SkipInitializeWithVersion bool   `yaml:"skipInitializeWithVersion"`
}

// RedisConfig represents the configuration for Redis
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

// JWTConfig represents the configuration for JWT authentication
type JWTConfig struct {
	Secret             string `yaml:"secret"`
	TokenLifespanHours int    `yaml:"token_lifespan_hours"`
}

// SessionConfig represents the configuration for session-based authentication
type SessionConfig struct {
	CookieName        string `yaml:"cookie_name"`
	ExpirationMinutes int    `yaml:"expiration_minutes"`
	Secure            bool   `yaml:"secure"`
	HTTPOnly          bool   `yaml:"http_only"`
	SameSite          string `yaml:"same_site"`
}

// AuthConfig represents authentication configuration settings
type AuthConfig struct {
	JWT     JWTConfig     `yaml:"jwt"`
	Session SessionConfig `yaml:"session"`
}

// HostConfig represents a configuration with hosts
type HostConfig struct {
	Hosts []string `yaml:"hosts"`
}

// LoggerConfig represents the configuration for logger
type LoggerConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

// AuthenticateAndPostConfig represents the configuration for the authenticate and post service
type AuthenticateAndPostConfig struct {
	Port               int            `yaml:"port"`
	Logger             LoggerConfig   `yaml:"logger"`
	Postgres           PostgresConfig `yaml:"postgres"`
	Redis              RedisConfig    `yaml:"redis"`
	NewsfeedPublishing HostConfig     `yaml:"newsfeed_publishing"`
	Auth               AuthConfig     `yaml:"auth"`
}

// NewsfeedConfig represents the configuration for the newsfeed service
type NewsfeedConfig struct {
	Port                int            `yaml:"port"`
	Logger              LoggerConfig   `yaml:"logger"`
	Postgres            PostgresConfig `yaml:"postgres"`
	Redis               RedisConfig    `yaml:"redis"`
	Kafka               KafkaConfig    `yaml:"kafka"`
	AuthenticateAndPost HostConfig     `yaml:"authenticate_and_post"`
}

// NewsfeedPublishingConfig represents the configuration for the newsfeed publishing service
type NewsfeedPublishingConfig struct {
	Port                int          `yaml:"port"`
	Logger              LoggerConfig `yaml:"logger"`
	Redis               RedisConfig  `yaml:"redis"`
	Kafka               KafkaConfig  `yaml:"kafka"`
	AuthenticateAndPost HostConfig   `yaml:"authenticate_and_post"`
}

// WebConfig represents the configuration for the web app
type WebConfig struct {
	Port                int          `yaml:"port"`
	Logger              LoggerConfig `yaml:"logger"`
	APIVersions         []string     `yaml:"api_version"`
	AuthenticateAndPost HostConfig   `yaml:"authenticate_and_post"`
	Newsfeed            HostConfig   `yaml:"newsfeed"`
	NewsfeedPublishing  HostConfig   `yaml:"newsfeed_publishing"`
	Redis               RedisConfig  `yaml:"redis"`
	Auth                AuthConfig   `yaml:"auth"`
}

// Config represents the main configuration for the whole system
type Config struct {
	Postgres            PostgresConfig            `yaml:"postgres"`
	Redis               RedisConfig               `yaml:"redis"`
	AuthenticateAndPost AuthenticateAndPostConfig `yaml:"authenticate_and_post_config"`
	Newsfeed            NewsfeedConfig            `yaml:"newsfeed_config"`
	NewsfeedPublishing  NewsfeedPublishingConfig  `yaml:"newsfeed_publishing_config"`
	Web                 WebConfig                 `yaml:"web_config"`
}

func parseConfig(cfgPath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return &Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return &Config{}, err
	}
	return &config, nil
}

// GetConfig loads the main configuration from a file
func GetConfig(configPath string) (*Config, error) {
	return parseConfig(configPath)
}

// GetWebConfig loads the web application configuration
func GetWebConfig(cfgPath string) (*WebConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &WebConfig{}, err
	}
	return &config.Web, nil
}

// GetNewsfeedPublishingConfig loads the newsfeed publishing configuration
func GetNewsfeedPublishingConfig(cfgPath string) (*NewsfeedPublishingConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &NewsfeedPublishingConfig{}, err
	}
	return &config.NewsfeedPublishing, nil
}

// GetNewsfeedConfig loads the newsfeed configuration
func GetNewsfeedConfig(cfgPath string) (*NewsfeedConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &NewsfeedConfig{}, err
	}
	return &config.Newsfeed, nil
}

// GetAuthenticateAndPostConfig loads the authenticate and post configuration
func GetAuthenticateAndPostConfig(cfgPath string) (*AuthenticateAndPostConfig, error) {
	config, err := parseConfig(cfgPath)
	if err != nil {
		return &AuthenticateAndPostConfig{}, err
	}
	return &config.AuthenticateAndPost, nil
}
