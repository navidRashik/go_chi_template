package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStreamBaseRunnerRepo interface {
	Run(messageID string, messageData map[string]interface{}, consumerGroupName string, redisClient *redis.Client) error
	RedisDataRepo
}
type RedisStreamBaseRunnerImpl struct {
	*RedisData
}

func (r *RedisStreamBaseRunnerImpl) Run(messageID string, messageData map[string]interface{}, consumerGroupName string, redisClient *redis.Client) error {
	r.messageID = messageID
	r.messageData = messageData
	r.consumerGroupName = consumerGroupName
	r.redisClient = redisClient

	ctx := context.Background()

	if err := r.operationStrategy.preRunOperation(ctx); err != nil {
		return err
	}

	if err := r.operationStrategy.RunOperation(ctx); err != nil {
		return err
	}

	return r.operationStrategy.postRunOperation(ctx)
}
