package redisclient

import (
	"sync"

	"example_project/internal/config"
	"example_project/internal/database"
	"example_project/internal/leveledlog"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Config                 config.Config
	Db                     *database.DB
	Logger                 *leveledlog.LogStruct
	Wg                     *sync.WaitGroup
	DistributedRedisClient *redis.Client
}

// type RedisStreamRunner struct {
// 	streamName        string
// 	consumerGroupName string
// 	redis             *redis.Client
// 	messageID         string
// 	messageData       map[string]string
// 	redisClient       *redis.Client
// }
