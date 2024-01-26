package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresSQLPool struct {
	pool *pgxpool.Pool
}

func NewPostgresSQLPool(url string) (*PostgresSQLPool, error) {
	pool, err := pgxpool.Connect(context.Background(), url)

	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgresSQLPool{
		pool: pool,
	}, nil
}

const (
	insertUserTemplate = `INSERT INTO users (id) VALUES ('%v')`
	deleteUserTemplate = `DELETE FROM users WHERE id = '%v'`
	selectUsersID      = `SELECT id FROM users`
)

func (p *PostgresSQLPool) AddUser(ctx context.Context, ID int64) error {
	q := fmt.Sprintf(insertUserTemplate, ID)
	rows, err := p.pool.Query(ctx, q)
	defer rows.Close()

	return err
}

func (p *PostgresSQLPool) RemoveUser(ctx context.Context, ID int64) error {
	q := fmt.Sprintf(deleteUserTemplate, ID)
	rows, err := p.pool.Query(ctx, q)
	defer rows.Close()

	return err
}

func (p *PostgresSQLPool) GetUsersID(ctx context.Context) ([]int64, error) {
	rows, err := p.pool.Query(ctx, selectUsersID)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var ids []int64
	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
