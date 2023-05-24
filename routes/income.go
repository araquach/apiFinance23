package routes

import (
	"github.com/araquach/apiFinance23/handlers"
)

func income() {
	s := R.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/income-by-month/{start}/{end}", handlers.ApiMonthlyIncomeByDateRange).Methods("GET")
}
