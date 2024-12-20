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
	client := cache.GetClientForKey(userRegion)

	// Try to get the data from cache
	val, err := client.Get(cache.Context(), key).Result()
	if err == nil {
		s.logger.Error("Cache hit", "key", key, "user_region", userRegion)
		return val, nil
	}

	// Cache miss, simulate database fetch
	s.logger.Info("Cache miss", "key", key, "user_region", userRegion)
	bookDetails := "Book: 1984, Author: Miah Md Shahjahan, Genre: Programming, Price: $18.99"

	// Store in cache for future requests
	err = client.Set(cache.Context(), key, bookDetails, 10*time.Minute).Err()
	if err != nil {
		s.logger.Error("Error setting cache", "key", key, "err", err)
	}
	return bookDetails, nil
}
