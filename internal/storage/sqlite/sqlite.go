package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/iarpitnagpure/go-rest-api/internal/config"
	"github.com/iarpitnagpure/go-rest-api/internal/types"
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

// Get student by using id
// Sqlite struct inherit GetStudentById method
func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no Studnt found")
		}
		return types.Student{}, fmt.Errorf("query Error")
	}

	return student, nil
}

// Get all student
// Sqlite struct inherit GetStudents method
func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

// Update student by using student payload
// Sqlite struct inherit UpdateStudent method
func (s *Sqlite) UpdateStudent(student types.Student) (types.Student, error) {
	stmt, err := s.Db.Prepare(`
		UPDATE students
		SET name = ?, email = ?, age = ?
		WHERE id = ?
	`)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(student.Name, student.Email, student.Age, student.Id)
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to execute update: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, fmt.Errorf("failed to get affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return types.Student{}, fmt.Errorf("no student found with id %d", student.Id)
	}

	return student, nil
}

// Delete student by using student payload
// Sqlite struct inherit DeleteStudentByID method
func (s *Sqlite) DeleteStudentById(id int64) (bool, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("failed to prepare delete statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return false, fmt.Errorf("failed to execute delete: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return false, fmt.Errorf("no student found with id %d", id)
	}

	return true, nil
}
