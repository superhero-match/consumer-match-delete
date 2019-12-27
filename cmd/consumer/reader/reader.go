package reader

import (
	"github.com/consumer-match-delete/internal/config"
	"github.com/consumer-match-delete/internal/consumer"
	"github.com/consumer-match-delete/internal/db"
)

// Reader holds all the data relevant.
type Reader struct {
	DB       *db.DB
	Consumer *consumer.Consumer
}

// NewReader configures Reader.
func NewReader(cfg *config.Config) (r *Reader, err error) {
	dbs, err := db.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	cs := consumer.NewConsumer(cfg)

	return &Reader{
		DB:       dbs,
		Consumer: cs,
	}, nil
}
