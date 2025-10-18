-- name: GetAllPhotosWithTags :many
SELECT p.*, 
       COALESCE(
           JSONB_AGG(
               JSONB_BUILD_OBJECT(
                   'title', t.title,
                   'comment', t.comment
               )
           ) FILTER (WHERE t.title IS NOT NULL), '[]'
       )::jsonb AS tags
FROM photos p
LEFT JOIN photo_tags pt ON p.id = pt.photo_id
LEFT JOIN tags t ON pt.tag_title = t.title
GROUP BY p.id;

-- name: GetPhotoByIDWithTags :one
SELECT p.*, 
       COALESCE(
           JSONB_AGG(
               JSONB_BUILD_OBJECT(
                   'title', t.title,
                   'comment', t.comment
               )
           ) FILTER (WHERE t.title IS NOT NULL), '[]'
       )::json AS tags
FROM photos p
LEFT JOIN photo_tags pt ON p.id = pt.photo_id
LEFT JOIN tags t ON pt.tag_title = t.title
WHERE p.id = $1
GROUP BY p.id;

-- name: GetPhotosByTagTitle :many
SELECT p.*,
         COALESCE(
              JSONB_AGG(
                JSONB_BUILD_OBJECT(
                     'title', t2.title,
                     'comment', t2.comment
                )
              ) FILTER (WHERE t2.title IS NOT NULL), '[]'
         )::json AS tags
FROM photos p
LEFT JOIN photo_tags pt ON p.id = pt.photo_id
LEFT JOIN tags t ON pt.tag_title = t.title
WHERE t.title = $1
GROUP BY p.id;

-- name: GetPhotosByTagTitles :many
SELECT p.*,
            COALESCE(
                JSONB_AGG(
                    JSONB_BUILD_OBJECT(
                        'title', t2.title,
                        'comment', t2.comment
                    )
                ) FILTER (WHERE t2.title IS NOT NULL), '[]'
            )::json AS tags
FROM photos p
LEFT JOIN photo_tags pt ON p.id = pt.photo_id
LEFT JOIN tags t ON pt.tag_title = t.title
WHERE t.title = ANY($1)
GROUP BY p.id;

-- name: ListTagsWithCount :many
SELECT t.title,
       t.comment,
       COUNT(pt.photo_id) AS photo_count
FROM tags t
LEFT JOIN photo_tags pt ON t.title = pt.tag_title
GROUP BY t.title, t.comment;


-- name: AddPhoto :one
INSERT INTO photos (title, photo_url, comment, metadata) 
VALUES ($1, $2, $3, $4) 
RETURNING id;

-- name: CreateTag :exec
INSERT INTO tags (title, comment) VALUES ($1, $2);

-- name: DeleteTag :exec
DELETE FROM tags WHERE title = $1;

-- name: AddTagToPhoto :exec
INSERT INTO photo_tags (photo_id, tag_title) VALUES ($1, $2);

-- name: AddTagsToPhoto :exec
INSERT INTO photo_tags (photo_id, tag_title)
SELECT $1, unnest($2::text[]);

-- name: RemoveTagFromPhoto :exec
DELETE FROM photo_tags WHERE photo_id = $1 AND tag_title = $2;