package main

import (
	//"context"
	"encoding/json"
	"fmt"
	//"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	//"net/url"
)

type handler struct {
	db *mydatabase
}

func (h *handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "ok")
}

func (h *handler) artistsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := *h.db
	artists := db.get("artists")

	b, err := json.Marshal(artists)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func (h *handler) showHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := *h.db
	shows := db.get("shows")

	// Filter by id
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var filteredShows []*Show
	for _, obj := range shows {
		show := obj.(*Show)
		if show.ArtistId == int64(id) {
			filteredShows = append(filteredShows, show)
		}
	}

	b, err := json.Marshal(filteredShows)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func (h *handler) ticketHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := *h.db
	shows := db.get("shows")

	// Filter by id
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var filteredShows []*Show
	for _, obj := range shows {
		show := obj.(*Show)
		if show.Id == int64(id) {
			filteredShows = append(filteredShows, show)
		}
	}

	b, err := json.Marshal(filteredShows)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func (h *handler) buyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := *h.db

	// Filter by id
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	buySuccess := db.buy(int64(id))

	b, err := json.Marshal(buySuccess)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}
