package service

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type DocStorage interface {
	Docs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	DocsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type docService struct {
	storage DocStorage
	logger  logging.Logger
}

func NewDocService(storage DocStorage, logger logging.Logger) *docService {
	return &docService{storage: storage, logger: logger}
}

func (s docService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		return s.storage.DocsWithOffset(ctx, searchQuery, params[0], params[1])
	}
	return s.storage.Docs(ctx, searchQuery)
}
