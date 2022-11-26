package main

import "fmt"

func checkConfig() (err error) {
	// Check for required configuration
	if k.String(name) == "" {
		err = fmt.Errorf("Configuration: name is required")
	}
	if k.String(api.host) == "" {
		err = fmt.Errorf("Configuration: api.host is required")
	}
	if k.Int(api.port) == 0 {
		err = fmt.Errorf("Configuration: api.port is required")
	}
	if k.String(database.url) == "" {
		err = fmt.Errorf("Configuration: database.url is required")
	}
	return
}
