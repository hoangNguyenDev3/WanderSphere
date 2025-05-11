package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

// ParseServiceConfig loads a specific service configuration directly from a YAML file
func ParseServiceConfig(cfgPath string, cfg interface{}) error {
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}

	return nil
}

// GetAuthenticateAndPostConfigDirect loads the authenticate and post configuration directly from its own file
func GetAuthenticateAndPostConfigDirect(cfgPath string) (*AuthenticateAndPostConfig, error) {
	cfg := &AuthenticateAndPostConfig{}
	err := ParseServiceConfig(cfgPath, cfg)
	return cfg, err
}

// GetNewsfeedConfigDirect loads the newsfeed configuration directly from its own file
func GetNewsfeedConfigDirect(cfgPath string) (*NewsfeedConfig, error) {
	cfg := &NewsfeedConfig{}
	err := ParseServiceConfig(cfgPath, cfg)
	return cfg, err
}

// GetNewsfeedPublishingConfigDirect loads the newsfeed publishing configuration directly from its own file
func GetNewsfeedPublishingConfigDirect(cfgPath string) (*NewsfeedPublishingConfig, error) {
	cfg := &NewsfeedPublishingConfig{}
	err := ParseServiceConfig(cfgPath, cfg)
	return cfg, err
}

// GetWebConfigDirect loads the web application configuration directly from its own file
func GetWebConfigDirect(cfgPath string) (*WebConfig, error) {
	cfg := &WebConfig{}
	err := ParseServiceConfig(cfgPath, cfg)
	return cfg, err
}
