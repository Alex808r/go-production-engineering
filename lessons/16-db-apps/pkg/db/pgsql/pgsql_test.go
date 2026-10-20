package pgsql_test

import (
	"context"
	"testing"
	"time"

	"go-core-4/lessons/16-db-apps/pkg/db"
	"go-core-4/lessons/16-db-apps/pkg/db/pgsql"
)

func TestInterfaceCompliance(t *testing.T) {
	var _ db.Interface = (*pgsql.DB)(nil)
}

func TestNilPoolErrors(t *testing.T) {
	d := pgsql.New(nil)

	_, err := d.Books(context.Background())
	if err == nil {
		t.Fatal("expected error on nil pool, got nil")
	}

	err = d.AddBooks(context.Background(), []db.Book{{Title: "Test"}})
	if err == nil {
		t.Fatal("expected error on nil pool, got nil")
	}

	// Empty books slice returns nil without accessing pool
	err = d.AddBooks(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected nil for empty slice, got: %v", err)
	}

	d.Close()
}

func TestConnect_Unavailable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err := pgsql.Connect(ctx, "postgres://user:pass@127.0.0.1:54329/invalid_db?connect_timeout=1")
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}
