package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Shriharsh07/chaintrack/config"
	"github.com/Shriharsh07/chaintrack/models"
	"github.com/Shriharsh07/chaintrack/service"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// @Summary      Create a new transaction
// @Description  Creates a new transaction with sender, receiver, and amount
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        transaction  body      models.Transaction  true  "Transaction details"
// @Success      201          {object}  models.Transaction
// @Failure      400          {object}  map[string]interface{}  "Validation failed or invalid input"
// @Failure      500          {string}  string  "Failed to create transaction"
// @Router       /transaction [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validate.Struct(tx); err != nil {
		// If the error is from the validator, extract detailed errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, ve := range validationErrors {
				field := strings.ToLower(ve.Field())
				switch ve.Tag() {
				case "required":
					errors[field] = "This field is required"
				case "email":
					errors[field] = "Invalid email format"
				case "min":
					errors[field] = fmt.Sprintf("Value must be at least %s", ve.Param())
				default:
					errors[field] = fmt.Sprintf("Validation failed on '%s'", ve.Tag())
				}
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":  "Validation failed",
				"fields": errors,
			})
			return
		}
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

// @Summary      Get All Transactions
// @Description  Retrieve all transactions from the database
// @Tags         Transactions
// @Produce      json
// @Success      200  {array}  models.Transaction
// @Failure      500  {string}  string  "Failed to retrieve transactions"
// @Router       /transactions [get]
func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions := []models.Transaction{}
	if err := config.DB.Find(&transactions).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

// @Summary      Get Transactions by Block ID
// @Description  Retrieve all transactions for a specific block by its ID
// @Tags         Transactions
// @Produce      json
// @Param        blockId  path      int  true  "Block ID"
// @Success      200      {array}   models.Transaction
// @Failure      500      {string}  string  "Failed to retrieve transactions"
// @Router       /transactions/{blockId} [get]
func GetTransactionsByBlockID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["blockId"]

	var transactions []models.Transaction
	if err := config.DB.Where("block_id = ?", id).Find(&transactions).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve transactions: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

// @Summary      Validate Blockchain Integrity
// @Description  Validate the blockchain by checking the hashes and proof of work of each block
// @Tags         Blockchain
// @Produce      json
// @Success      200  {object}  map[string]string  "Blockchain is valid"
// @Failure      400  {object}  map[string]interface{}  "Blockchain is invalid"
// @Failure      500  {string}  string  "Failed to fetch blocks"
// @Router       /validate [get]
func ValidateChain(w http.ResponseWriter, r *http.Request) {
	var blocks []models.Block
	if err := config.DB.Order("id asc").Find(&blocks).Error; err != nil {
		http.Error(w, "Failed to fetch blocks", http.StatusInternalServerError)
		return
	}

	difficulty := 4
	var invalidBlocks []map[string]interface{}

	for _, block := range blocks {
		var txs []models.Transaction
		config.DB.Where("block_id = ?", block.ID).Find(&txs)

		data := ""
		for _, tx := range txs {
			data += fmt.Sprintf("%s->%s:%.2f|", tx.Sender, tx.Receiver, tx.Amount)
		}

		expectedHash := service.CalculateHash(data + block.PreviousHash + fmt.Sprintf("%d", block.Nonce))

		if block.Hash != expectedHash || !service.IsValidPoW(block.Hash, difficulty) {
			invalidBlocks = append(invalidBlocks, map[string]interface{}{
				"block_id":      block.ID,
				"expected_hash": expectedHash,
				"stored_hash":   block.Hash,
			})
		}
	}

	if len(invalidBlocks) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":        "Blockchain is invalid",
			"invalid_blocks": invalidBlocks,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blockchain is valid",
	})
}
