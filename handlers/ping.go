package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("pong")
	err := json.NewEncoder(w).Encode("pong")
	if err != nil {
		log.Fatal(err)
		return
	}
}
