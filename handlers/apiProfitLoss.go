package handlers

import (
	"encoding/json"
	"github.com/araquach/apiHelpers"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ApiCostsAndTakings(w http.ResponseWriter, r *http.Request) {
	type Costs struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type Takings struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type Result struct {
		Costs          []Costs   `json:"costs"`
		Takings        []Takings `json:"takings"`
		TotalCosts     float32   `json:"total_costs"`
		AverageCosts   float32   `json:"average_costs"`
		TotalTakings   float32   `json:"total_takings"`
		AverageTakings float32   `json:"average_takings"`
	}

	vars := mux.Vars(r)
	sd := vars["start"]
	ed := vars["end"]

	var c []Costs
	var t []Takings
	var res Result

	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(debit) AS total FROM costs WHERE date BETWEEN ? AND ? GROUP BY month ORDER BY month", sd, ed).Scan(&c)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(services) + sum(products) AS total FROM takings WHERE date BETWEEN ? AND ? GROUP BY month ORDER BY month", sd, ed).Scan(&t)

	helpers.AddLinearRegressionPoints(t, []string{"Total"})
	helpers.AddLinearRegressionPoints(c, []string{"Total"})

	var totalCosts, totalTakings float32
	for _, cost := range c {
		totalCosts += cost.Total
	}
	for _, taking := range t {
		totalTakings += taking.Total
	}

	res = Result{
		Costs:          c,
		Takings:        t,
		TotalCosts:     totalCosts,
		AverageCosts:   totalCosts / float32(len(c)),
		TotalTakings:   totalTakings,
		AverageTakings: totalTakings / float32(len(t)),
	}

	myJson, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(myJson)
	if err != nil {
		return
	}
}
