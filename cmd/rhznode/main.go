package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/drgomesp/acervo/internal/rhznode"
	"github.com/drgomesp/acervo/pkg/node"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	app := &cli.App{
		Name:  "rhznode",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "moniker",
				Value: "marley",
			},
		},
		Action: func(c *cli.Context) (err error) {
			var fullNode *node.Node

			if fullNode, err = rhznode.NewFullNode(c.String("moniker")); err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}
			if err != nil {
				return errors.Wrap(err, "failed to initialize full node")
			}
			_ = fullNode

			return startNode(c.Context, fullNode)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(errors.Wrap(err, "failed to run full node")).Send()
	}
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
		log.Info().Msg("interrupt signal, shutting down...")
	}()

	return nil
}