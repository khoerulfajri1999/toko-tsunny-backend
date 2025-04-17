package controllers

import (
	"encoding/json"
	"fmt"
	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		helpers.Response(w, 400, "Error creating category", nil)
		return
	}

	defer r.Body.Close()

	if err := configs.DB.Save(&category).Error; err != nil {
		helpers.Response(w, 500, "Error creating category", nil)
		return
	}

	var categoryResponse models.CategoryResponse
	categoryResponse.ID = category.ID
	categoryResponse.Name = category.Name
	helpers.Response(w, 201, "Success creating category", categoryResponse)
}

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	var category []models.Category

	if err := configs.DB.Find(&category).Error; err != nil {
		helpers.Response(w, 500, "Failed to get all category", nil)
		return
	}

	var categoryResponse []models.CategoryResponse
	for _, u := range category {
		categoryResponse = append(categoryResponse, models.CategoryResponse{
			ID:   u.ID,
			Name: u.Name,
		})
	}
	helpers.Response(w, 200, "Get All Category Success", categoryResponse)
}

func GetCategoryById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid category id", nil)
		return
	}

	var category models.Category
	if err := configs.DB.First(&category, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Category with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	var categoryResponse models.CategoryResponse
	categoryResponse.ID = category.ID
	categoryResponse.Name = category.Name
	helpers.Response(w, 200, "Success get category", categoryResponse)

}

func UpdateCategoryById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid user id", nil)
		return
	}

	var category models.Category
	if err := configs.DB.First(&category, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Category with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	var updateCategory models.Category

	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
		helpers.Response(w, 500, "Failed to update category", nil)
		return
	}

	defer r.Body.Close()

	if updateCategory.Name != "" {
		category.Name = updateCategory.Name
	}

	if err := configs.DB.Save(&category).Error; err != nil {
		helpers.Response(w, 500, "Error update category", nil)
		return
	}

	var categoryResponse models.CategoryResponse
	categoryResponse.ID = category.ID
	categoryResponse.Name = category.Name
	helpers.Response(w, 200, "Success update category", categoryResponse)

}

func DeleteCategoryById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 500, "Failed delete category", nil)
		return
	}

	var category models.Category
	if err := configs.DB.First(&category, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Category with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	if err := configs.DB.Delete(&category, idFromParam).Error; err != nil {
		helpers.Response(w, 500, "Error delete category", nil)
		return
	}

	helpers.Response(w, 200, "Success delete category", nil)
}
