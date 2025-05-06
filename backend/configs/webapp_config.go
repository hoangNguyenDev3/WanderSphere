package configs

func GetWebConfig(cfgPath string) (*WebConfig, error) {
	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	return &cfg.WebConfig, nil
}
