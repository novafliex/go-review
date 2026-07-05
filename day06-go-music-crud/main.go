package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Music struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Year   int    `json:"year"`
}

var musics []Music

func getMusics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(musics)
}

func getMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range musics {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Music is not found.", http.StatusNotFound)
}

func deleteMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range musics {
		if item.ID == params["id"] {
			musics = append(musics[:index], musics[index+1:]...)
			json.NewEncoder(w).Encode(musics)
			return
		}
	}
	http.Error(w, "deletemusic id is not found.", http.StatusNotFound)
}

func createMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var music Music
	err := json.NewDecoder(r.Body).Decode(&music)
	if err != nil {
		http.Error(w, "invalid json data", http.StatusBadRequest)
		return
	}
	music.ID = strconv.Itoa(rand.Intn(100000000))
	musics = append(musics, music)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(music)

}

func updateMusic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index := range musics {
		if musics[index].ID == params["id"] {
			var music Music
			if err := json.NewDecoder(r.Body).Decode(&music); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			music.ID = params["id"]
			musics[index] = music
			json.NewEncoder(w).Encode(music)
			return
		}
	}
	http.Error(w, "Music is not found.", http.StatusNotFound)
}

func main() {
	// create a new router
	r := mux.NewRouter()

	// create new music
	musics = append(musics, Music{
		ID:     "1",
		Title:  "晴天",
		Artist: "周杰伦",
		Album:  "叶惠美",
		Year:   2003,
	})

	musics = append(musics, Music{
		ID:     "2",
		Title:  "十年",
		Artist: "陈奕迅",
		Album:  "黑白灰",
		Year:   2003,
	})

	// setup routes
	r.HandleFunc("/musics", getMusics).Methods("GET")
	r.HandleFunc("/musics/{id}", getMusic).Methods("GET")
	r.HandleFunc("/musics", createMusic).Methods("POST")
	r.HandleFunc("/musics/{id}", updateMusic).Methods("PUT")
	r.HandleFunc("/musics/{id}", deleteMusic).Methods("DELETE")

	// start server
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))

}
