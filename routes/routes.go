package routes

import (
	"github.com/gorilla/mux"
)

func FinanceRouter() *mux.Router {
	r := mux.NewRouter()

	Costs(r)
	Takings(r)
	ProfitLoss(r)
	Income(r)

	return r
}
