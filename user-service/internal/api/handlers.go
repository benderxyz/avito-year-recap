package api

import (
	"net/http"

	"user-service/internal/profiles"
)

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, profiles.List())
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	profile, ok := profiles.GetByID(id)
	if !ok {
		writeError(w, http.StatusNotFound, "profile not found")
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
