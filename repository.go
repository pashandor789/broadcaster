package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresSQLPool struct {
	pool   *pgxpool.Pool
}

const (
	insertUserTemplate = `INSERT INTO users (id) VALUES ('%v')`
	deleteUserTemplate = `DELETE FROM users WHERE id = '%v'`
)

func (p *PostgresSQLPool) AddUser(ctx context.Context, ID int64) error {
	cmd = fmt.Sprintf(insertUserDataTemplate, ID)

	rows, err := p.pool.Query(ctx, cmd)
	defer rows.Close() 

	return err
}

func (p *PostgresSQLPool) RemoveUser(ctx context.Context, ID int64) error {
	cmd = fmt.Sprintf(deleteUserTemplate, ID)

	rows, err := p.pool.
}