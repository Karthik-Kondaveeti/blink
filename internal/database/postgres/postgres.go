/*
	Fill the required environment variables in the .env file

	DB_URL, TABLE_NAME
*/

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	insertQuery string
	selectQuery string
	db          *pgxpool.Pool
}

func New(connStr string, tablename string) (*Postgres, error) {
	if connStr == "" || tablename == "" {
		return nil, errors.New("Connection string or Table name is empty")
	}

	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Postgres{
		db: db,
		insertQuery: fmt.Sprintf(
			`INSERT INTO %s (short_code, original_url)
			VALUES ($1, $2)`,
			tablename,
		),
		selectQuery: fmt.Sprintf(
			`SELECT original_url FROM %s WHERE short_code = $1`,
			tablename,
		),
	}, nil
}

func (p *Postgres) AddLink(ctx context.Context, shortCode string, originalURL string) error {
	_, err := p.db.Exec(ctx, p.insertQuery, shortCode, originalURL)
	if err != nil {
		return fmt.Errorf("Add Link: %w", err)
	}
	return nil
}

func (p *Postgres) GetLink(ctx context.Context, shortCode string) (string, error) {
	// for testing
	// time.Sleep(500 * time.Millisecond)

	var originalURL string
	err := p.db.QueryRow(ctx, p.selectQuery, shortCode).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("Get Link: %w", err)
	}
	return originalURL, nil
}

func (p *Postgres) Close() error {
	p.db.Close()
	return nil
}
