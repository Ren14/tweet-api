package routes

import (
	"net/http"

	"github.com/renzonaitor/tweet-api/cmd/http/dependencies"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/writer"
)

func SetupWriteRoutes(mux *http.ServeMux, dep dependencies.Dependencies) {
	writerHandler := writer.NewHandler(dep.WriterHandler.UserService)
	mux.HandleFunc("/api/v1/tweet", writerHandler.HandlePublishTweet)
	mux.HandleFunc("/api/v1/follow", writerHandler.HandleFollowUser)
	mux.HandleFunc("/ping", writerHandler.Ping)
}
