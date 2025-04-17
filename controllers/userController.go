package controllers

import (
	"net/http"
	"strconv"

	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"

	"github.com/gorilla/mux"
)

func Me(w http.ResponseWriter, r *http.Request) {
	// Ambil data dari JWT
	userClaims := r.Context().Value("userinfo").(*helpers.MyCustomClaims)

	// Ambil data user lengkap dari database
	var user models.User
	if err := configs.DB.First(&user, userClaims.ID).Error; err != nil {
		helpers.Response(w, 404, "User not found", nil)
		return
	}

	// Buat response
	userResponse := &models.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		PhoneNumber: user.PhoneNumber,
		ImageUrl: user.ImageUrl,
	}

	helpers.Response(w, 200, "My Profile", userResponse)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL dan convert ke integer
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid user id", nil)
		return
	}

	// 2. Cari user dari database
	var user models.User
	if err := configs.DB.First(&user, idFromParam).Error; err != nil {
		helpers.Response(w, 404, "User not found", nil)
		return
	}

	// 3. Buat response tanpa password
	userResponse := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
	}

	// 4. Kirim response
	helpers.Response(w, 200, "User by id", userResponse)
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	var user []models.User
	if err := configs.DB.Find(&user).Error; err != nil {
		helpers.Response(w, 500, "Failed to get all user", nil)
		return
	}

	var userResponses []models.UserResponse
	for _, u := range user {
		userResponses = append(userResponses, models.UserResponse{
			ID:          u.ID,
			Name:        u.Name,
			Email:       u.Email,
			PhoneNumber: u.PhoneNumber,
			ImageUrl:    u.ImageUrl,
		})
	}

	helpers.Response(w, 200, "Successfully", userResponses)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid user id", nil)
		return
	}

	// Cari user berdasarkan ID
	var user models.User
	if err := configs.DB.First(&user, idFromParam).Error; err != nil {
		helpers.Response(w, 404, "User not found", nil)
		return
	}

	// Parse multipart form (10MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.Response(w, 400, "Could not parse multipart form", nil)
		return
	}

	// Ambil field biasa
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")

	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.PhoneNumber = phone
	}

	// Handle image (opsional)
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()

		// Upload ke ImageKit
		imageUrl, err := helpers.UploadToImageKit(file, handler.Filename)
		if err != nil {
			helpers.Response(w, 500, err.Error(), nil)
			return
		}
		user.ImageUrl = imageUrl // Pastikan field ini ada di model `User`
	}

	// Simpan perubahan
	if err := configs.DB.Save(&user).Error; err != nil {
		helpers.Response(w, 500, "Failed to update user", nil)
		return
	}

	// Response
	userResponse := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ImageUrl:    user.ImageUrl,
	}

	helpers.Response(w, 200, "Success update user", userResponse)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid user id", nil)
		return
	}

	var user models.User
	if err := configs.DB.First(&user, idFromParam).Error; err != nil {
		helpers.Response(w, 404, "User not found", nil)
		return
	}

	if err := configs.DB.Delete(&user, idFromParam).Error; err != nil {
		helpers.Response(w, 500, "Error deleting user", nil)
		return
	}

	helpers.Response(w, 200, "Sccess delete user", nil)
}
