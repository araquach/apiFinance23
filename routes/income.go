package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func Income(r *mux.Router) {
	r.HandleFunc("/income-by-month/{start}/{end}", handlers.ApiMonthlyIncomeByDateRange).Methods("GET")
}
