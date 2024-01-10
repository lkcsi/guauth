package service

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/lkcsi/goauth/custerror"
	"github.com/lkcsi/goauth/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindByUsername(string) (*entity.User, error)
	Save(*entity.User) (*entity.User, error)
	Login(*entity.User) (string, error)
}

type inMemoryUserService struct {
	users map[string]entity.User
}

func (u *inMemoryUserService) FindByUsername(username string) (*entity.User, error) {
	user, ok := u.users[username]
	if !ok {
		return nil, custerror.NotFoundError(username)
	}
	return &user, nil
}

func (u *inMemoryUserService) Login(requser *entity.User) (string, error) {
	user, ok := u.users[requser.Username]
	if !ok {
		return "", custerror.NotFoundError(requser.Username)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requser.Password)) != nil {
		return "", custerror.InvalidPasswordError()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
	})

	secret := os.Getenv("AUTH_SECRET")
	jwtToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (u *inMemoryUserService) Save(userRequest *entity.User) (*entity.User, error) {
	_, exist := u.users[userRequest.Username]
	if exist {
		return nil, custerror.OccupiedFoundError(userRequest.Username)
	}
	pwd, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := entity.User{Username: userRequest.Username, Password: string(pwd)}
	u.users[user.Username] = user

	return &user, nil
}

func NewInMemoryUserService() UserService {
	u := make(map[string]entity.User)
	u["roland"] = entity.User{Username: "roland", Password: "password"}

	return &inMemoryUserService{users: u}
}
