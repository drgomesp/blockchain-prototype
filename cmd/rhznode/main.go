package main

import (
	"context"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rhizomplatform/rhizom/internal/rhznode"
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

			fullNode, err := makeFullNode(c.Context, logger.Sugar())
			if err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}

			return errors.Wrap(fullNode.Start(c.Context), "failed to run full node")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(errors.Wrap(err, "failed to run full node"))
	}
}

func makeFullNode(ctx context.Context, logger *zap.SugaredLogger) (
	rhz *rhznode.FullNode, err error,
) {
	if rhz, err = rhznode.NewFullNode(logger); err != nil {
		return nil, errors.Wrap(err, "failed to initialize full node")
	}

	if err = rhz.Start(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to start full node")
	}

	return rhz, nil
}
