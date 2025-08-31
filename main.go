package main

import (
	"context"

	microservice "github.com/rlapenok/toolbox/micro_service"
)

type Config struct{}

func (c *Config) GetName() string {
	return "microservice"
}

func (c *Config) EnableDefaultGinServer() bool {
	return false
}

func main() {
	service := microservice.NewMicroService(&Config{})
	service.Run(context.Background())
}
