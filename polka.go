package main

import (
	"encoding/json"
	"net/http"
)

func (cfg apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type parameters struct {
		Event string         `json:"event"`
		Data  map[string]int `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, "OK")
		return
	}

	// Get User ID from params data and attempt to upgrade the user
	userID := params.Data["user_id"]
	err = cfg.DB.UpgradeUser(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// User upgraded successfully
	respondWithJSON(w, http.StatusOK, "OK")
}
