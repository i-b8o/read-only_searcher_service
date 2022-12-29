package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type DocStorage interface {
	Docs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	DocsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type docService struct {
	storage DocStorage
}

func NewDocService(storage DocStorage) *docService {
	return &docService{storage: storage}
}

func (s docService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		return s.storage.DocsWithOffset(ctx, searchQuery, params[0], params[1])
	}
	return s.storage.Docs(ctx, searchQuery)
}
