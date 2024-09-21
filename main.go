package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("hello %s", req.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
