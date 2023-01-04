package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"mrasnake/TokenService/internal"
	"net/http"
	"os"
)

// main takes all build flags or environment variables, defines
// the configuration and runs the client.
func main() {

	app := cli.NewApp()
	app.Name = "server"
	app.Usage = "instantiate a new server to fulfill requests"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   ":10000",
			EnvVars: []string{"PORT"},
		},
	}

	app.Action = func(ctx *cli.Context) error {

		port := ctx.String("port")

		service, err := internal.NewService()
		if err != nil {
			return fmt.Errorf("unable to start service: %w", err)
		}
		server := internal.NewServer(service)
		server.InitServer()

		http.Handle("/", server.Router)
		fmt.Println("Starting token service at ", port)
		return http.ListenAndServe(port, nil)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error occurred %v", err)
	}
}
