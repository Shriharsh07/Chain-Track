package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shriharsh07/chaintrack/config"
	"github.com/Shriharsh07/chaintrack/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func markTransactionAsMined(transactionID uuid.UUID) {
	var transaction models.Transaction
	if err := config.DB.First(&transaction, transactionID).Error; err != nil {
		log.Println("Transaction not found")
		return
	}
	transaction.IsMined = true
	config.DB.Save(&transaction)
}

func MineBlock(w http.ResponseWriter, r *http.Request) {
	var transactions []models.Transaction
	result := config.DB.Model(&models.Transaction{}).Where("is_mined = ?", false).Find(&transactions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	if len(transactions) == 0 {
		http.Error(w, "No transactions to mine", http.StatusBadRequest)
		return
	}

	var lastHash string

	var latestBlock models.Block
	if err := config.DB.Order("id desc").First(&latestBlock).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			lastHash = "" // No blocks in the database, start from an empty hash
		} else {
			http.Error(w, "Failed to fetch the latest block", http.StatusInternalServerError)
			return
		}
	} else {
		lastHash = latestBlock.Hash
	}

	data := ""
	for _, tx := range transactions {
		data += fmt.Sprintf("%s->%s:%.2f|", tx.Sender, tx.Receiver, tx.Amount)
	}

	hash := sha256.Sum256([]byte(data + lastHash + time.Now().String()))
	hashStr := hex.EncodeToString(hash[:])

	block := models.Block{
		Transactions: transactions,
		PreviousHash: lastHash,
		Hash:         hashStr,
	}

	if err := config.DB.Create(&block).Error; err != nil {
		http.Error(w, "Failed to mine block", http.StatusInternalServerError)
		return
	}

	for _, tx := range transactions {
		markTransactionAsMined(tx.ID)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(block)
}

func GetBlocks(w http.ResponseWriter, r *http.Request) {
	var blocks []models.Block
	if err := config.DB.Order("id desc").Find(&blocks).Error; err != nil {
		http.Error(w, "Failed to retrieve blocks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(blocks)
}

func GetBlockByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var block models.Block
	if err := config.DB.First(&block, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("Block not found: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(block)
}
