package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	Db *sql.DB
}

func (m *UserModel) Insert(name string, email string, password string) (int, error) {
	return 0, nil
	// sqlStr := `INSERT INTO users (name, email, hashed_password, created) VALUES
	// 		(?, ?, ?, UTC_TIMESTAMP())`

	// res, err := m.Db.Exec(sqlStr, name, email, hashedPassword)
	// if err != nil {
	// 	if errors.Is(err, m.db)
	// 	return 0, err
	// }

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
