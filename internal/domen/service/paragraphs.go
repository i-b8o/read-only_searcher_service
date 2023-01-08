package service

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type ParagraphStorage interface {
	Paragraphs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	ParagraphsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type paragraphsService struct {
	storage ParagraphStorage
	logger  logging.Logger
}

func NewParagraphsService(storage ParagraphStorage, logger logging.Logger) *paragraphsService {
	return &paragraphsService{storage: storage, logger: logger}
}

func (s paragraphsService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		return s.storage.ParagraphsWithOffset(ctx, searchQuery, params[0], params[1])
	}
	return s.storage.Paragraphs(ctx, searchQuery)
}
