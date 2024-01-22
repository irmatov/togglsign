package storage_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/irmatov/togglsign/app"
	"github.com/irmatov/togglsign/infra/adapters/storage"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const dsnEnvName = "SIGNER_DSN"

func TestMain(m *testing.M) {
	if os.Getenv(dsnEnvName) == "" {
		os.Setenv(dsnEnvName, "postgres://postgres:mypass@localhost:5432/signer")
	}
	os.Exit(m.Run())
}

func testDB(t *testing.T) *pgx.Conn {
	c, err := pgx.Connect(context.Background(), os.Getenv(dsnEnvName))
	require.NoError(t, err)
	return c
}

func cleanUpTable(t *testing.T, c *pgx.Conn, name string) {
	_, err := c.Exec(context.Background(), "DELETE FROM "+name)
	assert.NoError(t, err)
}

// TestStorage verifies that SaveResponseSet and LoadResponseSet work correctly
func TestStorage(t *testing.T) {
	testConn := testDB(t)
	t.Cleanup(func() {
		cleanUpTable(t, testConn, "responses")
		cleanUpTable(t, testConn, "signatures")
	})
	tests := []app.ResponseSet{
		{
			Email:    "test@example.org",
			SignedAt: time.Date(2023, 9, 24, 1, 2, 3, 0, time.UTC),
			Sig:      "not a signature",
		},
		{
			Email:    "other@example.org",
			SignedAt: time.Date(2023, 10, 24, 1, 2, 3, 0, time.UTC),
			Sig:      "not a signature too",
			Responses: []app.Response{
				{Question: "q?", Answer: "a!"},
				{Question: "answer to everything?", Answer: "42"},
			},
		},
	}
	ctx := context.Background()
	s, err := storage.New(ctx, os.Getenv(dsnEnvName))
	require.NoError(t, err)
	for i, want := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.NoError(t, s.SaveResponseSet(ctx, want))
			got, err := s.LoadResponseSet(ctx, want.Email, want.Sig)
			assert.NoError(t, err)
			assert.Equal(t, want.Email, got.Email)
			assert.True(t, want.SignedAt.Equal(got.SignedAt), "want SignedAt=%v, got %v", want.SignedAt, got.SignedAt)
			assert.Equal(t, want.Sig, got.Sig)
			assert.Equal(t, want.Responses, got.Responses)
		})
	}
}
