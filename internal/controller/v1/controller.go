package controller

import (
	"context"

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

func (s *DocGRPCService) Docs(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponseMessage, error) {
	return s.search(ctx, req, s.docService)
}

func (s *DocGRPCService) Chapters(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponseMessage, error) {
	return s.search(ctx, req, s.paragraphService)
}

func (s *DocGRPCService) Pargaraphs(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponseMessage, error) {
	return s.search(ctx, req, s.paragraphService)
}

func (s *DocGRPCService) General(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponseMessage, error) {
	return s.search(ctx, req, s.generalService)
}

func (s *DocGRPCService) search(ctx context.Context, req *pb.SearchRequest, searchservice SearchService) (*pb.SearchResponseMessage, error) {
	ofset := req.GetOffset()
	limit := req.GetLimit()
	query := req.GetSearchQuery()
	r, err := searchservice.Search(ctx, query, ofset, limit)
	if err != nil {
		return nil, err
	}
	resp := &pb.SearchResponseMessage{Response: r}
	return resp, nil
}
