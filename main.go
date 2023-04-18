package main

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/cmd/app"
	"github.com/rmscoal/go-restful-monolith-boilerplate/config"
)

func main() {
	cfg := config.GetConfig()

	app.Run(cfg)
}
