package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

type user struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Gmail    string    `json:"gmail"`
	Time     time.Time `json:"time"`
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.usersPost(w, r)
	case http.MethodGet:
		h.usersGet(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) usersGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		intID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "id must be int", http.StatusBadRequest)
			return
		}
		var u user
		err = h.DB.QueryRow(
			r.Context(),
			`SELECT id, name, gmail, password, created_at FROM users WHERE id = $1`,
			intID,
		).Scan(
			&u.ID,
			&u.Name,
			&u.Gmail,
			&u.Password,
			&u.Time,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, "cannot encode json", http.StatusInternalServerError)
		}
		return
	}

	query, err := h.DB.Query(r.Context(), `SELECT id, name, gmail, password, created_at FROM users ORDER BY id`)
	if err != nil {
		return
	}
	defer query.Close()

	var users []user

	for query.Next() {
		u := user{}
		err = query.Scan(&u.ID,
			&u.Name,
			&u.Gmail,
			&u.Password,
			&u.Time)
		if err != nil {
			http.Error(w, "cannot scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if err = query.Err(); err != nil {
		http.Error(w, "rows error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "cannot encode json", http.StatusInternalServerError)
		return
	}
}

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
