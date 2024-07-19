package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	Db *sql.DB
}

// Insert a new snippet to the database and return its id
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	sql := `INSERT INTO snippets (title, content, created, expires) VALUES 
			(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := m.Db.Exec(sql, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// Return a snippet by id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Return the 10 recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
