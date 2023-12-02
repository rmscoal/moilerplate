package main

import (
	"github.com/rmscoal/go-restful-monolith-boilerplate/cmd"
	"github.com/rmscoal/go-restful-monolith-boilerplate/docs"
)

func main() {
	// Swagger documentation info
	docs.SwaggerInfo.Title = "Moilerplate"
	docs.SwaggerInfo.Description = "A monolithic RESTful API for Go"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8082"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Execute the app
	cmd.Execute()
}
