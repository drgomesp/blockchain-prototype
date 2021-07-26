package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rhizomplatform/rhizom/internal/rhznode"
	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/urfave/cli/v2"
)

func init() {
}

func main() {
	app := &cli.App{
		Name:  "rhznode",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) (err error) {
			fullNode, err := makeFullNode()
			if err != nil {
				return err
			}

			return fullNode.Start()
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(errors.Wrap(err, "failed run rhz full node"))
	}
}

func makeFullNode() (*rhznode.FullNode, error) {
	n, err := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	var rhz *rhznode.FullNode

	if rhz, err = rhznode.NewFullNode(n); err != nil {
		return nil, errors.Wrap(err, "failed to initialize rhznode")
	}

	if err = rhz.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start rhznode")
	}

	return rhz, nil
}
