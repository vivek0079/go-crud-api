package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "476228",
		Title: "Movie one",
		Director: &Director{
			Firstname: "John",
			Lastname:  "Doe",
		},
	})

	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "987729",
		Title: "Movie two",
		Director: &Director{
			Firstname: "Lorum",
			Lastname:  "Ipsum",
		},
	})

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movie", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("Welcome"))
	if err != nil {
		return
	}
}

func getMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	err := json.NewEncoder(writer).Encode(movies)
	if err != nil {
		return
	}
}

func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	params := mux.Vars(request)
	id := params["id"]
	var m Movie

	for _, movie := range movies {
		if movie.ID == id {
			m = movie
			break
		}
	}

	if m == (Movie{}) {
		writer.WriteHeader(http.StatusNotFound)
		_, err := writer.Write([]byte("Movie not found"))
		if err != nil {
			return
		}
	} else {
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(m)
		if err != nil {
			return
		}
	}
}

func createMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(99999))
	movies = append(movies, movie)
	writer.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(writer).Encode(movie)
	if err != nil {
		return
	}
}

func deleteMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-type", "application/json")
	params := mux.Vars(request)
	id := params["id"]

	for idx, movie := range movies {
		if movie.ID == id {
			movies[idx] = movies[len(movies)-1]
			movies = movies[:len(movies)-1]
			break
		}
	}
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(movies)
	if err != nil {
		return
	}
}

func updateMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id := mux.Vars(request)["id"]

	var newMovie Movie
	_ = json.NewDecoder(request.Body).Decode(&newMovie)

	for idx, item := range movies {
		if item.ID == id {
			movies[idx] = newMovie
		}
	}
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(movies)
	if err != nil {
		return
	}
}
