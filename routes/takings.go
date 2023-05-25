package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func Takings(r *mux.Router) {
	r.HandleFunc("/takings-by-year", handlers.ApiTakingsByYear).Methods("GET")
	r.HandleFunc("/monthly-takings-by-date-range/{start}/{end}", handlers.ApiMonthlyTakingsByDateRange).Methods("GET")
	r.HandleFunc("/stylist-takings-month-by-month/{fName}/{lName}/{start}/{end}", handlers.ApiStylistTakingsMonthByMonth).Methods("GET")
	r.HandleFunc("/takings-by-stylist-comparison/{salon}/{start}/{end}", handlers.ApiTakingsByStylistComparison).Methods("GET")
	r.HandleFunc("/totals-by-date-range/{start}/{end}", handlers.ApiTotalsByDateRange).Methods("GET")
}
