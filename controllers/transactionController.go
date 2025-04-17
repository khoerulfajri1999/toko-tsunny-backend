package controllers

import (
	"encoding/json"
	"fmt"
	"go-jwt/configs"
	"go-jwt/helpers"
	"go-jwt/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Response(w, 400, "Invalid request body", nil)
		return
	}
	defer r.Body.Close()

	// 🔁 Parse TransactionAt (string to time.Time)
	parsedDate, err := time.Parse("2006-01-02", req.TransactionAt)
	if err != nil {
		helpers.Response(w, 400, "Invalid date format. Use YYYY-MM-DD (e.g. 2025-04-12)", nil)
		return
	}

	tx := configs.DB.Begin()
	if tx.Error != nil {
		helpers.Response(w, 500, "Failed to begin transaction", nil)
		return
	}

	var totalAmount uint = 0
	var details []models.TransactionDetail
	var detailResponses []models.TransactionDetailResponse

	for _, d := range req.Details {
		var product models.Product
		if err := tx.First(&product, d.ProductID).Error; err != nil {
			tx.Rollback()
			message := fmt.Sprintf("Product with ID %d not found", d.ProductID)
			helpers.Response(w, 404, message, nil)
			return
		}

		if product.Stock < d.Quantity {
			tx.Rollback()
			message := fmt.Sprintf("Stock not enough for product ID %d", d.ProductID)
			helpers.Response(w, 400, message, nil)
			return
		}

		subTotal := d.Quantity * product.Price
		totalAmount += subTotal

		detail := models.TransactionDetail{
			ProductID: d.ProductID,
			Quantity:  d.Quantity,
			SubTotal:  subTotal,
		}
		details = append(details, detail)

		detailResponse := models.TransactionDetailResponse{
			ProductID: d.ProductID,
			Quantity:  d.Quantity,
			SubTotal:  subTotal,
		}
		detailResponses = append(detailResponses, detailResponse)

		product.Stock -= d.Quantity
		product.UnitsSold += d.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			helpers.Response(w, 500, "Failed to update product", nil)
			return
		}
	}

	transaction := models.Transaction{
		Income:             req.Income,
		Expense:            req.Expense,
		TotalAmount:        totalAmount,
		TransactionAt:      parsedDate, 
		TransactionDetails: details,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		helpers.Response(w, 500, "Failed to create transaction", nil)
		return
	}

	tx.Commit()

	response := models.TransactionResponse{
		ID:            transaction.ID,
		TotalAmount:   transaction.TotalAmount,
		Income:        transaction.Income,
		Expense:       transaction.Expense,
		TransactionAt: transaction.TransactionAt,
		Details:       detailResponses,
	}

	helpers.Response(w, 201, "Transaction created successfully", response)
}


func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	// Ambil query params
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Default pagination
	page := 1
	limit := 10

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	query := configs.DB.Preload("TransactionDetails")

	// Filter tanggal jika ada
	if startDateStr != "" && endDateStr != "" {
		startDate, errStart := time.Parse("2006-01-02", startDateStr)
		endDate, errEnd := time.Parse("2006-01-02", endDateStr)

		if errStart == nil && errEnd == nil {
			query = query.Where("transaction_at BETWEEN ? AND ?", startDate, endDate)
		}
	}

	var transactions []models.Transaction
	if err := query.Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		helpers.Response(w, 500, "Failed to fetch transactions", nil)
		return
	}

	// Bangun response
	var transactionResponses []models.TransactionResponse
	for _, t := range transactions {
		var detailResponses []models.TransactionDetailResponse
		for _, d := range t.TransactionDetails {
			detailResponses = append(detailResponses, models.TransactionDetailResponse{
				ProductID: d.ProductID,
				Quantity:  d.Quantity,
				SubTotal:  d.SubTotal,
			})
		}

		transactionResponses = append(transactionResponses, models.TransactionResponse{
			ID:            t.ID,
			TotalAmount:   t.TotalAmount,
			Income:        t.Income,
			Expense:       t.Expense,
			TransactionAt: t.TransactionAt,
			Details:       detailResponses,
		})
	}

	helpers.Response(w, 200, "Success get all transactions", transactionResponses)
}

func GetTransactionById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idFromParam, err := strconv.Atoi(id)
	if err != nil {
		helpers.Response(w, 400, "Invalid transaction id", nil)
		return
	}

	var transaction models.Transaction
	if err := configs.DB.Preload("TransactionDetails").First(&transaction, idFromParam).Error; err != nil {
		helpers.Response(w, 404, "Transaction not found", nil)
		return
	}

	var detailResponses []models.TransactionDetailResponse
	for _, d := range transaction.TransactionDetails {
		detailResponses = append(detailResponses, models.TransactionDetailResponse{
			ProductID: d.ProductID,
			Quantity:  d.Quantity,
			SubTotal:  d.SubTotal,
		})
	}

	response := models.TransactionResponse{
		ID:            transaction.ID,
		TotalAmount:   transaction.TotalAmount,
		Income:        transaction.Income,
		Expense:       transaction.Expense,
		TransactionAt: transaction.TransactionAt,
		Details:       detailResponses,
	}

	helpers.Response(w, 200, "Success get transaction by ID", response)
}

