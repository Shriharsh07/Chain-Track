package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Shriharsh07/chaintrack/config"
	"github.com/Shriharsh07/chaintrack/models"
	"github.com/Shriharsh07/chaintrack/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func markTransactionAsMined(transactionID uuid.UUID, blockID uint) {
	var transaction models.Transaction
	if err := config.DB.First(&transaction, transactionID).Error; err != nil {
		log.Println("Transaction not found")
		return
	}
	data := config.DB.Model(&models.Transaction{}).Where("id = ?", transactionID).Updates(map[string]interface{}{"is_mined": true, "block_id": blockID})
	if data.Error != nil {
		log.Println("Failed to mark transaction as mined", data.Error.Error())
		return
	}
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

	// Proof-of-Work: Finding the valid nonce
	difficulty := 4 // Number of leading zeros required in the hash
	var nonce int
	var hash string
	var validHash bool

	// Start with an empty hash and loop to find the correct one
	for {
		hash = service.CalculateHash(data + lastHash + fmt.Sprintf("%d", nonce))
		validHash = service.IsValidPoW(hash, difficulty)

		if validHash {
			break
		}

		nonce++
	}

	block := models.Block{
		PreviousHash: lastHash,
		Hash:         hash,
		Nonce:        nonce,
	}

	if err := config.DB.Create(&block).Error; err != nil {
		http.Error(w, "Failed to mine block", http.StatusInternalServerError)
		return
	}

	for _, tx := range transactions {
		markTransactionAsMined(tx.ID, block.ID)
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

func TamperBlockData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockID := vars["id"]

	var block models.Block
	if err := config.DB.First(&block, blockID).Error; err != nil {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	var tx models.Transaction
	if err := config.DB.Where("block_id = ?", block.ID).First(&tx).Error; err != nil {
		http.Error(w, "No transactions found for block", http.StatusNotFound)
		return
	}

	// Tamper the transaction
	tx.Amount += 1000 // e.g., inflate amount
	if err := config.DB.Save(&tx).Error; err != nil {
		http.Error(w, "Failed to tamper transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Block data tampered successfully")
}
