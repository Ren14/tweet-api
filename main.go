package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/cmd/http/dependencies"
	"github.com/renzonaitor/tweet-api/cmd/http/routes"
)

func main() {
	cfg := config.LoadConfig()
	dep := dependencies.InitDependencies(cfg)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register your routes
	routes.SetupReadRoutes(mux, dep)  // Assuming you have a function to set up read routes
	routes.SetupWriteRoutes(mux, dep) // And another for write routes

	const port = ":8080"
	fmt.Printf("Starting server at port %s\n", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
