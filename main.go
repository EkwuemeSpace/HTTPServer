package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// pongHandler handles GET /ping.
// Simplest possible endpoint — always returns "pong" with a 200 OK.
// No method check here, so technically any HTTP method hits this handler.
func pongHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// helloHandler handles GET /hello.
// Reads an optional "name" query parameter and greets the caller.
// Rejects any non-GET method with 405 Method Not Allowed.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Only GET is allowed on this route.
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	// Query params always arrive as strings — a URL is just text.
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest!" // default when no name is provided
	}

	// NOTE: "Guest!" already has an exclamation mark, and we append another
	// "!" below — for the Guest case this produces "Hello, Guest!!" (double
	// exclamation). It still passes substring-match tests, but is a small bug.
	w.Write([]byte("Hello, " + name + "!"))
}

// countHandler handles both GET and POST /count.
// GET returns usage instructions. POST reads the request body and returns
// its character count.
func countHandler(w http.ResponseWriter, r *http.Request) {
	// Only GET and POST are allowed; anything else is rejected.
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	// GET branch: just return instructions, no body to read.
	if r.Method == http.MethodGet {
		// NOTE: typo here — "ith" should be "with".
		w.Write([]byte("Send a POST request ith text to count words"))
		return
	}

	// POST branch: read the full request body.
	// defer ensures the body is closed once this function returns,
	// freeing the underlying network connection.
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Reading failed for some reason on the server's end.
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	// len(body) works directly on the []byte — no need to convert to string first.
	data := strconv.Itoa(len(body))
	w.Write([]byte("count: " + data))
}

// calculateHandler handles GET /calculate?a=X&b=Y&op=add|subtract|multiply.
// Validates both numbers and the operation before computing a result.
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	// Pull all three query params as raw strings first.
	first := r.URL.Query().Get("a")
	second := r.URL.Query().Get("b")
	opp := r.URL.Query().Get("op")

	// Convert "a" to an int. Atoi returns an error if it's not a valid number.
	a, err := strconv.Atoi(first)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 — invalid number
		return
	}

	// Same conversion + validation for "b".
	b, err := strconv.Atoi(second)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400 — invalid number
		return
	}

	// switch handles all valid operations in one place.
	// default catches any unrecognized op string and returns 400 —
	// this means we don't need a separate validity check before the switch.
	switch opp {
	case "add":
		w.Write([]byte("result: " + strconv.Itoa(a+b)))
	case "subtract":
		w.Write([]byte("result: " + strconv.Itoa(a-b)))
	case "multiply":
		w.Write([]byte("result: " + strconv.Itoa(a*b)))
	default:
		w.WriteHeader(http.StatusBadRequest) // 400 — unknown operation
		return
	}
}

// agentHandler handles GET /agent.
// Echoes back the client's User-Agent header.
func agentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	// Header names are case-insensitive — "user-agent" and "User-Agent"
	// both retrieve the same value.
	head := r.Header.Get("user-agent")
	w.Write([]byte("Your user agent:" + head))
}

// dashboardHandler handles GET /dashboard.
// A protected route — requires a matching X-Api-Key header value.
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	// Header names are case-insensitive, but header VALUES are not —
	// "secret123" and "Secret123" would be treated as different keys.
	apiKey := r.Header.Get("X-Api-Key")
	if apiKey != "secret123" {
		w.WriteHeader(http.StatusUnauthorized) // 401 — missing or wrong key
		return
	}

	w.Write([]byte("Welcome to the dashboard."))
}

// legacy handles GET /legacy.
// Issues a 301 permanent redirect to /v2. This sends only a status code
// and a Location header — the client makes a separate follow-up request
// to actually fetch /v2's content.
func legacy(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/v2", http.StatusMovedPermanently) // 301
}

// v2LegacyHandler handles GET /v2.
// The new destination that /legacy redirects to.
func v2LegacyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to version 2 legacy"))
}

// main wires up all routes and starts the server.
func main() {
	// http.HandleFunc registers a handler function for each URL path.
	// Any path not registered here automatically returns 404 Not Found.
	http.HandleFunc("/legacy", legacy)
	http.HandleFunc("/v2", v2LegacyHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/agent", agentHandler)
	http.HandleFunc("/calculate", calculateHandler)
	http.HandleFunc("/ping", pongHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/count", countHandler)

	fmt.Println("starting server at port :8080")

	// ListenAndServe blocks forever, listening on port 8080.
	// log.Fatal logs and exits if the server fails to start
	// (e.g. port already in use).
	log.Fatal(http.ListenAndServe(":8080", nil))
}
