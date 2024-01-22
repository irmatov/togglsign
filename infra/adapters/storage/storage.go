package storage

import (
	"context"
	"errors"

	"github.com/irmatov/togglsign/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

const insertSignatureQuery = `
INSERT INTO signatures (email, signed_at, sig)
VALUES ($1, $2, $3)
RETURNING id`

const insertResponseQuery = `
INSERT INTO responses (sort_key, question, answer, sig_id)
VALUES ($1, $2, $3, $4)`

func New(ctx context.Context, dsn string) (*Storage, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{pool}, nil
}

func (s *Storage) SaveResponseSet(ctx context.Context, rs app.ResponseSet) error {
	c, err := s.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()
	tx, err := c.Begin(ctx)
	if err != nil {
		return err
	}
	if err = func() error {
		var id int
		if err := tx.QueryRow(ctx, insertSignatureQuery, rs.Email, rs.SignedAt, rs.Sig).Scan(&id); err != nil {
			return err
		}

		for i, r := range rs.Responses {
			if _, err := tx.Exec(ctx, insertResponseQuery, i, r.Question, r.Answer, id); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		_ = tx.Rollback(ctx) // FIXME: at least we need to log the error
		return err
	}
	return tx.Commit(ctx)
}

const fetchSignatureQuery = `
SELECT id, signed_at
FROM signatures
WHERE email = $1 AND sig = $2`

const fetchResponsesQuery = `
SELECT question, answer
FROM responses
WHERE sig_id = $1
ORDER BY sort_key`

func (s *Storage) LoadResponseSet(ctx context.Context, email, sig string) (*app.ResponseSet, error) {
	c, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Release()
	tx, err := c.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	rs := app.ResponseSet{Email: email, Sig: sig}
	var id int
	if err := tx.QueryRow(ctx, fetchSignatureQuery, email, sig).Scan(&id, &rs.SignedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	rows, err := tx.Query(ctx, fetchResponsesQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r app.Response
		if err := rows.Scan(&r.Question, &r.Answer); err != nil {
			return nil, err
		}
		rs.Responses = append(rs.Responses, r)
	}
	rows.Close()
	return &rs, rows.Err()
}
