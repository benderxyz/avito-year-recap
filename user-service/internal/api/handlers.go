package api

import (
	"encoding/json"
	"net/http"

	"user-service/internal/profiles"
)

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profiles.List())
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	profile, ok := profiles.GetByID(id)
	if !ok {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
