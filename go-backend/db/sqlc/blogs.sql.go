// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: blogs.sql

package db

import (
	"context"
)

const createBlog = `-- name: CreateBlog :one
INSERT INTO blogs (
    author, title, content
) VALUES (
    $1, $2, $3
)
RETURNING id, author, title, content, created_at
`

type CreateBlogParams struct {
	Author  int64  `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (q *Queries) CreateBlog(ctx context.Context, arg CreateBlogParams) (Blog, error) {
	row := q.db.QueryRow(ctx, createBlog, arg.Author, arg.Title, arg.Content)
	var i Blog
	err := row.Scan(
		&i.ID,
		&i.Author,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBlog = `-- name: DeleteBlog :exec
DELETE FROM blogs
WHERE id = $1
`

func (q *Queries) DeleteBlog(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteBlog, id)
	return err
}

const editBlog = `-- name: EditBlog :one
UPDATE blogs
    set title = $1,
    content = $2
WHERE id = $3
RETURNING id, author, title, content, created_at
`

type EditBlogParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ID      int64  `json:"id"`
}

func (q *Queries) EditBlog(ctx context.Context, arg EditBlogParams) (Blog, error) {
	row := q.db.QueryRow(ctx, editBlog, arg.Title, arg.Content, arg.ID)
	var i Blog
	err := row.Scan(
		&i.ID,
		&i.Author,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const getAdminsBlog = `-- name: GetAdminsBlog :many
SELECT id, author, title, content, created_at FROM blogs
WHERE author = $1
`

func (q *Queries) GetAdminsBlog(ctx context.Context, author int64) ([]Blog, error) {
	rows, err := q.db.Query(ctx, getAdminsBlog, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Blog
	for rows.Next() {
		var i Blog
		if err := rows.Scan(
			&i.ID,
			&i.Author,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBlog = `-- name: GetBlog :one
SELECT id, author, title, content, created_at FROM blogs
WHERE id = $1
`

func (q *Queries) GetBlog(ctx context.Context, id int64) (Blog, error) {
	row := q.db.QueryRow(ctx, getBlog, id)
	var i Blog
	err := row.Scan(
		&i.ID,
		&i.Author,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const getBlogForUpdate = `-- name: GetBlogForUpdate :one
SELECT id, author, title, content, created_at FROM blogs
WHERE id = $1
FOR NO KEY UPDATE
`

func (q *Queries) GetBlogForUpdate(ctx context.Context, id int64) (Blog, error) {
	row := q.db.QueryRow(ctx, getBlogForUpdate, id)
	var i Blog
	err := row.Scan(
		&i.ID,
		&i.Author,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const listBlogs = `-- name: ListBlogs :many
SELECT id, author, title, content, created_at FROM blogs
ORDER BY created_at
`

func (q *Queries) ListBlogs(ctx context.Context) ([]Blog, error) {
	rows, err := q.db.Query(ctx, listBlogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Blog
	for rows.Next() {
		var i Blog
		if err := rows.Scan(
			&i.ID,
			&i.Author,
			&i.Title,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
