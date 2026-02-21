package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	const characterLimit = 140

	decoder := json.NewDecoder(r.Body)
	payload := chirp{}
	err := decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error decoding chirp: %s", err))
		return
	}
	if len(payload.Body) > characterLimit {
		respondWithError(w, 400, fmt.Sprintf("Chirp is too long: %s", err))
		return
	}
	respBody := returnVals{
		Body:  payload,
		Valid: true,
	}
	w.Header().Set("Content-Type", "application/json")
	respondWithJSON(w, 200, respBody)
}

type chirp struct {
	Body string `json:"body"`
}
type returnVals struct {
	Body  chirp `json:"body"`
	Valid bool  `json:"valid"`
	error string
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respBody := returnVals{
		error: msg,
	}
	w.WriteHeader(code)
	w.Write(([]byte)(respBody.error))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
