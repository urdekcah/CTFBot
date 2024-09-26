package env

import "os"

type Config struct {
	Values map[string]string
}

func Load(filename string) (*Config, error) {
	p := NewParser()
	v := NewValidator()

	rawConfig, err := p.Parse(filename)
	if err != nil {
		return nil, err
	}

	if err := v.Validate(rawConfig); err != nil {
		return nil, err
	}

	cfg := &Config{Values: rawConfig}
	if err := cfg.setEnvVariables(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) setEnvVariables() error {
	for key, value := range c.Values {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}
