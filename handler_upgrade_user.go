package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/robertjchandler/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type event struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey", err)
		return
	}
	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := event{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
