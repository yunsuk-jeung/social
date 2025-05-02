package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `
		SELECT
				p.id,
				p.user_id,
				p.title,
				p.content,
				p.created_at,
				p.version,
				p.tags,
				u.username,
				count(c.id) AS comments_count
		FROM posts AS p
		LEFT JOIN comments AS c ON p.id = c.post_id
		LEFT JOIN users AS u ON p.user_id = u.id
		INNER JOIN followers AS f ON p.user_id = f.follower_id OR p.user_id = $1 
		WHERE 
			f.user_id = $1 
			AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
			AND ($5::varchar[] IS NULL OR p.tags @> $5)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + ` 
		LIMIT $2 OFFSET $3
`
	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}

	return feed, nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at, tags, version
		FROM posts
		WHERE id = $1
		`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

// Create
func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	err := s.db.QueryRowContext(ctx, query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// Delete
func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts WHERE id = $1
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	res, err := s.db.ExecContext(ctx, query, postID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// Update
func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version +1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancle := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancle()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
