package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func getIDFromRequest(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, fmt.Errorf("id is required")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}
