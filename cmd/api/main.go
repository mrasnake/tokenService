package main

import (
	"fmt"
	"log"
	"mrasnake/TokenService/internal"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

// main takes all build flags or environment variables, defines
// the configuration and runs the client.
func main() {

	app := cli.NewApp()
	app.Name = "server"
	app.Usage = "instanciate a new server to fulfill requests"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "logs",
			Aliases: []string{"l"},
			Value:   fmt.Sprintf("./logfile-%v.log", time.Now().String()),
			EnvVars: []string{"LOG_FILE"},
		},
	}

	app.Action = func(ctx *cli.Context) error {

		service, err := internal.NewService(ctx.String("logs"))
		if err != nil {
			return fmt.Errorf("unable to start service: %w", err)
		}
		server := internal.NewServer(service)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error occurred %v", err)
	}
}
