package service

import (
	"context"
	"crudly/model"
	"crudly/util/result"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRateLimitStore struct {
	redisClient *redis.Client
}

func NewRedisRateLimiterStore(redisClient *redis.Client) redisRateLimitStore {
	return redisRateLimitStore{
		redisClient,
	}
}

func (r *redisRateLimitStore) IncrementCallCount(projectId model.ProjectId, ttl time.Duration) result.R[uint] {
	redisKey := getRedisKey(projectId)

	val, err := r.redisClient.IncrBy(context.Background(), redisKey, 1).Result()
	if err != nil {
		return result.Errf[uint]("error incrementing redis key: %w", err)
	}

	if val == 1 {
		err = r.redisClient.Expire(context.Background(), redisKey, ttl).Err()
		if err != nil {
			return result.Errf[uint]("error setting ttl on redis key: %w", err)
		}
	}

	return result.Ok(uint(val))
}

func (r *redisRateLimitStore) GetCurrentCallCount(projectId model.ProjectId) result.R[uint] {
	val, err := r.redisClient.Get(context.Background(), getRedisKey(projectId)).Result()

	if err != nil {
		return result.Errf[uint]("error getting redis key: %w", err)
	}

	integer, err := strconv.Atoi(val)

	if err != nil {
		return result.Errf[uint]("couldn't parse redis key to int: %w", err)
	}

	return result.Ok(uint(integer))
}

func getRedisKey(projectId model.ProjectId) string {
	return projectId.String()
}
