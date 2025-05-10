package routes

import (
	"github.com/Shriharsh07/chaintrack/controllers"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/transaction", controllers.CreateTransaction).Methods("POST")
	r.HandleFunc("/transactions", controllers.GetAllTransactions).Methods("GET")
	r.HandleFunc("/transactions/{blockId}", controllers.GetTransactionsByBlockID).Methods("GET")
	r.HandleFunc("/mine", controllers.MineBlock).Methods("POST")
	r.HandleFunc("/blocks", controllers.GetBlocks).Methods("GET")
	r.HandleFunc("/block/{id}", controllers.GetBlockByID).Methods("GET")
	r.HandleFunc("/block/{id}", controllers.TamperBlockData).Methods("POST")
	r.HandleFunc("/validate", controllers.ValidateChain).Methods("GET")
}
