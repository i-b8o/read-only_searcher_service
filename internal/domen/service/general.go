package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type GeneralStorage interface {
	Paragraphs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	ParagraphsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type generalService struct {
	storage GeneralStorage
}

func NewGeneralService(storage GeneralStorage) *generalService {
	return &generalService{storage: storage}
}

func (s generalService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		return s.storage.ParagraphsWithOffset(ctx, searchQuery, params[0], params[1])
	}
	return s.storage.Paragraphs(ctx, searchQuery)
}
