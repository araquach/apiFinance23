package routes

import (
	"github.com/araquach/apiFinance23/handlers"
	"github.com/gorilla/mux"
)

func Costs(r *mux.Router) {
	s := r.PathPrefix("/api/finance").Subrouter()

	s.HandleFunc("/costs-year-by-year", handlers.ApiCostsYearByYear).Methods("GET")
	s.HandleFunc("/costs-month-by-month/{start}/{end}", handlers.ApiCostsMonthByMonth).Methods("GET")
	s.HandleFunc("/costs-by-cat/{salon}/{start}/{end}", handlers.ApiCostsByCat).Methods("GET")
	s.HandleFunc("/ind-cost-month-by-month/{cat}/{start}/{end}", handlers.ApiIndCostMonthByMonth).Methods("GET")
}
