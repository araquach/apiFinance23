package handlers

import (
	"encoding/json"
	"github.com/araquach/apiFinance23/models"
	"github.com/araquach/apiHelpers"
	db "github.com/araquach/dbService"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func ApiTakingsByYear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	type JakataTotals struct {
		Year        string  `json:"year"`
		Products    float32 `json:"products"`
		Services    float32 `json:"services"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type PKTotals struct {
		Year        string  `json:"year"`
		Products    float32 `json:"products"`
		Services    float32 `json:"services"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type BaseTotals struct {
		Year        string  `json:"year"`
		Products    float32 `json:"products"`
		Services    float32 `json:"services"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
	}

	type TotalTotals struct {
		Year        string  `json:"year"`
		Products    float32 `json:"products"`
		Services    float32 `json:"services"`
		Total       float32 `json:"total"`
		LinearTotal float32 `json:"linear_total"`
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

	startYear, endYear := sd.Year(), ed.Year()

	jakataYearMap, jakataYears := generateEmptyYearMap(startYear, endYear)
	pkYearMap, pkYears := generateEmptyYearMap(startYear, endYear)
	baseYearMap, baseYears := generateEmptyYearMap(startYear, endYear)
	totalYearMap, totalYears := generateEmptyYearMap(startYear, endYear)

	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(services) AS services, SUM(products) AS products, SUM(services) + SUM(products) AS total FROM takings WHERE salon = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", "jakata", sd, ed).Scan(&jakata)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(services) AS services, SUM(products) AS products, SUM(services) + SUM(products) AS total FROM takings WHERE salon = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", "pk", sd, ed).Scan(&pk)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(services) AS services, SUM(products) AS products, SUM(services) + SUM(products) AS total FROM takings WHERE salon = ? AND date >= ? AND date <= ? GROUP BY year ORDER BY year", "base", sd, ed).Scan(&base)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('year', date),'YYYY') AS year, SUM(services) AS services, SUM(products) AS products, SUM(services) + SUM(products) AS total FROM takings WHERE date >= ? AND date <= ? GROUP BY year ORDER BY year", sd, ed).Scan(&totals)

	jakata = mergeYearData(jakata, jakataYearMap, jakataYears).([]JakataTotals)
	pk = mergeYearData(pk, pkYearMap, pkYears).([]PKTotals)
	base = mergeYearData(base, baseYearMap, baseYears).([]BaseTotals)
	totals = mergeYearData(totals, totalYearMap, totalYears).([]TotalTotals)

	helpers.AddLinearRegressionPoints(jakata, []string{"Total"})
	helpers.AddLinearRegressionPoints(pk, []string{"Total"})
	helpers.AddLinearRegressionPoints(base, []string{"Total"})
	helpers.AddLinearRegressionPoints(totals, []string{"Total"})

	res := Data{
		JakataTotals: jakata,
		PKTotals:     pk,
		BaseTotals:   base,
		TotalTotals:  totals,
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func generateEmptyYearMap(startYear, endYear int) (map[string]float32, []string) {
	yearMap := make(map[string]float32)
	years := make([]string, 0, endYear-startYear+1)
	for year := startYear; year <= endYear; year++ {
		yearStr := strconv.Itoa(year)
		yearMap[yearStr] = 0
		years = append(years, yearStr)
	}
	return yearMap, years
}

func mergeYearData(dbData interface{}, yearMap map[string]float32, years []string) interface{} {
	dataValue := reflect.ValueOf(dbData)
	dataType := dataValue.Type().Elem()

	if dataValue.Kind() != reflect.Slice {
		return nil
	}

	dataLength := dataValue.Len()

	// Update yearMap with values from dbData
	for i := 0; i < dataLength; i++ {
		item := dataValue.Index(i)
		year := item.FieldByName("Year").String()
		total := float32(item.FieldByName("Total").Float())
		yearMap[year] = total
	}

	// Create a new slice to store merged data
	mergedData := reflect.MakeSlice(reflect.SliceOf(dataType), 0, len(yearMap))

	// Iterate through years and update the mergedData slice
	for _, year := range years {
		total := yearMap[year]

		// Find the item in dbData with the same year
		var item reflect.Value
		for i := 0; i < dataLength; i++ {
			tempItem := dataValue.Index(i)
			if tempItem.FieldByName("Year").String() == year {
				item = tempItem
				break
			}
		}

		// If the item was found, update the Total field and append it to mergedData
		if item.IsValid() {
			item.FieldByName("Total").SetFloat(float64(total))
			mergedData = reflect.Append(mergedData, item)
		} else {
			// If the item was not found, create a new item and set the Year and Total fields
			newItem := reflect.New(dataType)
			newItem.Elem().FieldByName("Year").SetString(year)
			newItem.Elem().FieldByName("Total").SetFloat(float64(total))
			mergedData = reflect.Append(mergedData, newItem.Elem())
		}
	}

	return mergedData.Interface()
}

func ApiMonthlyTakingsByDateRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	type JakResult struct {
		MonthTotal     string  `json:"month"`
		Services       float32 `json:"services"`
		Products       float32 `json:"products"`
		Total          float32 `json:"total"`
		LinearServices float32 `json:"linear_services"`
		LinearProducts float32 `json:"linear_products"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type PkResult struct {
		MonthTotal     string  `json:"month"`
		Services       float32 `json:"services"`
		Products       float32 `json:"products"`
		Total          float32 `json:"total"`
		LinearServices float32 `json:"linear_services"`
		LinearProducts float32 `json:"linear_products"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type BaseResult struct {
		MonthTotal     string  `json:"month"`
		Services       float32 `json:"services"`
		Products       float32 `json:"products"`
		Total          float32 `json:"total"`
		LinearServices float32 `json:"linear_services"`
		LinearProducts float32 `json:"linear_products"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type TotalResult struct {
		MonthTotal     string  `json:"month"`
		Services       float32 `json:"services"`
		Products       float32 `json:"products"`
		Total          float32 `json:"total"`
		LinearServices float32 `json:"linear_services"`
		LinearProducts float32 `json:"linear_products"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type JakGrandTotals struct {
		Services   float32 `json:"services"`
		Products   float32 `json:"products"`
		GrandTotal float32 `json:"grand_total"`
		Yearly     float32 `json:"yearly"`
		Monthly    float32 `json:"monthly"`
	}

	type PkGrandTotals struct {
		Services   float32 `json:"services"`
		Products   float32 `json:"products"`
		GrandTotal float32 `json:"grand_total"`
		Yearly     float32 `json:"yearly"`
		Monthly    float32 `json:"monthly"`
	}

	type BaseGrandTotals struct {
		Services   float32 `json:"services"`
		Products   float32 `json:"products"`
		GrandTotal float32 `json:"grand_total"`
		Yearly     float32 `json:"yearly"`
		Monthly    float32 `json:"monthly"`
	}

	type TotalGrandTotals struct {
		Services   float32 `json:"services"`
		Products   float32 `json:"products"`
		GrandTotal float32 `json:"grand_total"`
		Yearly     float32 `json:"yearly"`
		Monthly    float32 `json:"monthly"`
	}

	type Jakata struct {
		Figures []JakResult    `json:"figures"`
		Totals  JakGrandTotals `json:"totals"`
	}

	type PK struct {
		Figures []PkResult    `json:"figures"`
		Totals  PkGrandTotals `json:"totals"`
	}

	type Base struct {
		Figures []BaseResult    `json:"figures"`
		Totals  BaseGrandTotals `json:"totals"`
	}

	type Total struct {
		Figures []TotalResult    `json:"figures"`
		Totals  TotalGrandTotals `json:"totals"`
	}

	type Data struct {
		Jakata Jakata `json:"jakata"`
		PK     PK     `json:"pk"`
		Base   Base   `json:"base"`
		Total  Total  `json:"total"`
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

	var jakataRes []JakResult
	var pkRes []PkResult
	var baseRes []BaseResult
	var totalRes []TotalResult

	var jgt JakGrandTotals
	var pkgt PkGrandTotals
	var bgt BaseGrandTotals
	var tgt TotalGrandTotals

	var jakata Jakata
	var pk PK
	var base Base
	var total Total

	monthsList := GenerateMonthsList(startDate, endDate)

	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month_total, sum(services) AS services, sum(products) as products, sum(services) + sum(products) as total FROM takings WHERE salon = ? AND date BETWEEN ? AND ? GROUP BY salon, month_total ORDER BY month_total", "jakata", sd, ed).Scan(&jakataRes)
	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month_total, sum(services) AS services, sum(products) as products, sum(services) + sum(products) as total FROM takings WHERE salon = ? AND date BETWEEN ? AND ? GROUP BY salon, month_total ORDER BY month_total", "pk", sd, ed).Scan(&pkRes)
	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month_total, sum(services) AS services, sum(products) as products, sum(services) + sum(products) as total FROM takings WHERE salon = ? AND date BETWEEN ? AND ? GROUP BY salon, month_total ORDER BY month_total", "base", sd, ed).Scan(&baseRes)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month_total, sum(services) AS services, sum(products) as products, sum(services) + sum(products) as total FROM takings WHERE date BETWEEN ? AND ? GROUP BY month_total ORDER BY month_total", sd, ed).Scan(&totalRes)

	jakataRes = MergeResultsWithMonthsList(monthsList, jakataRes).([]JakResult)
	pkRes = MergeResultsWithMonthsList(monthsList, pkRes).([]PkResult)
	baseRes = MergeResultsWithMonthsList(monthsList, baseRes).([]BaseResult)
	totalRes = MergeResultsWithMonthsList(monthsList, totalRes).([]TotalResult)

	helpers.AddLinearRegressionPoints(jakataRes, []string{"Services", "Products", "Total"})
	helpers.AddLinearRegressionPoints(pkRes, []string{"Services", "Products", "Total"})
	helpers.AddLinearRegressionPoints(baseRes, []string{"Services", "Products", "Total"})
	helpers.AddLinearRegressionPoints(totalRes, []string{"Services", "Products", "Total"})

	// Calculate total income
	for _, r := range jakataRes {
		jgt.Services += r.Services
		jgt.Products += r.Products
		jgt.GrandTotal += r.Total
		jgt.Yearly = jgt.GrandTotal / float32(mnths) * 12
		jgt.Monthly = jgt.GrandTotal / float32(mnths)
	}

	for _, r := range pkRes {
		pkgt.Services += r.Services
		pkgt.Products += r.Products
		pkgt.GrandTotal += r.Total
		pkgt.Yearly = pkgt.GrandTotal / float32(mnths) * 12
		pkgt.Monthly = pkgt.GrandTotal / float32(mnths)
	}

	for _, r := range baseRes {
		bgt.Services += r.Services
		bgt.Products += r.Products
		bgt.GrandTotal += r.Total
		bgt.Yearly = bgt.GrandTotal / float32(mnths) * 12
		bgt.Monthly = bgt.GrandTotal / float32(mnths)
	}

	for _, r := range totalRes {
		tgt.Services += r.Services
		tgt.Products += r.Products
		tgt.GrandTotal += r.Total
		tgt.Yearly = tgt.GrandTotal / float32(mnths) * 12
		tgt.Monthly = tgt.GrandTotal / float32(mnths)
	}

	jakata.Figures = jakataRes
	jakata.Totals = jgt
	pk.Figures = pkRes
	pk.Totals = pkgt
	base.Figures = baseRes
	base.Totals = bgt
	total.Figures = totalRes
	total.Totals = tgt

	f := Data{
		Jakata: jakata,
		PK:     pk,
		Base:   base,
		Total:  total,
	}

	json, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func GenerateMonthsList(startDate, endDate time.Time) []string {
	months := []string{}
	current := startDate
	for !current.After(endDate) {
		months = append(months, current.Format("2006-01"))
		current = current.AddDate(0, 1, 0)
	}
	return months
}

func MergeResultsWithMonthsList(monthsList []string, results interface{}) interface{} {
	resultsValue := reflect.ValueOf(results)
	resultsType := reflect.TypeOf(results).Elem()

	mergedResults := reflect.MakeSlice(reflect.SliceOf(resultsType), 0, 0)

	monthsMap := make(map[string]bool)

	for i := 0; i < resultsValue.Len(); i++ {
		monthTotal := resultsValue.Index(i).FieldByName("MonthTotal").String()
		monthsMap[monthTotal] = true
	}

	for _, month := range monthsList {
		if _, ok := monthsMap[month]; !ok {
			emptyResult := reflect.New(resultsType).Elem()
			emptyResult.FieldByName("MonthTotal").SetString(month)
			mergedResults = reflect.Append(mergedResults, emptyResult)
		} else {
			for i := 0; i < resultsValue.Len(); i++ {
				result := resultsValue.Index(i)
				if result.FieldByName("MonthTotal").String() == month {
					mergedResults = reflect.Append(mergedResults, result)
					break
				}
			}
		}
	}

	return mergedResults.Interface()
}

func ApiStylistTakingsMonthByMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	type result struct {
		Month          string  `json:"month"`
		Products       float32 `json:"products"`
		Services       float32 `json:"services"`
		Total          float32 `json:"total"`
		LinearProducts float32 `json:"linear_products"`
		LinearServices float32 `json:"linear_services"`
		LinearTotal    float32 `json:"linear_total"`
	}

	vars := mux.Vars(r)
	start := vars["start"]
	end := vars["end"]
	fn := vars["fName"]
	ln := vars["lName"]

	sd, err := time.Parse("2006-01-02", start)
	if err != nil {
		panic(err)
	}
	ed, err := time.Parse("2006-01-02", end)
	if err != nil {
		panic(err)
	}

	var res []result

	db.DB.Raw("SELECT to_char(date_trunc('month', m)::date, 'YYYY-MM') AS month, SUM(t.services) AS services, SUM(t.products) AS products, SUM(t.services) + SUM(t.products) AS total FROM (SELECT generate_series(?::timestamp, ?::timestamp, '1 month') AS m) AS s LEFT JOIN takings t ON date_trunc('month', t.date) = date_trunc('month', s.m) AND t.name ILIKE ? GROUP BY month ORDER BY month", sd, ed, fn+"%"+" "+ln+"%").Scan(&res)

	// Loop through rows and add missing months with 0 values
	for i := 0; i < len(res); i++ {
		// Check if any values are missing and replace with 0 if necessary
		if res[i].Products == 0 && res[i].Services == 0 && res[i].Total == 0 {
			res[i].LinearProducts = 0
			res[i].LinearServices = 0
		}
	}

	helpers.AddLinearRegressionPoints(res, []string{"Products", "Services", "Total"})

	json, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiTakingsByStylistComparison(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	type result struct {
		Stylist  string  `json:"stylist"`
		Products float32 `json:"products"`
		Services float32 `json:"services"`
		Total    float32 `json:"total"`
	}

	vars := mux.Vars(r)
	s := vars["salon"]
	sd := vars["start"]
	ed := vars["end"]

	var res []result

	if s == "all" {
		db.DB.Model(&models.Taking{}).Select("name as stylist, sum(services) as services, sum(products) as products, sum(services) + sum(products) as total").Where("date BETWEEN ? AND ?", sd, ed).Where("services > 0").Where("name != ''").Group("name").Order("total").Find(&res)
	} else {
		db.DB.Model(&models.Taking{}).Select("name as stylist, sum(services) as services, sum(products) as products, sum(services) + sum(products) as total").Where("date BETWEEN ? AND ?", sd, ed).Where("salon", s).Group("name, salon").Order("total").Find(&res)
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}

func ApiTotalsByDateRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	type JakataTotals struct {
		Month          string  `json:"month"`
		Products       float32 `json:"products"`
		Services       float32 `json:"services"`
		Total          float32 `json:"total"`
		LinearProducts float32 `json:"linear_products"`
		LinearServices float32 `json:"linear_services"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type PKTotals struct {
		Month          string  `json:"month"`
		Products       float32 `json:"products"`
		Services       float32 `json:"services"`
		Total          float32 `json:"total"`
		LinearProducts float32 `json:"linear_products"`
		LinearServices float32 `json:"linear_services"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type BaseTotals struct {
		Month          string  `json:"month"`
		Products       float32 `json:"products"`
		Services       float32 `json:"services"`
		Total          float32 `json:"total"`
		LinearProducts float32 `json:"linear_products"`
		LinearServices float32 `json:"linear_services"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type TotalTotals struct {
		Month          string  `json:"month"`
		Products       float32 `json:"products"`
		Services       float32 `json:"services"`
		Total          float32 `json:"total"`
		LinearProducts float32 `json:"linear_products"`
		LinearServices float32 `json:"linear_services"`
		LinearTotal    float32 `json:"linear_total"`
	}

	type Data struct {
		DateFrom     string         `json:"date_from"`
		DateTo       string         `json:"date_to"`
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

	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(products) as products, sum(services) as services, sum(services) + sum(products) AS total FROM takings WHERE salon = 'jakata' AND date BETWEEN ? AND ? GROUP BY salon, month ORDER BY month", sd, ed).Scan(&jakata)
	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(products) as products, sum(services) as services, sum(services) + sum(products) AS total FROM takings WHERE salon = 'pk' AND date BETWEEN ? AND ? GROUP BY salon, month ORDER BY month", sd, ed).Scan(&pk)
	db.DB.Raw("SELECT salon, to_char(DATE_TRUNC('month', date),'YYYY-MM') AS month, sum(products) as products, sum(services) as services, sum(services) + sum(products) AS total FROM takings WHERE salon = 'base' AND date BETWEEN ? AND ? GROUP BY salon, month ORDER BY month", sd, ed).Scan(&base)
	db.DB.Raw("SELECT to_char(DATE_TRUNC('month', date),'YYYY-MM') AS  month, sum(products) as products, sum(services) as services, sum(services) + sum(products) AS total FROM takings WHERE date BETWEEN ? AND ? GROUP BY month ORDER BY month", sd, ed).Scan(&totals)

	helpers.AddLinearRegressionPoints(jakata, []string{"Total", "Products", "Services"})
	helpers.AddLinearRegressionPoints(pk, []string{"Total", "Products", "Services"})
	helpers.AddLinearRegressionPoints(base, []string{"Total", "Products", "Services"})
	helpers.AddLinearRegressionPoints(totals, []string{"Total", "Products", "Services"})

	f := Data{
		DateFrom:     sd,
		DateTo:       ed,
		JakataTotals: jakata,
		PKTotals:     pk,
		BaseTotals:   base,
		TotalTotals:  totals,
	}

	json, err := json.Marshal(f)
	if err != nil {
		log.Println(err)
	}
	w.Write(json)
}
