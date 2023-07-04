package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func ProfitLoss(r *mux.Router) {
	s := r.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/costs-and-takings/{start}/{end}", handlers.ApiCostsAndTakings).Methods("GET")
}
