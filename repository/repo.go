package repository

import (
	"database/sql"
	"errors"
	"privateCabin/entity"

	"github.com/lib/pq"
)

type UserRepository interface {
	GetUserByLogin(login, password string) (*entity.UserPublicDTO, error)
	CreateUserByData(login, password string) (*entity.UserPublicDTO, error)
	ListAllUsers() ([]entity.UserPublicDTO, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByLogin(login, password string) (*entity.UserPublicDTO, error) {
	row := r.db.QueryRow("SELECT id, login FROM users WHERE login = $1 and password = $2", login, password)
	user := &entity.UserPublicDTO{}
	err := row.Scan(&user.ID, &user.Login)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUserByData(login, password string) (*entity.UserPublicDTO, error) {
	var user entity.UserPublicDTO
	query := `INSERT INTO users (login,password) VALUES ($1, $2) RETURNING id, login;`

	err := r.db.QueryRow(query, login, password).Scan(&user.ID, &user.Login)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return nil, errors.New("login_already_exists")
			}

		}
		return nil, err
	}
	return &user, nil

}

func (r *userRepository) ListAllUsers() ([]entity.UserPublicDTO, error) {
	rows, err := r.db.Query("SELECT id, login FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entity.UserPublicDTO
	for rows.Next() {
		var user entity.UserPublicDTO
		if err := rows.Scan(&user.ID, &user.Login); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}
