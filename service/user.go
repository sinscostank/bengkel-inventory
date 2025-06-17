package service

import (
	"errors"

	"time"
	"github.com/sinscostank/bengkel-inventory/forms"
	"github.com/sinscostank/bengkel-inventory/models"
	"github.com/sinscostank/bengkel-inventory/repository"
	"github.com/sinscostank/bengkel-inventory/utils"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (us *UserService) Register(req *forms.RegisterForm) error {
	existingUser, _ := us.UserRepo.FindUserByEmail(req.Email)
	if existingUser != nil {
		return errors.New("email already in use")
	}

	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      "karyawan",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return us.UserRepo.CreateUser(&user)
}

func (us *UserService) Login(req *forms.LoginForm) (string, *models.User, error) {
	user, err := us.UserRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", nil, err
	}

	if user == nil || !utils.CheckHash(req.Password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
