package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	sqlStr := `INSERT INTO users (name, email, hashed_password, created) VALUES
			(?, ?, ?, UTC_TIMESTAMP())`

	res, err := m.Db.Exec(sqlStr, name, email, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			fmt.Print(mySQLError.Message)
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return 0, ErrDuplicateEmail
			}
			return 0, err
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	u, err := m.getUserByEmail(email)
	if err != nil {

		return 0, err

	}

	fmt.Printf("Hashed password: %v\n", u.HashedPassword)
	fmt.Printf("Plain text password: %v\n", password)

	if err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password)); err != nil {
		fmt.Printf("Password might be wrong?")
		fmt.Println(err)
		return 0, ErrInvalidCredentials
	}
	return u.Id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (m *UserModel) getUserByEmail(email string) (User, error) {
	sqlStr := `SELECT id, name, email, hashed_password, created FROM users WHERE email = ?`

	row := m.Db.QueryRow(sqlStr, email)

	var u User

	if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.HashedPassword, &u.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserDoesNotExist
		} else {
			return User{}, err
		}
	}

	return u, nil
}
