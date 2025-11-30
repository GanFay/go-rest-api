package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, id int) {
	var cmdTag pgconn.CommandTag

	cmdTag, err := h.DB.Exec(r.Context(),
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cmdTag.RowsAffected() == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("user deleted"))
}
