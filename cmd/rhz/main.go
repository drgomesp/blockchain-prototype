package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
}

func main() {
	app := &cli.App{
		Name:  "rhz",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			log.Println("Hello world!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
