package routes

import (
	"github.com/araquach/apiFinance23/handlers"
)

func takings() {
	s := R.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/takings-by-year", handlers.ApiTakingsByYear).Methods("GET")
	s.HandleFunc("/monthly-takings-by-date-range/{start}/{end}", handlers.ApiMonthlyTakingsByDateRange).Methods("GET")
	s.HandleFunc("/stylist-takings-month-by-month/{fName}/{lName}/{start}/{end}", handlers.ApiStylistTakingsMonthByMonth).Methods("GET")
	s.HandleFunc("/takings-by-stylist-comparison/{salon}/{start}/{end}", handlers.ApiTakingsByStylistComparison).Methods("GET")
	s.HandleFunc("/totals-by-date-range/{start}/{end}", handlers.ApiTotalsByDateRange).Methods("GET")
}
