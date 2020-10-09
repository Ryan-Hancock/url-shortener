package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Handler struct {
	router  *mux.Router
	storage Storage
}

func (h *Handler) initialise(s Storage) {
	h.storage = s

	h.router = mux.NewRouter()
	h.initialiseRoutes()
}

func (h *Handler) initialiseRoutes() {
	h.router.HandleFunc("/url/{url}", h.postURL).Methods("POST")
	h.router.HandleFunc("/{key}", h.postURL).Methods("GET")
}

func (h Handler) run(addr string) {
	srv := &http.Server{
		Handler: h.router,
		Addr:    addr,
	}

	log.Fatal(srv.ListenAndServe())
}

func (h Handler) postURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rURL := vars["url"]

	u, err := url.ParseRequestURI(rURL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("not a url: %s", rURL))
	}

	if err := h.storage.Set(GenerateString(10, charset), u.String()); err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "something went wrong")
	}

	respondWithJSON(w, http.StatusCreated, nil)
}

func (h Handler) getURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	foundURL := h.storage.Get(key)
	if foundURL == "" {
		respondWithRedirect(w, r, foundURL)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithRedirect(w http.ResponseWriter, r *http.Request, redirect string) {
	http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
}
