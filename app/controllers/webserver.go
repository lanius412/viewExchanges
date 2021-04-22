package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"viewExchange/app/models"
	"viewExchange/config"
)

var templates = template.Must(template.ParseFiles("app/views/chart.html"))

type JsonData struct {
	Pair string `json:"pair"`
	Open string `json:"open"`
}

var data JsonData
var OpenExchange, CurrencyPair string

func viewChartHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "chart.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func receiveDataHandler(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Fatalln("receive error:", err)
	}
	OpenExchange = data.Open
	CurrencyPair = data.Pair
	StreamIngestionData()
}

type JsonError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func APIError(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonError, err := json.Marshal(JsonError{Error: errMsg, Code: code})
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(jsonError)
}

var apiValidPath = regexp.MustCompile("^/api/candle/$")

func apiMakeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := apiValidPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			APIError(w, "Not Found", http.StatusNotFound)
		}
		fn(w, r)
	}
}

func apiCandleHandler(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		APIError(w, "No Pair Param", http.StatusBadRequest)
		return
	}

	strLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(strLimit)
	if strLimit == "" || err != nil || limit < 0 || limit > 1000 {
		limit = 1000
	}

	duration := r.URL.Query().Get("duration")
	durations := map[string]int{
		"1m":  0,
		"5m":  1,
		"30m": 2,
		"1h":  3,
	}
	if duration == "" {
		duration = "5m"
	}
	durationTime := config.Config.Durations[durations[duration]]

	var exchange string
	if data.Open == "" {
		exchange = "bitflyer"
	} else {
		exchange = data.Open
	}

	df, err := models.GetAllCandle(exchange, pair, durationTime, limit)
	if err != nil {
		log.Fatalln(err)
	}

	js, err := json.Marshal(df)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func StartWebServer() {
	http.HandleFunc("/api/candle/", apiMakeHandler(apiCandleHandler))
	http.HandleFunc("/receive", receiveDataHandler)
	http.HandleFunc("/chart/", viewChartHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
