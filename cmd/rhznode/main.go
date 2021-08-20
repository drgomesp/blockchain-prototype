package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drgomesp/rhizom/internal/rhznode"
	"github.com/drgomesp/rhizom/pkg/node"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var app *cli.App

func init() {
	app = &cli.App{
		Name:  "rhznode",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) (err error) {
			ctx, cancelFunc := context.WithCancel(c.Context)
			defer cancelFunc()

			fullNode, err := makeFullNode(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}
			_ = fullNode

			return startNode(ctx, fullNode)
		},
	}
}

func main() {
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

func makeFullNode(ctx context.Context) (*node.Node, error) {
	logger, err := buildLogger()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize logger")
	}

	var fullNode *node.Node

	if fullNode, err = rhznode.NewFullNode(logger.Sugar()); err != nil {
		return nil, errors.Wrap(err, "failed to initialize full node")
	}

	return fullNode, nil
}

func startNode(ctx context.Context, node *node.Node) error {
	if err := node.Start(ctx); err != nil {
		return err
	}

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		defer signal.Stop(sig)

		<-sig
		log.Println("interrupt signal, shutting down...")
	}()

	return nil
}
