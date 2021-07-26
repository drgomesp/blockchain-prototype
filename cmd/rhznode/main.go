package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rhizomplatform/rhizom/internal/rhznode"
	"github.com/rhizomplatform/rhizom/pkg/node"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
}

func main() {
	app := &cli.App{
		Name:  "rhznode",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) (err error) {
			config := zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			logger, err := config.Build()
			if err != nil {
				return errors.Wrap(err, "failed to initialized logger")
			}

			fullNode, err := makeFullNode(logger.Sugar())
			if err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}

			return errors.Wrap(fullNode.Start(), "failed to run full node")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(errors.Wrap(err, "failed to run full node"))
	}
}

func makeFullNode(logger *zap.SugaredLogger) (*rhznode.FullNode, error) {
	n, err := node.New(&node.Config{
		Type: node.TypeFull,
		Name: "rhz_node",
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed initialize node")
	}

	var rhz *rhznode.FullNode

	if rhz, err = rhznode.NewFullNode(logger, n); err != nil {
		return nil, errors.Wrap(err, "failed to initialize full node")
	}

	if err = rhz.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start full node")
	}

	return rhz, nil
}
