package controllers

import (
	"fmt"
	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		helpers.Response(w, 400, "Could not parse multipart form", nil)
		return
	}

	// Ambil field biasa
	name := r.FormValue("name")
	description := r.FormValue("description")
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

	// Ambil file
	file, handler, err := r.FormFile("image")
	if err != nil {
		helpers.Response(w, 400, "Image file is required", nil)
		return
	}
	defer file.Close()

	// Cari kategori
	var category models.Category
	if err := configs.DB.First(&category, categoryID).Error; err != nil {
		message := fmt.Sprintf("Category with ID %d not found", categoryID)
		helpers.Response(w, 404, message, nil)
		return
	}

	// Upload ke ImageKit
	imageUrl, err := helpers.UploadToImageKit(file, handler.Filename)
	if err != nil {
		helpers.Response(w, 500, err.Error(), nil)
		return
	}

	// Simpan product
	product := models.Product{
		Name:        name,
		Description: description,
		Stock:       uint(stock),
		Price:       uint(price),
		UnitsSold:   0,
		CategoryID:  uint(categoryID),
		ImageUrl:    imageUrl,
	}
	if err := configs.DB.Save(&product).Error; err != nil {
		helpers.Response(w, 400, err.Error(), nil)
		return
	}

	productResponse := models.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Stock:        product.Stock,
		Price:        product.Price,
		UnitsSold:    product.UnitsSold,
		ImageUrl:     product.ImageUrl,
		CategoryID:   product.CategoryID,
		CategoryName: category.Name,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
	helpers.Response(w, 201, "Product created successfully", productResponse)
}

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter page dan limit, default: page=1, limit=10
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	var products []models.Product
	if err := configs.DB.Preload("Category").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		helpers.Response(w, 500, "Failed to get products", nil)
		return
	}

	var productResponses []models.ProductResponse
	for _, p := range products {
		productResponse := models.ProductResponse{
			ID:           p.ID,
			Name:         p.Name,
			Description:  p.Description,
			Stock:        p.Stock,
			Price:        p.Price,
			UnitsSold:    p.UnitsSold,
			ImageUrl:     p.ImageUrl,
			CategoryID:   p.CategoryID,
			CategoryName: p.Category.Name,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		}
		productResponses = append(productResponses, productResponse)
	}

	helpers.Response(w, 200, "Success get all products", productResponses)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid product id", nil)
		return
	}

	var product models.Product
	if err := configs.DB.Preload("Category").First(&product, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Product with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	productResponse := models.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Stock:        product.Stock,
		Price:        product.Price,
		UnitsSold:    product.UnitsSold,
		ImageUrl:     product.ImageUrl,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.Name,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	helpers.Response(w, 200, "Success get product", productResponse)
}

func UpdateProductById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid product id", nil)
		return
	}

	// Ambil data produk dari database
	var product models.Product
	if err := configs.DB.Preload("Category").First(&product, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Product with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helpers.Response(w, 400, "Could not parse multipart form", nil)
		return
	}

	// Ambil data dari form
	name := r.FormValue("name")
	description := r.FormValue("description")
	stockStr := r.FormValue("stock")
	priceStr := r.FormValue("price")
	categoryIDStr := r.FormValue("category_id")

	// Update jika ada isian
	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}
	if stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err == nil {
			product.Stock = uint(stock)
		}
	}
	if priceStr != "" {
		if price, err := strconv.Atoi(priceStr); err == nil {
			product.Price = uint(price)
		}
	}
	if categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			// Cek apakah category valid
			var category models.Category
			if err := configs.DB.First(&category, categoryID).Error; err != nil {
				message := fmt.Sprintf("Category with ID %d not found", categoryID)
				helpers.Response(w, 404, message, nil)
				return
			}
			product.CategoryID = uint(categoryID)
			product.Category = category
		}
	}

	// Cek apakah ada file gambar
	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		imageUrl, err := helpers.UploadToImageKit(file, handler.Filename)
		if err != nil {
			helpers.Response(w, 500, err.Error(), nil)
			return
		}
		product.ImageUrl = imageUrl
	}

	// Simpan perubahan ke database
	if err := configs.DB.Save(&product).Error; err != nil {
		helpers.Response(w, 500, "Failed to update product", nil)
		return
	}

	// Buat response
	response := models.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Stock:        product.Stock,
		Price:        product.Price,
		UnitsSold:    product.UnitsSold,
		ImageUrl:     product.ImageUrl,
		CategoryID:   product.CategoryID,
		CategoryName: product.Category.Name,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	helpers.Response(w, 200, "Success update product", response)
}

func DeleteProductById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid product id", nil)
		return
	}

	var product models.Product
	if err := configs.DB.Preload("Category").First(&product, idFromParam).Error; err != nil {
		message := fmt.Sprintf("Product with ID %d not found", idFromParam)
		helpers.Response(w, 404, message, nil)
		return
	}

	if err := configs.DB.Delete(&product, idFromParam).Error; err != nil {
		helpers.Response(w, 400, "Error deleting product", nil)
		return
	}
	helpers.Response(w, 200, "Success delete product", nil)
}
