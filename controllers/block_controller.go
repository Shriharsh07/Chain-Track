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

// @Summary      Mines a new block
// @Description  Takes all unmined transactions and creates a new block
// @Tags         Block Mining
// @Produce      json
// @Success      201  {object}  models.BlockResponse
// @Failure      400  {string}  string  "No transactions to mine"
// @Failure      500  {string}  string  "Failed to mine block"
// @Router       /mine [post]
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

// @Summary      Get All Blocks
// @Description  Retrieve all blocks from the database
// @Tags         Blocks
// @Produce      json
// @Success      200  {array} []models.BlockResponse
// @Failure      500  {string}  string  "Failed to retrieve blocks"
// @Router       /blocks [get]
func GetBlocks(w http.ResponseWriter, r *http.Request) {
	var blocks []models.Block
	if err := config.DB.Order("id desc").Find(&blocks).Error; err != nil {
		http.Error(w, "Failed to retrieve blocks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(blocks)
}

// @Summary      Get Block by ID
// @Description  Retrieve a single block by its ID
// @Tags         Blocks
// @Produce      json
// @Param        id   path      int  true  "Block ID"
// @Success      200  {object}  models.BlockResponse
// @Failure      404  {string}  string  "Block not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /block/{id} [get]
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

// @Summary      Tamper Block Data
// @Description  Tamper the data of a transaction within the specified block
// @Tags         Blocks
// @Param        id   path      int     true  "Block ID"
// @Produce      json
// @Success      200  {string}  string  "Block data tampered successfully"
// @Failure      404  {string}  string  "Block or transaction not found"
// @Failure      500  {string}  string  "Failed to tamper transaction"
// @Router       /block/{id} [post]
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
