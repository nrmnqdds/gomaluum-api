package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/nrmnqdds/gomaluum-api/helpers"
	"github.com/redis/go-redis/v9"
)

// A list of task types.
const (
	TypeRedisSave = "redis:save"
	TypeRedisGet  = "redis:get"
)

// Function to save user password to redis.
func SaveToRedis(ctx context.Context, t *asynq.Task) error {
	logger, _ := helpers.NewLogger()

	// opt, _ := redis.ParseURL(helpers.GetEnv("REDIS_URL"))
	opt := &redis.Options{
		Addr:     helpers.GetEnv("REDIS_URL"),
		Username: helpers.GetEnv("REDIS_USERNAME"),
		Password: helpers.GetEnv("REDIS_PASSWORD"),
	}
	client := redis.NewClient(opt)

	var p RedisSavePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logger.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		return err
	}

	err := client.Set(ctx, p.Key, p.Value, 0).Err()
	if err != nil {
		logger.Warnf("Failed to set user password to redis: %v", err)
		return err
	}

	return nil
}

type RedisSavePayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewSaveToRedisTask creates a new asynq.Task to save a key-value pair to redis.
func NewSaveToRedisTask(key, value string) (*asynq.Task, error) {
	payload, err := json.Marshal(RedisSavePayload{Key: key, Value: value})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRedisSave, payload), nil
}
