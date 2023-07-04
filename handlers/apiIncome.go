package handlers

import (
	"encoding/json"
	"github.com/araquach/apiHelpers"
	"github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func ApiMonthlyIncomeByDateRange(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type GrandTotals struct {
		GrandTotal float32 `json:"grand_total"`
		Yearly     float32 `json:"yearly"`
		Monthly    float32 `json:"monthly"`
	}

	type Data struct {
		Figures []Result    `json:"figures"`
		Totals  GrandTotals `json:"totals"`
	}

	vars := mux.Vars(r)
	sd := vars["start"]
	ed := vars["end"]

	startDate, err := time.Parse("2006-01-02", sd)
	if err != nil {
		panic(err)
	}
	endDate, err := time.Parse("2006-01-02", ed)
	if err != nil {
		panic(err)
	}

	mnths := helpers.MonthsCount(startDate, endDate)

	var res []Result

	var gt GrandTotals

	var data Data

	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(deposit) as total FROM incomes WHERE date BETWEEN ? AND ? GROUP BY month ORDER BY month", sd, ed).Scan(&res)

	helpers.AddLinearRegressionPoints(res, []string{"Total"})

	// Calculate total income

	for _, r := range res {
		gt.GrandTotal += r.Total
		gt.Yearly = gt.GrandTotal / float32(mnths) * 12
		gt.Monthly = gt.GrandTotal / float32(mnths)
	}

	data.Figures = res
	data.Totals = gt

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Error encoding data to JSON", http.StatusInternalServerError)
	}
	w.Write(jsonData)
}
