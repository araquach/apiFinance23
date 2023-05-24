package routes

import (
	"github.com/gorilla/mux"
)

var R mux.Router

func FinanceRouter() {
	costs()
	takings()
	profitLoss()
	income()

	return
}
