package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
)

type Users struct {
	ID            int    `json:"id"`
	Gmail         string `json:"gmail"`
	AdminPassword string `json:"admin_password"`
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Users
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	AP := os.Getenv("A_P")

	if req.AdminPassword != AP {
		http.Error(w, "Admin password wrong", http.StatusUnauthorized)
		return
	}

	var cmdTag pgconn.CommandTag

	switch {
	case req.ID != 0:
		cmdTag, err = h.DB.Exec(r.Context(),
			`DELETE FROM users WHERE id = $1`, req.ID)

	case req.Gmail != "":
		cmdTag, err = h.DB.Exec(r.Context(),
			`DELETE FROM users WHERE gmail = $1`, req.Gmail)

	default:
		http.Error(w, "ID, Name or Gmail required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		fmt.Println("delete error:", err)
		return
	}

	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("user deleted"))
}
