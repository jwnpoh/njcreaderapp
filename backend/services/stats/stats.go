package stats

import (
	"fmt"

	"github.com/jwnpoh/njcreaderapp/backend/internal/core"
	"github.com/jwnpoh/njcreaderapp/backend/services/serializer"
)

type StatsDatabase interface {
	GetStats() (*core.Stats, error)
}

type Stats struct {
	db StatsDatabase
}

// NewStatsService returns an statsService object to implement methods to interact with PlanetScale database.
func NewStatsService(statsDB StatsDatabase) *Stats {
	return &Stats{statsDB}
}

// Stats gets a bunch of predefined statistics on the articles database.
func (stats *Stats) GetStats() (serializer.Serializer, error) {
	// get stats from db
	data, err := stats.db.GetStats()
	if err != nil {
		return nil, fmt.Errorf("unable to get articles from db - %w", err)
	}

	fmt.Println(data)
	return serializer.NewSerializer(false, "got stats", data), nil
}
