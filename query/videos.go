package query

const SqlSearch = `
SELECT 
videos.id,
videos.title,
videos.short_desc,
videos.image_url,
videos.video_url,
rank_title,
rank_desc,
similarity
FROM 
videos, 
to_tsvector(videos.title || videos.short_desc) document,
to_tsquery($1) query,
NULLIF(ts_rank(to_tsvector(videos.title), query), 0) rank_title,
NULLIF(ts_rank(to_tsvector(videos.short_desc), query), 0) rank_desc,
SIMILARITY($1, videos.title || videos.short_desc) similarity
WHERE query @@ document OR similarity > 0
ORDER BY rank_title, rank_desc, similarity DESC NULLS LAST
`
