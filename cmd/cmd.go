package cmd

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/rmscoal/go-restful-monolith-boilerplate/cmd/app"
)

func Execute() {
	cli := cli.NewCLI("moilerplate", "1.0.0")
	cli.Commands = Commands()
	cli.Args = os.Args[1:]

	exitCode, err := cli.Run()
	if err != nil {
		log.Fatalf("unable to run app")
	}

	os.Exit(exitCode)
}

func Commands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return app.NewAppCLI(), nil
		},
	}
}
