package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/drgomesp/rhizom/internal/rhznode"
	"github.com/pkg/errors"
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
			ctx, cancelFunc := context.WithCancel(c.Context)
			defer cancelFunc()

			fullNode, err := makeFullNode(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}

			return errors.Wrap(fullNode.Start(ctx), "failed to run full node")
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(errors.Wrap(err, "failed to run full node"))
	}
}

func buildLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Kitchen)

	logger, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialized logger")
	}

	return logger, nil
}

func makeFullNode(ctx context.Context) (*rhznode.FullNode, error) {
	logger, err := buildLogger()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}

	var rhz *rhznode.FullNode

	if rhz, err = rhznode.NewFullNode(logger.Sugar()); err != nil {
		return nil, errors.Wrap(err, "failed to initialize full node")
	}

	if err = rhz.Start(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to start full node")
	}

	return rhz, nil
}
