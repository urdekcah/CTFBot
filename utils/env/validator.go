package env

import (
	"fmt"
	"regexp"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(config map[string]string) error {
	for key, value := range config {
		if err := v.validateKey(key); err != nil {
			return err
		}
		if err := v.validateValue(value); err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty key is not allowed")
	}
	if matched, _ := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_]*$`, key); !matched {
		return fmt.Errorf("invalid key format: %s", key)
	}
	return nil
}

func (v *Validator) validateValue(value string) error {
	return nil
}
