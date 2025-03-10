package batch

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupPG() {
	urlExample := "postgres://leow:password@127.0.0.1:5432/pgtemporal"
	dbpool, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var greeting string
	err = dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)
}

func TestAddRole(t *testing.T) {
	type args struct {
		username string
		role     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	setupPG()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddRole(tt.args.username, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
