package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (h *Handler) usersGet(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) usersGetByID(w http.ResponseWriter, r *http.Request, id int) {
	var u user
	err := h.DB.QueryRow(
		r.Context(),
		`SELECT id, name, gmail, password, created_at FROM users WHERE id = $1`,
		id,
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
	if err = json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "cannot encode json", http.StatusInternalServerError)
	}
	return
}
