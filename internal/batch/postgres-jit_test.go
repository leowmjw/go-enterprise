package batch

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	// Check no access ..

	// Change permission
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Error(err)
	}
	defer conn.Close(context.Background())
	// grant ..
	ct, yerr := conn.Exec(context.Background(), `GRANT $1 TO $2`, "read", "mleow")
	if yerr != nil {
		var pgErr *pgconn.PgError
		if errors.As(yerr, &pgErr) {
			fmt.Println("ERR:", pgErr.Message) // => syntax error at end of input
			fmt.Println("CODE:", pgErr.Code)   // => 42601
		}
		t.Fatal(yerr)
	}
	spew.Dump(ct.RowsAffected()) // aalways 0
	spew.Dump(ct.String())

	// Revoke ..
	ct, xerr := conn.Exec(context.Background(), "REVOKE write FROM mleow")
	if xerr != nil {
		var pgErr *pgconn.PgError
		if errors.As(xerr, &pgErr) {
			fmt.Println("ERR:", pgErr.Message) // => syntax error at end of input
			fmt.Println("CODE:", pgErr.Code)   // => 42601
		}
		//spew.Dump(xerr.Error())
		t.Fatal(xerr)
	}
	spew.Dump(ct.RowsAffected()) // aalways 0
	spew.Dump(ct.String())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddRole(tt.args.username, tt.args.role); (err != nil) != tt.wantErr {
				t.Errorf("AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
