package handlers

import (
	"encoding/json"
	"github.com/araquach/apiFinance23/models"
	helpers "github.com/araquach/apiHelpers"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func ApiCostsYearByYear(w http.ResponseWriter, _ *http.Request) {
	type JakataTotals struct {
		Year  string  `json:"year"`
		Total float32 `json:"total"`
	}

	type PKTotals struct {
		Year  string  `json:"year"`
		Total float32 `json:"total"`
	}

	type BaseTotals struct {
		Year  string  `json:"year"`
		Total float32 `json:"total"`
	}

	type TotalTotals struct {
		Year  string  `json:"year"`
		Total float32 `json:"total"`
	}

	type Data struct {
		JakataTotals []JakataTotals `json:"jakata"`
		PKTotals     []PKTotals     `json:"pk"`
		BaseTotals   []BaseTotals   `json:"base"`
		TotalTotals  []TotalTotals  `json:"all"`
	}

	sd, err := time.Parse("2006-01-02", "2017-01-01")
	if err != nil {
		panic(err)
	}
	ed := time.Now()

	var jakata []JakataTotals
	var pk []PKTotals
	var base []BaseTotals
	var totals []TotalTotals

	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(debit) AS total FROM costs WHERE account = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", helpers.GetBankAcc("jakata"), sd, ed).Scan(&jakata)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(debit) AS total FROM costs WHERE account = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", helpers.GetBankAcc("pk"), sd, ed).Scan(&pk)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(debit) AS total FROM costs WHERE account = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", helpers.GetBankAcc("base"), sd, ed).Scan(&base)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(debit) AS total FROM costs WHERE date >= ? AND date <= ? GROUP BY year ORDER BY year", sd, ed).Scan(&totals)

	fillMissingBaseTotals := func(start, end time.Time, base []BaseTotals) []BaseTotals {
		yearMap := make(map[string]float32)

		for y := start.Year(); y <= end.Year(); y++ {
			yearMap[strconv.Itoa(y)] = 0
		}

		for _, b := range base {
			yearMap[b.Year] = b.Total
		}

		var filledBase []BaseTotals
		for year, total := range yearMap {
			filledBase = append(filledBase, BaseTotals{Year: year, Total: total})
		}

		// Sort the years in the filledBase slice
		sort.Slice(filledBase, func(i, j int) bool {
			return filledBase[i].Year < filledBase[j].Year
		})

		return filledBase
	}

	base = fillMissingBaseTotals(sd, ed, base)

	res := Data{
		JakataTotals: jakata,
		PKTotals:     pk,
		BaseTotals:   base,
		TotalTotals:  totals,
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

func ApiCostsMonthByMonth(w http.ResponseWriter, r *http.Request) {
	type JakataTotals struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type PKTotals struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type BaseTotals struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type TotalTotals struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type Data struct {
		JakataTotals []JakataTotals `json:"jakata"`
		PKTotals     []PKTotals     `json:"pk"`
		BaseTotals   []BaseTotals   `json:"base"`
		TotalTotals  []TotalTotals  `json:"all"`
	}

	vars := mux.Vars(r)
	sd := vars["start"]
	ed := vars["end"]

	var jakata []JakataTotals
	var pk []PKTotals
	var base []BaseTotals
	var totals []TotalTotals

	db.DB.Raw("SELECT account as salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(debit) AS total FROM costs WHERE account = '06517160' AND date BETWEEN ? AND ? GROUP BY account, month ORDER BY month", sd, ed).Scan(&jakata)
	db.DB.Raw("SELECT account as salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS  month, sum(debit) AS total FROM costs WHERE account = '02017546' AND date BETWEEN ? AND ? GROUP BY account, month ORDER BY month", sd, ed).Scan(&pk)
	db.DB.Raw("SELECT account as salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS  month, sum(debit) AS total FROM costs WHERE account = '17623364' AND date BETWEEN ? AND ? GROUP BY account, month ORDER BY month", sd, ed).Scan(&base)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS  month, sum(debit) AS total FROM costs WHERE date BETWEEN ? AND ? GROUP BY month ORDER BY month", sd, ed).Scan(&totals)

	fillMissingBaseMonths := func(start, end string, base []BaseTotals) []BaseTotals {
		startDate, _ := time.Parse("2006-01-02", start)
		endDate, _ := time.Parse("2006-01-02", end)
		monthMap := make(map[string]float32)

		for d := startDate; d.Before(endDate) || d.Equal(endDate); d = d.AddDate(0, 1, 0) {
			monthMap[d.Format("2006-01")] = 0
		}

		for _, b := range base {
			monthMap[b.Month] = b.Total
		}

		var filledBase []BaseTotals
		for month, total := range monthMap {
			filledBase = append(filledBase, BaseTotals{Month: month, Total: total})
		}

		// Sort the months in the filledBase slice
		sort.Slice(filledBase, func(i, j int) bool {
			return filledBase[i].Month < filledBase[j].Month
		})

		return filledBase
	}

	base = fillMissingBaseMonths(sd, ed, base)

	helpers.AddLinearRegressionPoints(jakata, []string{"Total"})
	helpers.AddLinearRegressionPoints(pk, []string{"Total"})
	helpers.AddLinearRegressionPoints(base, []string{"Total"})
	helpers.AddLinearRegressionPoints(totals, []string{"Total"})

	f := Data{
		JakataTotals: jakata,
		PKTotals:     pk,
		BaseTotals:   base,
		TotalTotals:  totals,
	}

	myJson, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(myJson)
	if err != nil {
		return
	}
}

func ApiCostsByCat(w http.ResponseWriter, r *http.Request) {
	var t float32

	type Result struct {
		Category string  `json:"category"`
		Total    float32 `json:"total"`
		Percent  float32 `json:"percent"`
		Average  float32 `json:"average"`
	}

	type Data struct {
		Salon      string   `json:"salon"`
		GrandTotal float32  `json:"grand_total"`
		ByYear     float32  `json:"by_year"`
		Months     int      `json:"months"`
		Figures    []Result `json:"figures"`
	}

	vars := mux.Vars(r)
	s := vars["salon"]
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

	if s == "all" {
		db.DB.Model(&models.Cost{}).Order("total desc").Select("category, sum(debit) as total").Where("date BETWEEN ? AND ?", sd, ed).Group("category").Find(&res)
	} else {
		db.DB.Model(&models.Cost{}).Order("total desc").Select("account, category, sum(debit) as total").Where("date BETWEEN ? AND ?", sd, ed).Where("account", helpers.GetBankAcc(s)).Group("account, category").Find(&res)
	}

	// Calculate total costs
	for _, r := range res {
		t += r.Total
	}

	for k, v := range res {
		//remove Izzys Wage and loans from drawings
		if v.Category == "drawings" {
			(res)[k].Total = (res)[k].Total - ((2200 + 450.23 + 291) * float32(mnths))
		}
		if v.Category == "loans" {
			(res)[k].Total = (res)[k].Total + ((450 + 291) * float32(mnths))
		}
		if v.Category == "wages" {
			(res)[k].Total = (res)[k].Total + (2200 * float32(mnths))
		}

		// calulate averagesa
		(res)[k].Average = (res)[k].Total / float32(mnths)
		(res)[k].Percent = ((res)[k].Total / t) * 100
	}

	f := Data{
		Salon:      s,
		GrandTotal: t,
		ByYear:     t / float32(mnths) * 12,
		Months:     mnths,
		Figures:    res,
	}

	myJson, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(myJson)
	if err != nil {
		return
	}
}

func ApiIndCostMonthByMonth(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Month       string  `json:"month"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type Data struct {
		Category string   `json:"category"`
		Figures  []Result `json:"figures"`
	}

	vars := mux.Vars(r)
	cat := vars["cat"]
	start := vars["start"]
	end := vars["end"]

	sd, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic(err)
	}
	ed, err := time.Parse("2006-01-02", end)
	if err != nil {
		panic(err)
	}

	var res []Result

	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, SUM(debit) AS total FROM costs WHERE date >= ? AND date <= ? AND category ILIKE ? GROUP BY month ORDER BY month", sd, ed, cat+"%").Scan(&res)

	// Create a map of months between start and end dates with 0 value
	monthData := make(map[string]float32)
	var orderedMonths []string
	for d := sd; d.Before(ed) || d.Equal(ed); d = d.AddDate(0, 1, 0) {
		month := d.Format("2006-01")
		monthData[month] = 0
		orderedMonths = append(orderedMonths, month)
	}

	// Update the map with the results from the database
	for _, r := range res {
		monthData[r.Month] = r.Total
	}

	// Convert the map to a slice of Result structs while maintaining the order
	var filledResults []Result
	for _, month := range orderedMonths {
		filledResults = append(filledResults, Result{
			Month: month,
			Total: monthData[month],
		})
	}

	helpers.AddLinearRegressionPoints(filledResults, []string{"Total"})

	f := Data{
		Category: cat,
		Figures:  filledResults,
	}

	myJson, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	_, err = w.Write(myJson)
	if err != nil {
		return
	}
}
