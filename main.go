package main

import (
	"github.com/araquach/apiFinance23/routes"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db.DBInit(dsn)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Load API Routes
	financeRouter := routes.FinanceRouter()
	mainRouter := mux.NewRouter()

	mainRouter.PathPrefix("/api/finance").Handler(financeRouter)

	log.Printf("Starting server on %s", port)

	http.ListenAndServe(":"+port, mainRouter)
}
