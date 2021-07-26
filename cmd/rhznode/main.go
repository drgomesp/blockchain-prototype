package main

import (
	"log"
	"os"

	"github.com/rhizomplatform/rhizom/internal/rhznode"

	"github.com/pkg/errors"
	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/urfave/cli/v2"
)

func init() {
}

func main() {
	app := &cli.App{
		Name:  "rhznode",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			log.Println("Hello world!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		node, err := node.New(&node.Config{
			Type: node.TypeFull,
			Name: "rhz_node",
		})
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to initialize node"))
		}

		var rhz *rhznode.FullNode
		if rhz, err = rhznode.New(node); err != nil {
			log.Fatal(errors.Wrap(err, "failed to initialize rhznode"))
		}

		rhz.Start()
	}
}
