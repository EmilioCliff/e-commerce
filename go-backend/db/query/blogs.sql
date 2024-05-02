-- name: CreateBlog :one
INSERT INTO blogs (
    author, title, content
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetBlog :one
SELECT * FROM blogs
WHERE id = $1;

-- name: GetBlogForUpdate :one
SELECT * FROM blogs
WHERE id = $1
FOR NO KEY UPDATE;

-- name: GetAdminsBlog :many
SELECT * FROM blogs
WHERE author = $1;

-- name: ListBlogs :many
SELECT * FROM blogs
ORDER BY title;

-- name: EditBlog :one
UPDATE blogs
    set title = $1,
    content = $2
WHERE id = $3
RETURNING *;

-- name: DeleteBlog :exec
DELETE FROM blogs
WHERE id = $1;