package config

import "fmt"

func validate() error {
	// Check for required fields
	if Config.Name == "" {
		return fmt.Errorf("Configuration: name is required")
	}

	if Config.API.Host == "" {
		return fmt.Errorf("Configuration: api.host is required")
	}

	if Config.API.Port == 0 {
		return fmt.Errorf("Configuration: api.port is required")
	}

	if Config.Database.URL == "" {
		return fmt.Errorf("Configuration: database.url is required")
	}
	return nil
}
