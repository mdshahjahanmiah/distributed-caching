package book

import (
	"github.com/mdshahjahanmiah/distributed-caching/pkg/cache"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"time"
)

type service struct {
	logger *logging.Logger
}

type Service interface {
	FetchBookDetails(key string, userRegion string) (string, error)
}

func NewBookService(logger *logging.Logger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) FetchBookDetails(key string, userRegion string) (string, error) {
	// Use consistent hashing with consideration for userRegion
	client := cache.GetClientForKey(key, userRegion)

	// Try to get the data from cache
	val, err := client.Get(cache.Context(), key).Result()
	if err == nil {
		s.logger.Info("Cache hit", "key", key, "client", client.Options().Addr)
		return val, nil
	}

	// Cache miss, simulate database fetch
	s.logger.Info("Cache miss", "key", key, "client", client.Options().Addr)
	bookDetails := "Book: 1984, Author: Miah Md Shahjahan, Genre: Programming, Price: $18.99"

	// Define TTL based on data freshness requirements
	ttl := 10 * time.Minute // Example: 10 minutes, adjust as needed

	err = client.Set(cache.Context(), key, bookDetails, ttl).Err()
	if err != nil {
		s.logger.Error("Error setting cache", "key", key, "err", err)
		return "", err
	}

	return bookDetails, nil

}
