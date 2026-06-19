package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")

	if name == "" {
		name = "Guest"
	}
	w.Write([]byte("Hello " + name))
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodGet {
		w.Write([]byte("use POST to count character"))
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := strconv.Itoa(len(body))
	w.Write([]byte("count: " + data))
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	first := r.URL.Query().Get("a")
	second := r.URL.Query().Get("b")
	opp := r.URL.Query().Get("op")

	a, err := strconv.Atoi(first)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b, err := strconv.Atoi(second)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch opp {
	case "add":
		w.Write([]byte("result: " + strconv.Itoa(a+b)))
	case "subtract":
		w.Write([]byte("result: " + strconv.Itoa(a-b)))
	case "multiply":
		w.Write([]byte("result: " + strconv.Itoa(a*b)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func agentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	head := r.Header.Get("user-agent")
	w.Write([]byte("Your user agent:" + head))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	apiKey := r.Header.Get("X-Api-Key")

	if apiKey != "secret123" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte("Welcome to the dashboard."))
}
func main() {
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/agent", agentHandler)
	http.HandleFunc("/calculate", calculateHandler)
	http.HandleFunc("/ping", pongHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/count", countHandler)

	fmt.Println("starting server at port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}