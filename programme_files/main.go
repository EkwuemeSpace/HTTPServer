package main

import (
	"fmt"
	"log"
	"net/http"
)

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
