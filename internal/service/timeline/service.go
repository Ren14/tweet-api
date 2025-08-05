package timeline

type StorageRepo interface {
	// todo define contract
}

type CacheRepo interface {
	// todo define contract
}

type Service struct {
	Storage StorageRepo
	Cache   CacheRepo
}

func NewService(storage StorageRepo, cache CacheRepo) *Service {
	return &Service{
		Storage: storage,
		Cache:   cache,
	}
}
