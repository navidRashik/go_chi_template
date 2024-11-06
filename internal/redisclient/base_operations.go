package redisclient

import (
	"context"

	"github.com/sirupsen/logrus"
)

type RedisStreamOperationRepo interface {
	preRunOperation(ctx context.Context) error
	postRunOperation(ctx context.Context) error
	RunOperation(ctx context.Context) error
	RedisDataRepo
}

type RedisStreamRunnerBaseOperationImpl struct {
	RedisDataRepo
}

func (r *RedisStreamRunnerBaseOperationImpl) RunOperation(ctx context.Context) error {
	// Implement your operation logic here
	logrus.Infof("Running operation for message ID: %s with data: %v",
		r.GetRedisData().messageID, r.GetRedisData().messageData)
	return nil
}

func (r *RedisStreamRunnerBaseOperationImpl) preRunOperation(ctx context.Context) error {
	logrus.Infof("Received a redis stream message with message_id: %s message_data: %v where stream_name: %s consumer_group_name: %s from preRunOperation",
		r.GetRedisData().messageID, r.GetRedisData().messageData, r.GetRedisData().streamName, r.GetRedisData().consumerGroupName)
	return nil
}

func (r *RedisStreamRunnerBaseOperationImpl) postRunOperation(ctx context.Context) error {
	logrus.Infof("Run method executed for a redis stream message with message_id: %s message_data: %v from postRunOperation",
		r.GetRedisData().messageID, r.GetRedisData().messageData)

	return r.GetRedisData().redisClient.XAck(ctx, r.GetRedisData().streamName, r.GetRedisData().consumerGroupName, r.GetRedisData().messageID).Err()
}
