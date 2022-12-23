package salesforcedb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	endpoint, username, password, token string
)

func getEnvVar(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s environment variable not set", key)
	}

	return value
}

func TestMain(m *testing.M) {
	// show file & location, date & time
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// get connection parameters from environment
	username = getEnvVar("USERNAME")
	password = getEnvVar("PASSWORD")
	token = getEnvVar("TOKEN")
	endpoint = getEnvVar("ENDPOINT")

	if username == "" || password == "" || token == "" || endpoint == "" {
		os.Exit(1)
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

// test that database/sql works
func TestGetResultsSql(t *testing.T) {
	// limit number of rows returned for testing
	// should be > 2000 since that the batch size for simpleforce
	rowLimit := 2100

	// connect
	sqlCon, err := sql.Open("salesforcedb", fmt.Sprintf("%s/%s/%s@%s", username, password, token, endpoint))
	if err != nil {
		t.Fatalf("sql.Open error: %+v", err)
	}

	// query
	rows, err := sqlCon.Query(fmt.Sprintf("select id, name from account limit %d", rowLimit))
	if err != nil {
		t.Fatalf("conn.Query error: %+v", err)
	}
	defer rows.Close()

	// get columns
	cols, err := rows.Columns()
	if err != nil {
		t.Fatalf("rows.Columns error: %+v", err)
	}
	if len(cols) == 0 {
		t.Fatalf("no columns returned")
	}

	// go thru results to validate getting multiple batches
	vals := make([]interface{}, len(cols))
	for i := range vals {
		vals[i] = new(interface{})
	}
	numRows := 0
	for rows.Next() {
		err = rows.Scan(vals...)
		if err != nil {
			t.Fatalf("rows.Scan error: %+v", err)
		}
		numRows++
	}
	if numRows != rowLimit {
		t.Fatalf("didn't retrieve all results")
	}

	err = rows.Close()
	if err != nil {
		t.Fatalf("rows.Close error: %+v", err)
	}
}

// test that jmoiron/sqlx works
func TestGetResultsSqlx(t *testing.T) {
	// limit number of rows returned for testing
	// should be > 2000 since that the batch size for simpleforce
	rowLimit := 2100

	// connect
	sqlxCon, err := sqlx.Connect("salesforcedb", fmt.Sprintf("%s/%s/%s@%s", username, password, token, endpoint))
	if err != nil {
		t.Fatalf("sqlx.Open error: %+v", err)
	}

	// get all the results
	accounts := []struct {
		ID   string `db:"Id"`
		Name string `db:"Name"`
	}{}
	err = sqlxCon.Select(&accounts, fmt.Sprintf("select id, name from account limit %d", rowLimit))
	if err != nil {
		t.Fatalf("sqlx.Select error: %+v", err)
	}

	if len(accounts) != rowLimit {
		t.Fatalf("didn't retrieve all results")
	}
}

func TestSupported(t *testing.T) {
	// connect
	sqlCon, err := sql.Open("salesforcedb", fmt.Sprintf("%s/%s/%s@%s", username, password, token, endpoint))
	if err != nil {
		t.Fatalf("sql.Open error: %+v", err)
	}

	_, err = sqlCon.Begin()
	if err != errNoTransactions {
		t.Fatalf("conn.Begin wrong error: %+v", err)
	}

	_, err = sqlCon.Prepare("query")
	if err != errNoPrepared {
		t.Fatalf("conn.Prepare wrong error: %+v", err)
	}

	err = sqlCon.Ping()
	if err != nil {
		t.Fatalf("conn.Ping wrong error: %+v", err)
	}

	err = sqlCon.Close()
	if err != nil {
		t.Fatalf("conn.Close wrong error: %+v", err)
	}
}
