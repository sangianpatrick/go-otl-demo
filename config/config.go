package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Application struct {
		Name       string
		Port       int
		Enviroment string
	}
	Opentelemetry struct {
		CollectorHost string
	}
	JSONPlaceHolderAPI struct {
		Host string
	}
	Mariadb struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
}

func (c *Config) application() *Config {
	var port int = 9000 // default port
	name := os.Getenv("APPLICATION_NAME")
	desiredPort, _ := strconv.Atoi(os.Getenv("APPLICATION_PORT"))
	environment := os.Getenv("APPLICATION_ENVIRONMENT")
	if desiredPort > 0 {
		port = desiredPort
	}

	c.Application.Name = name
	c.Application.Port = port
	c.Application.Enviroment = environment

	return c
}

func (c *Config) opentelemetry() *Config {
	collectorHost := os.Getenv("OPENTELEMETRY_COLLECTOR_HOST")

	c.Opentelemetry.CollectorHost = collectorHost

	return c
}

func (c *Config) jsonPlaceholderAPI() *Config {
	host := os.Getenv("JSON_PLACEHOLDER_API_HOST")
	c.JSONPlaceHolderAPI.Host = host

	return c
}

func (c *Config) mariadb() *Config {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	username := os.Getenv("MARIADB_USERNAME")
	password := os.Getenv("MARIADB_PASSWORD")
	database := os.Getenv("MARIADB_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_IDLE_CONNECTIONS"), 10, 64)

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	c.Mariadb.Driver = "mysql"
	c.Mariadb.Host = host
	c.Mariadb.Port = port
	c.Mariadb.Username = username
	c.Mariadb.Password = password
	c.Mariadb.Database = database
	c.Mariadb.DSN = dsn
	c.Mariadb.MaxOpenConnections = int(maxOpenConnections)
	c.Mariadb.MaxIdleConnections = int(maxIdleConnections)

	return c
}

func load() *Config {
	c := new(Config)
	return c.application().
		opentelemetry().jsonPlaceholderAPI().mariadb()
}

func Get() *Config {
	return load()
}
