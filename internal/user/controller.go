package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/parmeet20/golang-chatapp/internal/auth"
	"github.com/parmeet20/golang-chatapp/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	Service     *UserService
	authService *auth.AuthService
	validator   *validator.Validate
}

func NewUserController(service *UserService, authService *auth.AuthService) *UserController {
	return &UserController{
		Service:     service,
		authService: authService,
		validator:   validator.New(),
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(user); err != nil {
		response.JSON(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	if c.IsEmailExists(user.Email) {
		response.JSON(w, http.StatusBadRequest, errors.New("email already exists").Error())
		return
	}

	if err := c.Service.Register(&user); err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
	})
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(body); err != nil {
		response.JSON(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return
	}

	token, err := c.Service.Login(body.Username, body.Password)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func (c *UserController) GetMeByToken(w http.ResponseWriter, r *http.Request) {

	claims, err := auth.GetClaims(r)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, "invalid token")
		return
	}

	objUserId, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	user, err := c.Service.FindById(objUserId)
	if err != nil {
		response.JSON(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (c *UserController) IsEmailExists(email string) bool {

	if email == "" {
		return false
	}

	_, err := c.Service.FindByEmail(email)

	return err == nil
}
