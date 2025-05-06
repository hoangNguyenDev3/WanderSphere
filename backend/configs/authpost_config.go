package configs

func GetAuthenticateAndPostConfig(cfgPath string) (*AuthenticateAndPostConfig, error) {
	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	return &cfg.AuthenticateAndPost, nil
}
