package service

import (
	"database/sql"
	"errors"
	"privateCabin/entity"
	"privateCabin/repository"
)

type UserService interface {
	GetUser(login, password string) (*entity.UserPublicDTO, error)
	CreateUser(login, password string) (*entity.UserPublicDTO, error)
	ListUsers() (*entity.ResponseUserList, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetUser(login, password string) (*entity.UserPublicDTO, error) {
	if login == "" || password == "" {
		return nil, errors.New("login and password required")
	}

	user, err := s.repo.GetUserByLogin(login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid login or password") // ← бизнес-ошибка (не технический ErrNoRows)
		}
		return nil, err // ← техническая ошибка БД
	}
	return user, nil
}

func (s *userService) CreateUser(login, password string) (*entity.UserPublicDTO, error) {
	if login == "" || password == "" {
		return nil, errors.New("login and password required")
	}
	return s.repo.CreateUserByData(login, password)
}

func (s *userService) ListUsers() (*entity.ResponseUserList, error) {
	users, err := s.repo.ListAllUsers()
	if err != nil {
		return nil, err
	}
	return &entity.ResponseUserList{Users: users}, nil
}
