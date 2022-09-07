package handler

import (
	"encoding/json"
	"golang-crud-sql/model"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

var weather model.Weather

const weatherTemplate = "html/weather.html"
const weatherJson = "config/weather.json"

func WeatherHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content Type", "text/html")
	// read from json file and write to webData
	file, _ := ioutil.ReadFile(weatherJson)
	json.Unmarshal(file, &weather)

	view := model.WeatherTemplate{
		Water: weather.Status.Water,
		Wind:  weather.Status.Wind,
	}

	if (weather.Status.Water <= 5) || (weather.Status.Wind <= 6) {
		view.Status = "AMAN"
	} else if ((weather.Status.Water > 5) && (weather.Status.Water <= 8)) || ((weather.Status.Wind > 6) && (weather.Status.Water <= 14)) {
		view.Status = "SIAGA"
	} else if (weather.Status.Water > 8) || (weather.Status.Wind > 15) {
		view.Status = "BAHAYA"
	} else {
		view.Status = "TIDAK DIKETAHUI"
	}

	html, err := template.ParseFiles(weatherTemplate)
	if err != nil {
		http.Error(w, "Error Parsing Data To HTML", http.StatusInternalServerError)
		return
	}

	html.Execute(w, view)
}

func GenerateRandomWeather() {
	for {
		weather.Status.Water = rand.Intn(100)
		weather.Status.Wind = rand.Intn(100)

		// write to json file
		jsonString, _ := json.Marshal(&weather)
		ioutil.WriteFile(weatherJson, jsonString, os.ModePerm)

		// sleep for 15 seconds
		time.Sleep(15 * time.Second)
	}
}
