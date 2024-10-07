package kafkaclient

import (
	"sync"

	"example_project/internal/config"
	"example_project/internal/database"
	"example_project/internal/leveledlog"
)

type Kafka struct {
	Config config.Config
	Db     *database.DB
	Logger *leveledlog.LogStruct
	Wg     *sync.WaitGroup
}
