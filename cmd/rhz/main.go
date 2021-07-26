package main

import (
	"fmt"
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
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
