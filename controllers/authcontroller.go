package controllers

import (
	"encoding/json"
	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var register models.Register

	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}

	defer r.Body.Close()

	if register.Password != register.PasswordConfirm {
		helpers.Response(w, 400, "Password and Password Confirm do not match", nil)

		return
	}

	passwordHash, err := helpers.HashPassword(register.Password)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}
	user := models.User{
		Name:     register.Name,
		Email:    register.Email,
		Role:     "user",
		Password: passwordHash,
	}
	if err := configs.DB.Create(&user).Error; err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}
	userResponse := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	helpers.Response(w, 201, "User created successfully", userResponse)

}

func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var register models.Register

	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}

	defer r.Body.Close()

	if register.Password != register.PasswordConfirm {
		helpers.Response(w, 400, "Password and Password Confirm do not match", nil)

		return
	}

	passwordHash, err := helpers.HashPassword(register.Password)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}
	user := models.User{
		Name:     register.Name,
		Email:    register.Email,
		Role:     "admin",
		Password: passwordHash,
	}
	if err := configs.DB.Create(&user).Error; err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}
	userResponse := models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	helpers.Response(w, 201, "User created successfully", userResponse)

}

func Login(w http.ResponseWriter, r *http.Request) {

	var login models.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}

	defer r.Body.Close()

	var user models.User
	if err := configs.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		helpers.Response(w, 404, "Wrong email or password", nil)

		return
	}

	if err := helpers.VerifyPassword(login.Password, user.Password); err != nil {
		helpers.Response(w, 404, "Wrong email or password", nil)

		return
	}

	token, err := helpers.CreateToken(&user)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)

		return
	}

	userLoginResponse := models.UserLoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	helpers.Response(w, 200, "Login successful", userLoginResponse)
}
