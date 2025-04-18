package db

import (
	"fmt"


	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
	"gitlab.com/melsonmascarenhas/backend/error"
	"gitlab.com/melsonmascarenhas/backend/models"
)
var(
	account models.Account
	mu sync.Mutex
)


// Connect to SQLite database
func ConnectDB() (*sqlx.DB, bool) {

	fmt.Println("connectdb")
	var ok bool
	db, err := sqlx.Open("sqlite", "./test.db")
	error.IfError(err)

	err = db.Ping()

	if error.IfError(err) && db.DriverName() != "sqlite" {

		Create()
	}

	// error.IfError(err)
	ok = true

	return db, ok
}
// Execute a given SQL statement
func dbExec(sqlStatement string) {
	fmt.Println("dbexec")

	db, ok := ConnectDB()
	defer db.Close()

	if ok {
		mu.Lock()
		_, err := db.Exec(sqlStatement)
		mu.Unlock()
		error.IfError(err)
	}
}
func Create() {
	sqlStatement := `CREATE TABLE IF NOT EXISTS "account" (
			"ID"         TEXT PRIMARY KEY,
			"CustomerID" TEXT UNIQUE,
			"Name"       TEXT,
			"Balance"    REAL,
			"CreatedAt"  TEXT
		);`
	
	dbExec(sqlStatement)
}
func Insert(account models.Account) {
	fmt.Println("Inserting into database...")

	sqlStatement := `INSERT INTO account ("ID", "CustomerID", "Name", "Balance", "CreatedAt") 
	VALUES ('%s', '%s', '%s', '%f', '%s');`

	sqlStatement = fmt.Sprintf(
		sqlStatement,
		account.ID,
		account.CustomerID,
		account.Name,
		account.Balance,
		account.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	dbExec(sqlStatement)
	fmt.Println("Insertion completed.")
}
