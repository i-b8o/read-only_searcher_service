package postgressql

import (
	"context"
	"fmt"
	client "read-only_search/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type generalStorage struct {
	client client.PostgreSQLClient
}

func NewGeneralStorage(client client.PostgreSQLClient) *generalStorage {
	return &generalStorage{client: client}
}

func (ss *generalStorage) General(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT d_id, d_name, c_id, c_name, p_id, text, count(*) OVER() AS full_count FROM doc_search WHERE ts @@ phraseto_tsquery('russian',$1) ORDER BY ts_rank(ts, phraseto_tsquery('russian',$1))`

	var searchResults []*pb.SearchResponse

	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.ChapterID, &search.ChapterName, &search.ParagraphID, &search.Text, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (ss *generalStorage) SearchWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT d_id, d_name, c_id, c_name, p_id, text, count(*) OVER() AS full_count FROM doc_search WHERE ts @@ phraseto_tsquery('russian',$1) ORDER BY ts_rank(ts, phraseto_tsquery('russian',$1))`

	var searchResults []*pb.SearchResponse
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)

	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.ChapterID, &search.ChapterName, &search.ParagraphID, &search.Text, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (ss *generalStorage) SearchLike(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT r_id, r_name, c_id, c_name, p_id, text, count(*) OVER() AS full_count from doc_search where text like '%` + searchQuery + `%'`

	var searchResults []*pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		search := &pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.ChapterID, &search.ChapterName, &search.ParagraphID, &search.Text, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}

func (ss *generalStorage) SearchLikeWithOffset(ctx context.Context, searchQuery string, params ...string) ([]pb.SearchResponse, error) {
	sql := `SELECT r_id, r_name, c_id, c_name, p_id, text, count(*) OVER() AS full_count from doc_search where text like '%` + searchQuery + `%'`
	if len(params) == 2 {
		sql += fmt.Sprintf(` OFFSET %s LIMIT %s`, params[0], params[1])
	}
	var searchResults []pb.SearchResponse
	rows, err := ss.client.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		search := pb.SearchResponse{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.ChapterID, &search.ChapterName, &search.ParagraphID, &search.Text, &search.Count,
		); err != nil {
			return nil, err
		}

		searchResults = append(searchResults, search)
	}

	return searchResults, nil
}
