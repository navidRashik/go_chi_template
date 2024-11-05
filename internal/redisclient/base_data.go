package redisclient

import "github.com/redis/go-redis/v9"

type RedisDataRepo interface {
	GetStreamName() string
	GetRedisData() *RedisData
}
type RedisData struct {
	streamName        string
	messageID         string
	messageData       map[string]interface{}
	consumerGroupName string
	consumerName      string
	redisClient       *redis.Client
	operationStrategy RedisStreamOperationRepo
}

func (r *RedisData) GetStreamName() string {
	return r.streamName
}

func (r *RedisData) GetRedisData() *RedisData {
	return r
}
