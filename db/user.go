package db

import (
	"context"
	"time"
)

type User struct {
	ID             int       `db:"id"`
	FirstName      string    `db:"first_name"`
	LastName       string    `db:"last_name"`
	HashedPassword string    `db:"password"`
	Email          string    `db:"email"`
	CreatedAt      time.Time `db:"created_at"`
}

func (s *Store) CreateUser(ctx context.Context, arg User) (User, error) {
	data := User{}
	stmt, err := s.db.PrepareNamed(`
	INSERT INTO users (
		first_name,
		last_name,
		password,
		email)
	VALUES (
		:first_name, 
		:last_name, 
		:password,
		:email
	) 
	RETURNING *`)
	if err != nil {
		return data, err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			return
		}
	}()
	return data, stmt.Get(&data, arg)
}

func (s *Store) GetUser(ctx context.Context, email string) (User, error) {
	data := make([]User, 0)
	if err := s.db.Select(&data, `
	SELECT * from users
	WHERE email = $1 LIMIT 1`, email); err != nil {
		return User{}, err
	}
	if len(data) > 0 {
		return data[0], nil
	}
	return User{}, nil
}
