package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis ...
type Redis struct {
	client *redis.Client
}

// New - helper to init redis
func New(addr string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
	}, nil
}

// SetHash ...
func (r *Redis) SetHash(ctx context.Context, key string, value map[string]interface{}, ttl int) (bool, error) {
	result, err := r.client.HMSet(ctx, key, value).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

// SetString ...
func (r *Redis) SetString(ctx context.Context, key, value string, ttl int) (string, error) {
	result, err := r.client.Set(ctx, key, value, time.Duration(ttl)*time.Minute).Result()
	if err != nil {
		return "", err
	}

	return result, nil
}

// SetList ...
func (r *Redis) SetList(ctx context.Context, key string, value []interface{}, ttl int) (int64, error) {
	strSlice := make([]string, len(value))
	for i, val := range value {
		switch val.(type) {
		case float64:
			strSlice[i] = fmt.Sprintf("%f", val.(float64))
		case int64:
			strSlice[i] = fmt.Sprintf("%d", val.(int64))
		case map[string]interface{}:
			serializedValue, err := json.Marshal(val)
			if err != nil {
				return 0, err
			}

			strSlice[i] = string(serializedValue)
		default:
			strSlice[i] = val.(string)
		}
	}

	result, err := r.client.RPush(ctx, key, strSlice).Result()
	if err != nil {
		return 0, err
	}

	return result, nil
}

// GetHash ...
func (r *Redis) GetHash(ctx context.Context, key string) (map[string]string, error) {
	res, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetString ...
func (r *Redis) GetString(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return res, nil
}

// GetList ...
func (r *Redis) GetList(ctx context.Context, key string) ([]string, error) {
	res, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Getkeys ...
func (r *Redis) Getkeys(ctx context.Context, pattern string) ([]string, error) {
	res, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete ...
func (r *Redis) Delete(ctx context.Context, key string) (int64, error) {
	res, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return res, nil
} 