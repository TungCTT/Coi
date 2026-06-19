package main

import (
	"fmt"
	"time"
	"net/http"
	"log"
)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s — %v", r.Method, r.URL.Path, time.Since(start))
	})
}
func main() {
	mux:= http.NewServeMux()
	handler:= loggingMiddleware(mux)

	server:= &http.Server{
		Addr: ":8080",
		Handler: handler,
		ReadTimeout: 5* time.Second,
		WriteTimeout: 10* time.Second,
	}
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}