package service

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type ChapterStorage interface {
	Chapters(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	ChaptersWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type chapterService struct {
	storage ChapterStorage
	logger  logging.Logger
}

func NewChapterService(storage ChapterStorage, logger logging.Logger) *chapterService {
	return &chapterService{storage: storage, logger: logger}
}

func (s chapterService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		return s.storage.ChaptersWithOffset(ctx, searchQuery, params[0], params[1])
	}
	return s.storage.Chapters(ctx, searchQuery)
}
