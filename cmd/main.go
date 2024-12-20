package main

import (
	books "github.com/mdshahjahanmiah/distributed-caching/pkg/book"
	"github.com/mdshahjahanmiah/distributed-caching/pkg/cache"
	"github.com/mdshahjahanmiah/distributed-caching/pkg/config"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"log/slog"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "err", err)
		return
	}

	logger, err := logging.NewLogger(conf.LoggerConfig)
	if err != nil {
		slog.Error("initializing logger", "err", err)
		return
	}

	// Initialize Redis nodes
	cache.InitializeNodes()
	defer cache.Cleanup()

	bookService := books.NewBookService(logger)

	// Fetch book details
	userRegion := "Europe"
	bookKey := "book::1984"
	bookDetails, err := bookService.FetchBookDetails(bookKey, userRegion)
	if err != nil {
		logger.Fatal("Failed to fetch book details", "err", err)
	}

	logger.Info("Fetched book details", "book_details", bookDetails)
}
