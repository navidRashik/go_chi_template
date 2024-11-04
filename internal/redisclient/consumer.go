package redisclient

import (
	"context"
	"example_project/internal/leveledlog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func (r *Redis) ConsumerEvent(doneChan chan bool, quit chan os.Signal) {
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	defer func() {
		doneChan <- true
	}()
	opt, err := redis.ParseURL(r.Config.DistributedRedisUrl)
	if err != nil {
		r.Logger.Error("failed to parse redis url, details: %s", err.Error())
		return
	}
	r.DistributedRedisClient = redis.NewClient(opt)
	// Handle timeout error
	// Handle other errors
	_ = testConnection(err, r, logger)

	ctx := context.Background()

	ctx, redisClient, consumerGroupName, redisStreamInitiator, runners, consumerName, cancel := PrepareRedis(ctx, r.DistributedRedisClient)
	go func() {
		if err := redisStreamInitiator.ReadStream(ctx, runners, consumerName, consumerGroupName, redisClient); err != nil {
			logrus.Fatalf("Error reading stream: %v", err)
		}
	}()

	// Catch SIGINT and SIGTERM signals
	// Wait for a signal
	<-quit

	// Cancel the context to stop the ReadStream goroutine

	// Perform a clean shutdown
	if err := redisStreamInitiator.teardownStreamReader(ctx, consumerName, redisClient, consumerGroupName, redisStreamInitiator.GetStreamListenerMap()); err != nil {
		logrus.Errorf("Error tearing down stream reader: %v", err)
	}
	cancel()

	logrus.Info("Shutdown complete")

}

func testConnection(err error, r *Redis, logger *leveledlog.LogStruct) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = r.DistributedRedisClient.Ping(ctx).Result()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			logger.Fatal("failed to connect to redis, details: %s", err.Error())

		} else {

			logger.Fatal("failed to connect to redis, details: %s", err.Error())
		}
	}
	logger.Debug("connected to redis")
	return ctx
}