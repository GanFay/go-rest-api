package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func (h *Handler) usersPost(w http.ResponseWriter, r *http.Request) {
	var req user

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Password == "" || req.Gmail == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	var id int
	var Time time.Time
	time.Sleep(time.Second * 5)

	err := h.DB.QueryRow(r.Context(), "INSERT into users (name, password, gmail) VALUES ($1, $2, $3) RETURNING id, created_at", req.Name, req.Password, req.Gmail).Scan(&id, &Time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.ID = id
	req.Time = Time
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(req)
}
