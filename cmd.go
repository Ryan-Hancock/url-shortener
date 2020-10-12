package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()
var db = newDB(":memory:")

func initApp() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "URL shortening API"
	app.Usage = "Can shorten a full URL"
	app.Author = "Ryan Hancock"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		serveCommand(),
	}
}

func serveCommand() cli.Command {
	return cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Serve the web API",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Value:    "127.0.0.1:8000",
				Usage:    "set address for the api",
				Required: true,
			},
		},
		Action: func(c *cli.Context) {
			h := Handler{}
			h.initialise(db)
			h.run(c.String("address"))
		},
	}
}
