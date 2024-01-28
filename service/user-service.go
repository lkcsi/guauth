package service

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lkcsi/goauth/custerror"
	"github.com/lkcsi/goauth/entity"
	"github.com/lkcsi/goauth/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByUsername(string) (*entity.User, error)
	Save(*entity.User) (*entity.User, error)
	Login(*entity.User) (string, error)
}

type inMemoryUserService struct {
	userRepository repository.UserRepository
}

func (service *inMemoryUserService) FindByUsername(username string) (*entity.User, error) {
	return service.userRepository.FindByUsername(username)
}

func (service *inMemoryUserService) Login(requser *entity.User) (string, error) {
	user, err := service.userRepository.FindByUsername(requser.Username)
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requser.Password)) != nil {
		return "", custerror.InvalidPasswordError()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	secret := os.Getenv("AUTH_SECRET")
	jwtToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (service *inMemoryUserService) Save(userRequest *entity.User) (*entity.User, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := entity.User{Username: userRequest.Username, Password: string(pwd)}
	err = service.userRepository.Save(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &inMemoryUserService{userRepository: *repo}
}
