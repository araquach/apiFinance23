package routes

import (
	"github.com/araquach/apiFinance23/handlers"
)

func profitLoss() {
	s := R.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/costs-and-takings/{start}/{end}", handlers.ApiCostsAndTakings).Methods("GET")
}
