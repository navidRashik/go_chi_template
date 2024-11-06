package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisStreamInitiatorRepo interface {
	ReadStream(ctx context.Context, runners []RedisStreamBaseRunnerRepo, consumerName string, redisClient *redis.Client) error
	processStream(ctx context.Context, consumerName string, redisClient *redis.Client, consumerGroupName string,
		streamKeyInstanceMap map[string]RedisStreamBaseRunnerRepo, streamListenersMap map[string]string)
	teardownStreamReader(ctx context.Context, consumerName string, redisClient *redis.Client,
		consumerGroupName string, streamListenersMap map[string]string)
	GetStreamListenerMap() map[string]string
	// WriteStream(ctx context.Context, runners []RedisStreamRunner, messageData map[string]interface{}, redisClient *redis.Client) error
}

type RedisStreamInitiatorImpl struct {
	streamListenerMap map[string]string
}

func (r *RedisStreamInitiatorImpl) GetStreamListenerMap() map[string]string {
	return r.streamListenerMap
}

func (r *RedisStreamInitiatorImpl) ReadStream(ctx context.Context, runners []RedisStreamBaseRunnerRepo, consumerName string, consumerGroupName string, redisClient *redis.Client) error {

	streamKeyInstanceMap, streamListenersMap, err := r.setupStreamReader(ctx, runners, redisClient, consumerGroupName, consumerName)
	r.streamListenerMap = streamListenersMap
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return r.teardownStreamReader(ctx, consumerName, redisClient, consumerGroupName, streamListenersMap)
		default:
			if err := r.processStream(ctx, consumerName, redisClient, consumerGroupName, streamKeyInstanceMap, streamListenersMap); err != nil {
				logrus.Warnf("Error processing stream: %v", err)
				time.Sleep(60 * time.Second)
				streamKeyInstanceMap, streamListenersMap, err = r.setupStreamReader(ctx, runners, redisClient, consumerGroupName, consumerName)
				if err != nil {
					return err
				}
			}
		}
	}
}

func (r *RedisStreamInitiatorImpl) processStream(ctx context.Context, consumerName string, redisClient *redis.Client, consumerGroupName string,
	streamKeyInstanceMap map[string]RedisStreamBaseRunnerRepo, streamListenersMap map[string]string) error {

	streams := make([]string, 0, len(streamListenersMap)*2)
	for k, v := range streamListenersMap {
		streams = append(streams, k, v)
	}

	xStreamSlice, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroupName,
		Consumer: consumerName,
		Streams:  streams,
		Count:    1,
		Block:    0,
	}).Result()

	if err != nil {
		return err
	}

	for _, xStream := range xStreamSlice {
		runnerInstance := streamKeyInstanceMap[xStream.Stream]
		if len(xStream.Messages) == 0 {
			logrus.Infof("-------- No more messages to read for %s ---------", runnerInstance.GetStreamName())
			continue
		}

		for _, msg := range xStream.Messages {
			messageData := make(map[string]interface{})
			for k, v := range msg.Values {
				messageData[k] = v
			}

			fmt.Printf("Received message with ID %s and data %v\n from processStream", msg.ID, messageData)

			if err := runnerInstance.Run(msg.ID, messageData, consumerGroupName, redisClient); err != nil {
				logrus.Errorf("Error running operation: %v", err)
			}
		}
	}

	return nil
}

func (r *RedisStreamInitiatorImpl) teardownStreamReader(ctx context.Context, consumerName string, redisClient *redis.Client,
	consumerGroupName string, streamListenersMap map[string]string) error {

	for streamName := range streamListenersMap {
		logrus.Infof("Removing consumer %v from group %v for stream %v", consumerName, consumerGroupName, streamName)
		if err := redisClient.XGroupDelConsumer(ctx, streamName, consumerGroupName, consumerName).Err(); err != nil {
			logrus.Errorf("Error removing consumer %v from group %v for stream %v: %v", consumerName, consumerGroupName, streamName, err)
			return err
		}
	}

	logrus.Info("Tearing down and removing consumer names from the group as this service is shutting down")
	return nil
}

