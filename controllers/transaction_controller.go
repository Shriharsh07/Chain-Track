package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shriharsh07/chaintrack/config"
	"github.com/Shriharsh07/chaintrack/models"
	"github.com/google/uuid"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()

	if err := config.DB.Create(&tx).Error; err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := []models.Transaction{}
	if err := config.DB.Find(&transactions).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}
