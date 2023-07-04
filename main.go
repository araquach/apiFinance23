package main

import (
	"github.com/araquach/apiFinance23/routes"
	"github.com/araquach/apiHelpers/middleware"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // replace with your domain
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // or other methods you need
		AllowedHeaders:   []string{"*"},
	})

	// Load API Routes
	financeRouter := routes.FinanceRouter()
	mainRouter := mux.NewRouter()

	mainRouter.PathPrefix("/api/finance").Handler(financeRouter)

	mainRouter.Use(middleware.ContentTypeMiddleware)
	mainRouter.Use(c.Handler)

	log.Printf("Starting server on %s", port)

	http.ListenAndServe(":"+port, mainRouter)
}
