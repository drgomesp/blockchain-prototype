package rhz

import (
	"context"

	"github.com/drgomesp/rhizom/pkg/block"
	"go.uber.org/zap"
)

type PeeringService struct {
	logger *zap.SugaredLogger
}

func NewPeeringService(logger *zap.SugaredLogger) *PeeringService {
	return &PeeringService{
		logger: logger,
	}
}

func (p *PeeringService) GetBlocks(_ context.Context, index uint64) ([]*block.Block, error) {
	return []*block.Block{
		{Header: block.Header{Index: index}},
	}, nil
}

func (p *PeeringService) Blocks(_ context.Context, blocks []*block.Block) error {
	p.logger.Debugf("blocks: %v", blocks)

	return nil
}
