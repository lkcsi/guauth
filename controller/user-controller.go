package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lkcsi/goauth/custerror"
	"github.com/lkcsi/goauth/entity"
	"github.com/lkcsi/goauth/service"
)

type UserController interface {
	FindByUsername(*gin.Context)
	Save(*gin.Context)
	Login(*gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(s *service.UserService) UserController {
	return &userController{userService: *s}
}

func (u *userController) Save(c *gin.Context) {
	c.Writer.Header().Set("content-type", "application/json")
	var newUser entity.User
	if err := c.BindJSON(&newUser); err != nil {
		setError(c, err)
		return
	}
	if _, err := u.userService.Save(&newUser); err != nil {
		setError(c, err)
		return
	}
	c.IndentedJSON(201, newUser)
}

func (u *userController) FindByUsername(context *gin.Context) {
	username := context.Param("username")
	user, err := u.userService.FindByUsername(username)
	if err != nil {
		setError(context, err)
		return
	}
	context.IndentedJSON(200, user)
}

func (u *userController) Login(c *gin.Context) {
	c.Writer.Header().Set("content-type", "application/json")
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		setError(c, err)
		return
	}
	jwt, err := u.userService.Login(&user)
	if err != nil {
		setError(c, err)
		return
	}
	c.IndentedJSON(200, gin.H{"access_token": jwt})
}

func setError(context *gin.Context, err error) {
	switch err.(type) {
	case validator.ValidationErrors:
		s := err.(validator.ValidationErrors)
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": getErrorMsg(s[0])})
	case custerror.CustError:
		s := err.(custerror.CustError)
		context.AbortWithStatusJSON(s.Code(), gin.H{"error": s.Error()})
	default:
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.ActualTag() {
	case "required":
		return fe.Field() + " field is mandatory"
	case "gte":
		return fe.Field() + " must be greater than or equals " + fe.Param()
	}
	return "unkown error"
}
