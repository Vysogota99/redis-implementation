package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
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
func (r *Redis) SetHash(ctx context.Context, key string, value map[string]interface{}, ttl int) error {
	txf := func(tx *redis.Tx) error {
		_, err := r.client.HMSet(ctx, key, value).Result()
		if err != nil {
			return err
		}

		if ttl != 0 {
			_, err = r.client.Expire(ctx, key, time.Duration(ttl)*time.Minute).Result()
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := r.client.Watch(ctx, txf)
	if err != nil {
		return err
	}
	return nil
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
func (r *Redis) SetList(ctx context.Context, key string, value []interface{}, ttl int) error {
	txf := func(tx *redis.Tx) error {
		strSlice := make([]string, len(value))
		for i, val := range value {
			element := models.ListElement{}
			switch val.(type) {
			case float64:
				element.Dtype = "float64"
				element.Data = fmt.Sprintf("%f", val.(float64))
				serialized, err := json.Marshal(element)
				if err != nil {
					return err
				}
				strSlice[i] = string(serialized)
			case int64:
				element.Dtype = "int64"
				element.Data = fmt.Sprintf("%d", val.(int64))
				serialized, err := json.Marshal(element)
				if err != nil {
					return err
				}
				strSlice[i] = string(serialized)
			case map[string]interface{}:
				serializedData, err := json.Marshal(val)
				if err != nil {
					return err
				}

				element.Dtype = "map"
				element.Data = string(serializedData)
				serialized, err := json.Marshal(element)
				if err != nil {
					return err
				}
				strSlice[i] = string(serialized)
			default:
				element.Dtype = "string"
				element.Data = val.(string)
				serialized, err := json.Marshal(element)
				if err != nil {
					return err
				}
				strSlice[i] = string(serialized)
			}
		}

		_, err := r.client.RPush(ctx, key, strSlice).Result()
		if err != nil {
			return err
		}

		if ttl != 0 {
			_, err = r.client.Expire(ctx, key, time.Duration(ttl)*time.Minute).Result()
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := r.client.Watch(ctx, txf)
	if err != nil {
		return err
	}
	return nil
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
func (r *Redis) GetList(ctx context.Context, key string) ([]interface{}, error) {
	res, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	resultSlice := make([]interface{}, len(res))
	for i, el := range res {
		var lElement models.ListElement
		if err := json.Unmarshal([]byte(el), &lElement); err != nil {
			return nil, err
		}

		switch lElement.Dtype {
		case "float64":
			floatVal, err := strconv.ParseFloat(lElement.Data, 64)
			if err != nil {
				return nil, err
			}

			resultSlice[i] = floatVal
		case "int64":
			intVal, err := strconv.ParseInt(lElement.Data, 10, 64)
			if err != nil {
				return nil, err
			}

			resultSlice[i] = intVal
		case "string":
			resultSlice[i] = lElement.Data
		case "map":
			var deserializedValue interface{}
			if err := json.Unmarshal([]byte(lElement.Data), &deserializedValue); err != nil {
				return nil, err
			}

			resultSlice[i] = deserializedValue
		}
	}

	return resultSlice, nil
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
