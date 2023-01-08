package postgressql

import (
	"context"
	"fmt"
	"read-only_search/internal/domen/entity"
	client "read-only_search/pkg/client/postgresql"

	pb "github.com/i-b8o/read-only_contracts/pb/searcher/v1"
)

type searchStorage struct {
	client client.PostgreSQLClient
}

func NewSearchStorage(client client.PostgreSQLClient) *searchStorage {
	return &searchStorage{client: client}
}

const (
	general           = `SELECT d_id::INTEGER, d_name, c_id::INTEGER, c_name, p_id, text, count(*) OVER() AS full_count FROM doc_search WHERE ts @@ phraseto_tsquery('russian',$1) ORDER BY ts_rank(ts, phraseto_tsquery('russian',$1))`
	generalWithOffset = `SELECT d_id::INTEGER, d_name, c_id::INTEGER, c_name, p_id, text, count(*) OVER() AS full_count FROM doc_search WHERE ts @@ phraseto_tsquery('russian',$1) ORDER BY ts_rank(ts, phraseto_tsquery('russian',$1))`
)

func (ss *searchStorage) search(ctx context.Context, searchQuery, sql string) ([]*pb.SearchResponse, error) {
	fmt.Println(sql)
	var searchResults []*pb.SearchResponse

	rows, err := ss.client.Query(ctx, sql, searchQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		search := entity.Search{}
		if err = rows.Scan(
			&search.DocID, &search.DocName, &search.CID, &search.CName, &search.PID, &search.Text, &search.Count,
		); err != nil {
			return nil, err
		}
		resp := search.ToResponse()
		searchResults = append(searchResults, &resp)
	}

	return searchResults, nil
}

func (ss *searchStorage) Search(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	return ss.search(ctx, searchQuery, general)
}

func (ss *searchStorage) SearchWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	return ss.search(ctx, searchQuery, generalWithOffset)
}

func (ss *searchStorage) SearchLike(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT d_id::INTEGER, d_name, c_id::INTEGER, c_name, p_id, text, count(*) OVER() AS full_count from doc_search where text like '%` + searchQuery + `%'`

	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) SearchLikeWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT d_id::INTEGER, d_name, c_id::INTEGER, c_name, p_id, text, count(*) OVER() AS full_count from doc_search where text like '%` + searchQuery + `%'`
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) Docs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT id, name, title, updated_at, count(*) OVER() AS full_count from doc WHERE ts @@ phraseto_tsquery('russian',$1)`
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) DocsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT id, name, title, updated_at, count(*) OVER() AS full_count from doc WHERE ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	// sql += fmt.Sprintf(` AND (updated_at, id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY updated_at, id LIMIT %s`, params[0], params[1], params[2])
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)
	//  else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) Chapters(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT c.id, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from chapter AS c INNER JOIN doc as d ON c.r_id = d.id WHERE c.ts @@ phraseto_tsquery('russian',$1)`
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) ChaptersWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT c.id, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from chapter AS c INNER JOIN doc as d ON c.r_id = d.id WHERE c.ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	// sql += fmt.Sprintf(` AND (c.updated_at, c.id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY c.updated_at, c.id LIMIT %s`, params[0], params[1], params[2])
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)

	// else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) Paragraphs(ctx context.Context, searchQuery string) ([]*pb.SearchResponse, error) {
	sql := `SELECT p.id, p.content, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from paragraph AS p INNER JOIN chapter as c ON c.id = p.c_id INNER JOIN doc AS d ON c.d_id = d.id  WHERE p.ts @@ phraseto_tsquery('russian',$1)`
	return ss.search(ctx, searchQuery, sql)
}

func (ss *searchStorage) ParagraphsWithOffset(ctx context.Context, searchQuery string, offset, limit uint32) ([]*pb.SearchResponse, error) {
	sql := `SELECT p.id, p.content, c.name, d.name, c.updated_at, count(*) OVER() AS full_count from paragraph AS p INNER JOIN chapter as c ON c.id = p.c_id INNER JOIN doc AS d ON c.d_id = d.id  WHERE p.ts @@ phraseto_tsquery('russian',$1)`
	// Pagination
	// sql += fmt.Sprintf(` AND (c.updated_at, p.id) > ('%s' :: TIMESTAMPTZ, '%s') ORDER BY c.updated_at, p.id LIMIT %s`, params[0], params[1], params[2])
	sql += fmt.Sprintf(` OFFSET %d LIMIT %d`, offset, limit)

	// else if len(params) == 1 { // First page
	// 	sql += fmt.Sprintf(` LIMIT %s`, params[0])
	// }
	return ss.search(ctx, searchQuery, sql)
}
