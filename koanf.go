package main

type Config struct {
	Name string `koanf:"name"`
	API  struct {
		Host string `koanf:"host"`
		Port int    `koanf:"port"`
	}
	Database struct {
		URL string `koanf:"url"`
	}
}
