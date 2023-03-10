package service

import (
	"context"
	"fmt"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type GeneralStorage interface {
	Search(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	SearchWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
	SearchLike(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error)
	SearchLikeWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error)
}

type generalService struct {
	storage GeneralStorage
	logger  logging.Logger
}

func NewGeneralService(storage GeneralStorage, logger logging.Logger) *generalService {
	return &generalService{storage: storage, logger: logger}
}

func (s generalService) Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error) {
	if len(params) == 2 {
		respSlice, err := s.storage.SearchWithOffset(ctx, searchQuery, params[0], params[1])
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		if len(respSlice) == 0 {
			respSlice, err = s.storage.SearchLikeWithOffset(ctx, searchQuery, params[0], params[1])
			if err != nil {
				s.logger.Error(err)
				return nil, err
			}
			return respSlice, nil
		}
	}

	respSlice, err := s.storage.Search(ctx, searchQuery)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	if len(respSlice) == 0 {
		respSlice, err = s.storage.SearchLike(ctx, searchQuery)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
	}
	fmt.Println("AOK")
	return respSlice, nil
}