func (r *RedisStreamInitiatorImpl) setupStreamReader(ctx context.Context, runners []RedisStreamBaseRunnerRepo, redisClient *redis.Client,
	consumerGroupName string, consumerName string) (map[string]RedisStreamBaseRunnerRepo, map[string]string, error) {

	streamKeyInstanceMap := make(map[string]RedisStreamBaseRunnerRepo)
	streamListenersMap := make(map[string]string)

	for _, runner := range runners {
		streamName := runner.GetStreamName()

		exists, err := redisClient.Exists(ctx, streamName).Result()
		if err != nil {
			return nil, nil, err
		}

		if exists == 0 {
			if err := redisClient.XGroupCreateMkStream(ctx, streamName, consumerGroupName, "$").Err(); err != nil {
				return nil, nil, err
			}
		} else {
			groups, err := redisClient.XInfoGroups(ctx, streamName).Result()
			if err != nil {
				return nil, nil, err
			}

			consumerGroupExists := false
			for _, group := range groups {
				if group.Name == consumerGroupName {
					consumerGroupExists = true
					break
				}
			}

			if !consumerGroupExists {
				if err := redisClient.XGroupCreate(ctx, streamName, consumerGroupName, "$").Err(); err != nil {
					return nil, nil, err
				}
			}
		}

		if err := redisClient.XGroupCreateConsumer(ctx, streamName, consumerGroupName, consumerName).Err(); err != nil {
			return nil, nil, err
		}

		streamKeyInstanceMap[streamName] = runner
		streamListenersMap[streamName] = ">"
	}

	return streamKeyInstanceMap, streamListenersMap, nil
}

type RedisStreamRunnerOperationExample struct {
	RedisStreamOperationRepo
}

func (r *RedisStreamRunnerOperationExample) RunOperation(ctx context.Context) error {
	logrus.Info("Your New implementation here for the operations--------->")

	// TODO: Testing write opearation so comment the following code of r.Write
	// r.Write(ctx, r.GetRedisData().messageData)
	return r.RedisStreamOperationRepo.RunOperation(ctx)
}

// func (r *RedisStreamRunnerOperationExample) Write(ctx context.Context, messageData map[string]interface{}) error {
// 	// Implement the logic to write the message to Redis
// 	// using the Redis client

// 	// Example code:
// 	logrus.Info("Writing message to Redis------------------->")
// 	err := r.GetRedisData().redisClient.XAdd(ctx, &redis.XAddArgs{
// 		Stream: r.GetStreamName(),
// 		Values: messageData,
// 	}).Err()
// 	if err != nil {
// 		logrus.Errorf("Error writing message to Redis: %v", err)
// 		return err
// 	}

// 	return nil
// }

func PrepareRedis(ctx context.Context, redisClient *redis.Client) (context.Context, *redis.Client, string, *RedisStreamInitiatorImpl, []RedisStreamBaseRunnerRepo, string, context.CancelFunc) {

	consumerGroupName := "workflow_automation"

	redisData := &RedisData{
		streamName:        "example_stream",
		redisClient:       redisClient,
		consumerGroupName: consumerGroupName,
	}

	redisStreamInitiator := &RedisStreamInitiatorImpl{
		streamListenerMap: make(map[string]string),
	}

	runners := []RedisStreamBaseRunnerRepo{
		&RedisStreamBaseRunnerImpl{
			RedisData: redisData,
			op:        &RedisStreamRunnerOperationExample{&RedisStreamRunnerBaseOperationImpl{redisData}},
		},
	}

	consumerName := uuid.New().String()

	ctx, cancel := context.WithCancel(ctx)
	return ctx, redisClient, consumerGroupName, redisStreamInitiator, runners, consumerName, cancel
}
