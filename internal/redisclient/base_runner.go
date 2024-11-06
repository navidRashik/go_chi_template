package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisStreamBaseRunnerRepo interface {
	Run(messageID string, messageData map[string]interface{}, consumerGroupName string, redisClient *redis.Client) error

	Write(ctx context.Context, messageData map[string]interface{}) error
	RedisDataRepo
}
type RedisStreamBaseRunnerImpl struct {
	*RedisData
	op RedisStreamOperationRepo
}

func (r *RedisStreamBaseRunnerImpl) Run(messageID string, messageData map[string]interface{}, consumerGroupName string, redisClient *redis.Client) error {
	r.messageID = messageID
	r.messageData = messageData
	r.consumerGroupName = consumerGroupName
	r.redisClient = redisClient

	ctx := context.Background()

	if err := r.op.preRunOperation(ctx); err != nil {
		return err
	}

	if err := r.op.RunOperation(ctx); err != nil {
		return err
	}

	return r.op.postRunOperation(ctx)
}

func (r *RedisStreamBaseRunnerImpl) Write(ctx context.Context, messageData map[string]interface{}) error {
	// Implement the logic to write the message to Redis
	// using the Redis client

	// Example code:
	logrus.Info("Writing message to Redis------------------->")
	streamName := r.GetStreamName()
	err := r.GetRedisData().redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: messageData,
	}).Err()
	if err != nil {
		return err
	}

	return nil
}
