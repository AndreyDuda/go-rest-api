package author

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"rest_api/internal/author"
	"rest_api/pkg/client/postgresql"
	"rest_api/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", "")
}

func (r repository) Create(ctx context.Context, author *author.Author) error {
	q := `
		  INSERT INTO author (name) 
		  VALUES ($1)
		  RETURNING id
		  `

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, author.Name).Scan(&author.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); !ok {
			return pgErr
		}

		return err
	}

	return nil
}

func (r repository) FindAll(ctx context.Context) ([]author.Author, error) {
	q := `SELECT id, name FROM public.author`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	authors := make([]author.Author, 0)

	for rows.Next() {
		var auth author.Author

		err = rows.Scan(&auth.ID, &auth.Name)
		if err != nil {
			return nil, err
		}

		authors = append(authors, auth)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r repository) FindOne(ctx context.Context, id string) (author.Author, error) {
	var auth author.Author
	q := `SELECT id, name FROM public.author WHERE id = $1`

	err := r.client.QueryRow(ctx, q, id).Scan(&auth.ID, &auth.Name)
	if err != nil {
		return author.Author{}, err
	}

	return auth, nil
}

func (r repository) Update(ctx context.Context, user author.Author) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) author.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
