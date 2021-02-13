package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type CrewMember struct {
	Launches []string `json:"launches"`
}
type Response struct {
	Number int `json:"crew_members"`
}
type Response2 struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Kg   int    `json:"total_payload_weights"`
}
type Rocket struct {
	Name           string    `json:"name"`
	Id             string    `json:"id"`
	PayloadWeights []Payload `json:"payload_weights"`
}
type Launch struct {
	Rocket  string `json:"rocket"`
	Success bool   `json:"success"`
	Id      string `json:"id"`
}
type Payload struct {
	Kg int `json:"kg"`
}

// Returns number of crew members that were in space.
func TotalCrewSent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := http.Get("https://api.spacexdata.com/v4/crew")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var (
			data, _     = ioutil.ReadAll(response.Body)
			crewMembers = []CrewMember{}
			j           = 0
			err         = json.Unmarshal([]byte(data), &crewMembers)
		)
		if err != nil {
			fmt.Println("error:", err)
		}
		for i := range crewMembers {
			if len(crewMembers[i].Launches) > 0 {
				j++
			}
		}
		apiResponse := Response{
			Number: j,
		}
		json.NewEncoder(w).Encode(apiResponse)
	}
}

// Returns sum of all payloads weight that were in space on rocket named Falcon
func FalconLoad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var (
		rockets       = []Rocket{}
		response, err = http.Get("https://api.spacexdata.com/v4/rockets")
	)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var (
			data, _ = ioutil.ReadAll(response.Body)
			err     = json.Unmarshal([]byte(data), &rockets)
		)

		if err != nil {
			fmt.Println("error:", err)
		}
	}
	response, err = http.Get("https://api.spacexdata.com/v4/launches")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		var (
			data, _  = ioutil.ReadAll(response.Body)
			launches = []Launch{}
			err      = json.Unmarshal([]byte(data), &launches)
		)

		if err != nil {
			fmt.Println("error:", err)
		}
		for i := range rockets {
			matched, err := regexp.MatchString(`.*Falcon.*`, rockets[i].Name)
			if err != nil {
				fmt.Println("error:", err)
			}
			if !matched {
				continue
			}
			var (
				flights          = 0                         // Number of launches per rocket
				payloadWeight    = 0                         // Weight of one payload
				totPayloadWeight = 0                         // Weight of all payloads on one rocket
				totWeight        = 0                         // Total weight of all payloads on rocket each flight
				payloads         = rockets[i].PayloadWeights // Payload object
			)

			for j := range payloads {
				payloadWeight = payloads[j].Kg
				totPayloadWeight = payloadWeight + totPayloadWeight
			}
			for j := range launches {
				if (launches[j].Rocket == rockets[i].Id) && (launches[j].Success == true) {
					flights++
				}
			}
			totWeight = totPayloadWeight * flights
			apiResponse := Response2{
				Name: rockets[i].Name,
				Id:   rockets[i].Id,
				Kg:   totWeight,
			}
			json.NewEncoder(w).Encode(apiResponse)
		}
	}
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/load", FalconLoad).Methods(http.MethodGet)
	api.HandleFunc("/crew", TotalCrewSent).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
