ALTER TABLE doc
ADD COLUMN ts tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('russian', coalesce(name, '')), 'A') || setweight(
            to_tsvector('russian', coalesce(header, '')),
            'B'
        )
    ) STORED;
CREATE INDEX doc_ts_idx ON doc USING GIN (ts);
ALTER TABLE chapter
ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', name)) STORED;
CREATE INDEX ch_ts_idx ON chapter USING GIN (ts);
ALTER TABLE paragraph
ADD COLUMN ts tsvector GENERATED ALWAYS AS (to_tsvector('russian', content)) STORED;
CREATE INDEX p_ts_idx ON paragraph USING GIN (ts);
CREATE MATERIALIZED VIEW doc_search AS
SELECT d.id AS "d_id",
    d.name AS "d_name",
    NULL AS "c_id",
    NULL AS "c_name",
    CAST(NULL AS integer) AS "p_id",
    NULL AS "p_text",
    d.name AS "text",
    to_tsvector('russian', d.name) AS ts
FROM doc AS d
UNION
SELECT NULL AS "d_id",
    d.name AS "d_name",
    c.id AS "c_id",
    c.name AS "c_name",
    NULL AS "p_id",
    NULL AS "p_text",
    c.name AS "text",
    to_tsvector('russian', c.name) AS ts
FROM chapter AS c
    INNER JOIN doc AS d ON d.id = c.doc_id
UNION
SELECT NULL AS "d_id",
    d.name AS "d_name",
    c.id AS "c_id",
    c.name AS "c_name",
    p.paragraph_id AS "p_id",
    p.content AS "p_text",
    p.content AS "text",
    to_tsvector('russian', content) AS ts
FROM paragraph AS p
    INNER JOIN chapter AS c ON p.c_id = c.id
    INNER JOIN doc AS d ON c.doc_id = d.id;
create index idx_search on doc_search using GIN(ts);

GRANT CONNECT ON DATABASE search TO reader;
GRANT SELECT ON TABLE public.doc TO reader;
GRANT SELECT ON TABLE public.chapter TO reader;
GRANT SELECT ON TABLE public.paragraph TO reader;
GRANT SELECT ON public.doc_search TO reader;
