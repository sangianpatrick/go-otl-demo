package config

import (
	"os"
	"strconv"
)

type Config struct {
	Application struct {
		Name string
		Port int
	}
	Jaeger struct {
		Host string
	}
}

func (c *Config) application() *Config {
	var port int = 9000 // default port
	name := os.Getenv("APPLICATION_NAME")
	desiredPort, _ := strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	if desiredPort > 0 {
		port = desiredPort
	}

	c.Application.Name = name
	c.Application.Port = port

	return c
}

func (c *Config) jaeger() *Config {
	host := os.Getenv("JAEGER_HOST")

	c.Jaeger.Host = host

	return c
}

func load() *Config {
	c := new(Config)
	return c.application().jaeger()
}

func Get() *Config {
	return load()
}
