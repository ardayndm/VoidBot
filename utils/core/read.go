package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Read - YAML Dosyasını okur
func ReadYaml(path string, to any) error {

	if data, err := os.ReadFile(path); err == nil {
		err = yaml.Unmarshal(data, to)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}
