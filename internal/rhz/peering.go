package rhz

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/drgomesp/acervo/pkg/block"
)

type PeeringService struct {
}

func NewPeeringService() *PeeringService {
	return &PeeringService{}
}

func (p *PeeringService) GetBlocks(_ context.Context, index uint64) ([]block.Block, error) {
	return []block.Block{
		{Header: block.Header{Index: index}},
	}, nil
}

func (p *PeeringService) Blocks(_ context.Context, blocks []block.Block) error {
	log.Debug().Msgf("Blocks(blocks=%vs)", blocks)

	return nil
}