package controller

import (
	"context"
	"fmt"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type SearchService interface {
	Search(ctx context.Context, searchQuery string, params ...uint32) ([]*pb.SearchResponse, error)
}

type DocGRPCService struct {
	docService       SearchService
	chapterService   SearchService
	paragraphService SearchService
	generalService   SearchService
	pb.UnimplementedSearcherGRPCServer
}

func NewDocGRPCService(docService, chapterService, paragraphService, generalService SearchService) *DocGRPCService {
	return &DocGRPCService{
		docService:       docService,
		chapterService:   chapterService,
		paragraphService: paragraphService,
		generalService:   generalService,
	}
}

func (s *DocGRPCService) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponseMessage, error) {

	switch subj := req.GetSubject(); subj {
	case pb.SearchRequest_Docs:
		return s.search(ctx, req, s.docService)
	case pb.SearchRequest_Chapters:
		return s.search(ctx, req, s.chapterService)
	case pb.SearchRequest_Pargaraphs:
		return s.search(ctx, req, s.paragraphService)
	case pb.SearchRequest_General:
		return s.search(ctx, req, s.generalService)
	default:
		return nil, fmt.Errorf("wrong subject")
	}
}

func (s *DocGRPCService) search(ctx context.Context, req *pb.SearchRequest, searchservice SearchService) (*pb.SearchResponseMessage, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()
	query := req.GetSearchQuery()
	r, err := searchservice.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, err
	}
	resp := &pb.SearchResponseMessage{Response: r}
	return resp, nil
}
