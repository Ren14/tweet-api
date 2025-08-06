package routes

import (
	"net/http"

	"github.com/renzonaitor/tweet-api/cmd/http/dependencies"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/reader"
)

func SetupReadRoutes(mux *http.ServeMux, dep dependencies.Dependencies) {
	readHandler := reader.NewHandler(dep.ReaderHandler.Timeline)
	mux.HandleFunc("/api/v1/timeline", readHandler.HandleGetTimeline)
}
