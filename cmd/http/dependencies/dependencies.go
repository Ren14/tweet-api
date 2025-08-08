package dependencies

import (
	"fmt"

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
	postgresRepo, err := postgres.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect a postgres: %s", err.Error()))
	}
	redisRepo, err := redis.NewRepository(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect a redis: %s", err.Error()))
	}

	// service layer
	timelineService := timeline.NewService(postgresRepo, redisRepo)
	userService := user.NewService(postgresRepo, timelineService)

	// handler layer
	writerHandler := writer.NewHandler(userService)
	readerHandler := reader.NewHandler(timelineService)

	return Dependencies{
		WriterHandler: *writerHandler,
		ReaderHandler: *readerHandler,
	}
}
