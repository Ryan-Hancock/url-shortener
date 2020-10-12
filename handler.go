package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// Handler represents a web app handler
type Handler struct {
	router  *mux.Router
	storage Storage

	address string
}

// Response for json response
type Response struct {
	Message interface{} `json:"message"`
}

// URLRequest for url requests
type URLRequest struct {
	URL string `json:"url"`
}

func (h *Handler) initialise(s Storage) {
	h.storage = s

	h.router = mux.NewRouter()
	h.initialiseRoutes()
}

func (h *Handler) initialiseRoutes() {
	h.router.HandleFunc("/url", h.postURL).Methods("POST")
	h.router.HandleFunc("/{key}", h.getURL).Methods("GET")
}

func (h *Handler) run(addr string) {
	h.address = addr
	srv := &http.Server{
		Handler: h.router,
		Addr:    addr,
	}

	log.Fatal(srv.ListenAndServe())
}

func (h Handler) postURL(w http.ResponseWriter, r *http.Request) {
	var uReq URLRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&uReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()

	u, err := url.ParseRequestURI(uReq.URL)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Not a url: %s example: http://google.com", uReq.URL))
		return
	}

	shortURL := GenerateString(10, charset)

	if err := h.storage.Set(shortURL, u.String()); err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{Message: fmt.Sprintf("%s/%s", h.address, shortURL)})
}

func (h Handler) getURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	foundURL := h.storage.Get(key)

	if foundURL != "" {
		respondWithRedirect(w, r, foundURL)
		return
	}

	respondWithError(w, http.StatusNotFound, fmt.Sprintf("Not found %s", key))
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
	http.Redirect(w, r, redirect, http.StatusMovedPermanently)
}
