package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type StatusData struct {
	Status struct {
		Water int
		Wind  int
	}
}

const (
	PORT   = ":8080"
	MAX    = 100
	RELOAD = 15
)

func main() {
	go AutoReloadStatus()
	http.HandleFunc("/", AutoReloadWeb)
	http.ListenAndServe(PORT, nil)
}

func AutoReloadStatus() {
	for {
		data := StatusData{}
		data.Status.Wind = rand.Intn(MAX) + 1
		data.Status.Water = rand.Intn(MAX) + 1
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("fail marshalling data")
		}
		err = ioutil.WriteFile("data.json", jsonData, 0644)
		if err != nil {
			log.Fatal("fail reading data.json")
		}
		time.Sleep(RELOAD * time.Second)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data.json")
	if err != nil {
		log.Fatal("fail reading data.json")
	}
	var statusData StatusData
	err = json.Unmarshal(fileData, &statusData)
	data := make(map[string]string)
	if err != nil {
		log.Fatal(("fail unmarshalling data.json"))
	}
	if statusData.Status.Water < 5 {
		data["waterStatus"] = "aman"
	} else if statusData.Status.Water <= 8 {
		data["waterStatus"] = "siaga"
	} else {
		data["waterStatus"] = "bahaya"
	}
	if statusData.Status.Wind < 6 {
		data["windStatus"] = "aman"
	} else if statusData.Status.Wind <= 15 {
		data["windStatus"] = "siaga"
	} else {
		data["windStatus"] = "bahaya"
	}
	tpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("fail parsing file index.html")
	}
	tpl.Execute(w, data)
}
