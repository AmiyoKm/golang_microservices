package main

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse json data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if reqBody.UserID == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}
	log.Println("Success")

	res := contracts.APIResponse{Data: reqBody}
	writeJSON(w, http.StatusCreated, res)
}
