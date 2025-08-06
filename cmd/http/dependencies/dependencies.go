package dependencies

import (
	"github.com/renzonaitor/tweet-api/cmd/http/config"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/reader"
	"github.com/renzonaitor/tweet-api/cmd/http/handlers/writer"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/postgres"
	"github.com/renzonaitor/tweet-api/internal/infraestructure/redis"
	"github.com/renzonaitor/tweet-api/internal/service/timeline"
	"github.com/renzonaitor/tweet-api/internal/service/user"
)

type Dependencies struct {
	WriterHandler writer.WriterHandler
	ReaderHandler reader.ReaderHandler
}

func InitDependencies(cfg config.Config) Dependencies {

	// repository layer
	postgresRepo := postgres.NewRepository(cfg)
	redisRepo := redis.NewRepository(cfg)

	// service layer
	userService := user.NewService(postgresRepo)
	timelineService := timeline.NewService(postgresRepo, redisRepo)

	// handler layer
	writerHandler := writer.NewHandler(userService)
	readerHandler := reader.NewHandler(timelineService)

	return Dependencies{
		WriterHandler: *writerHandler,
		ReaderHandler: *readerHandler,
	}
}
