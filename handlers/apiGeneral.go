package handlers

import (
	"encoding/json"
	"github.com/araquach/apiFinance23/models"
	db "github.com/araquach/dbService"
	"log"
	"net/http"
)

func ApiGetStylists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var res []models.TeamMember

	db.DB.Find(&res)

	json, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
