package sqlite

import (
	"database/sql"

	"github.com/iarpitnagpure/go-rest-api/internal/config"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

// Database connection
func New(cfg *config.Config) (*Sqlite, error) {
	// Open new database connection using sql package, requires storage path of db file
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// Creating table in DB if doesnt exist already
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	// returing database struct
	return &Sqlite{Db: db}, nil
}

// Adding student on Database
// Sqlite struct inherit CreateStudent method
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
