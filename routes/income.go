package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func Income(r *mux.Router) {
	s := r.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/income-by-month/{start}/{end}", handlers.ApiMonthlyIncomeByDateRange).Methods("GET")
}
