package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func ProfitLoss(r *mux.Router) {
	r.HandleFunc("/costs-and-takings/{start}/{end}", handlers.ApiCostsAndTakings).Methods("GET")
}
