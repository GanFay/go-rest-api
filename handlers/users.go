package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type user struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Gmail    string    `json:"gmail"`
	Time     time.Time `json:"time"`
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	splPath := strings.Split(path, "/")
	var done []string
	for i := 0; i < len(splPath); i++ {
		if splPath[i] != "" {
			done = append(done, splPath[i])
		}
	}
	switch len(done) {
	case 1:
		switch r.Method {
		case http.MethodGet:
			h.usersGet(w, r)
		case http.MethodPost:
			h.usersPost(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}
	case 2:
		id, err := strconv.Atoi(done[1])
		if err != nil && id != 0 {
			return
		}
		switch r.Method {
		case "GET":
			h.usersGetByID(w, r, id)
		case "DELETE":
			h.delete(w, r, id)
		//case "PATCH": // меняет true or false смотря от запроса
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		//switch {
		//case splPath[1] == "edit":
		//
		//}
	}
}
