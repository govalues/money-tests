package sqlite_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/govalues/money"
	_ "modernc.org/sqlite"
)

const (
	url            = ":memory:"
	selectNull     = "SELECT null"
	dropTable      = "DROP TABLE IF EXISTS curr_tests"
	createTable    = "CREATE TABLE curr_tests (curr TEXT) STRICT"
	insertCurrency = "INSERT INTO curr_tests (curr) VALUES ($1) RETURNING curr"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("sqlite", url)
	if err != nil {
		log.Fatalf("Open(%q) failed: %v\n", url, err)
	}
	defer db.Close()
	_, err = db.Exec(dropTable)
	if err != nil {
		log.Fatalf("Exec(%q) failed: %v\n", dropTable, err)
	}
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Exec(%q) failed: %v\n", createTable, err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestCurrency_selectNull(t *testing.T) {
	t.Run("money.Currency", func(t *testing.T) {
		row := db.QueryRow(selectNull)
		var got money.Currency
		err := row.Scan(&got)
		if err == nil {
			t.Errorf("QueryRow(%q) did not fail", selectNull)
			return
		}
	})

	t.Run("money.NullCurrency", func(t *testing.T) {
		row := db.QueryRow(selectNull)
		var got money.NullCurrency
		err := row.Scan(&got)
		if err != nil {
			t.Errorf("QueryRow(%q) failed: %v", selectNull, err)
			return
		}
		var want money.NullCurrency
		if got != want {
			t.Errorf("Scan() = %v, want %v", got, want)
		}
	})
}

func TestCurrency_insert(t *testing.T) {
	tests := []money.Currency{
		money.USD, money.JPY, money.OMR,
	}

	t.Run("money.Currency", func(t *testing.T) {
		for _, tt := range tests {
			row := db.QueryRow(insertCurrency, tt)
			var got money.Currency
			err := row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertCurrency, tt, err)
				continue
			}
			if got != tt {
				t.Errorf("Scan(&got) = %v, want %v", got, tt)
				continue
			}
		}
	})

	t.Run("money.NullCurrency", func(t *testing.T) {
		for _, tt := range tests {
			row := db.QueryRow(insertCurrency, tt)
			var got money.NullCurrency
			err := row.Scan(&got)
			if err != nil {
				t.Errorf("QueryRow(%q, %v) failed: %v", insertCurrency, tt, err)
				continue
			}
			want := money.NullCurrency{Currency: tt, Valid: true}
			if got != want {
				t.Errorf("Scan(&got) = %v, want %v", got, want)
				continue
			}
		}
	})
}
