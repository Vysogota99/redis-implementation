package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
	"github.com/go-redis/redismock/v8"
)

// RedisMock ...
type RedisMock struct {
	client *Redis
	mock   redismock.ClientMock
}

// NewMock - helper to init redis mock
func NewMock() *RedisMock {
	db, mock := redismock.NewClientMock()

	client := &Redis{
		client: db,
	}
	return &RedisMock{
		client: client,
		mock:   mock,
	}
}

// Save ...
func (r *RedisMock) Save(ctx context.Context) error {
	return nil
}

// SetHash ...
func (r *RedisMock) SetHash(ctx context.Context, key string, value map[string]interface{}, ttl int) error {
	r.mock.ExpectHMSet(key, value).SetVal(true)
	err := r.client.SetHash(ctx, key, value, ttl)
	return err
}

// SetString ...
func (r *RedisMock) SetString(ctx context.Context, key, value string, ttl int) (string, error) {
	r.mock.ExpectSet(key, value, time.Duration(ttl)).SetVal("OK")
	res, err := r.client.SetString(ctx, key, value, ttl)
	if err != nil {
		return "", err
	}

	return res, nil
}

// SetList ...
func (r *RedisMock) SetList(ctx context.Context, key string, value []interface{}, ttl int) error {
	strSlice := make([]string, len(value))
	for i, val := range value {
		element := models.ListElement{}
		switch val.(type) {
		case float64:
			element.Dtype = "float64"
			element.Data = fmt.Sprintf("%.5f", val.(float64))
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
		case int:
			element.Dtype = "int"
			element.Data = fmt.Sprintf("%d", val.(int))
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

	r.mock.ExpectRPush(key, strSlice).SetVal(int64(len(value)))
	err := r.client.SetList(ctx, key, value, ttl)
	return err
}

// GetHash ...
func (r *RedisMock) GetHash(ctx context.Context, key string) (map[string]string, error) {
	value := map[string]string{
		"name":     "Ivan",
		"password": "$2a$08$$2a$08$oLVqr66wg50xQyk3Sa7zZeArxdeevVaz/Axbf6U.uwUAFZ0tdG6aW",
	}

	r.mock.ExpectHGetAll(key).SetVal(value)
	res, err := r.client.GetHash(ctx, key)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetString ...
func (r *RedisMock) GetString(ctx context.Context, key string) (string, error) {
	value := "value"
	r.mock.ExpectGet(key).SetVal(value)
	res, err := r.client.GetString(ctx, key)
	if err != nil {
		return "", err
	}

	return res, nil
}

// GetList ...
func (r *RedisMock) GetList(ctx context.Context, key string) ([]interface{}, error) {
	values := []string{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}
	r.mock.ExpectLRange(key, 0, -1).SetVal(values)
	res, err := r.client.GetList(ctx, key)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetKeys ...
func (r *RedisMock) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	r.mock.ExpectKeys(pattern).SetVal([]string{"one", "two", "three"})
	res, err := r.client.GetKeys(ctx, pattern)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete ...
func (r *RedisMock) Delete(ctx context.Context, key string) (int64, error) {
	r.mock.ExpectDel(key).SetVal(1)
	res, err := r.client.Delete(ctx, key)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// HGet ...
func (r *RedisMock) HGet(ctx context.Context, key string, field string) (string, error) {
	r.mock.ExpectHGet(key, field).SetVal("user:1")
	res, err := r.client.HGet(ctx, key, field)
	if err != nil {
		return "", err
	}

	return res, nil
}

// HSet ...
func (r *RedisMock) HSet(ctx context.Context, key string, values map[string]interface{}) (int64, error) {
	valueExp := map[string]interface{}{
		"role": "admin",
	}
	r.mock.ExpectHSet(key, valueExp).SetVal(2)
	res, err := r.client.HSet(ctx, key, values)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// LRange ...
func (r *RedisMock) LRange(ctx context.Context, key string, start, stop int64) ([]interface{}, error) {
	values := []string{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}

	r.mock.ExpectLRange(key, start, stop).SetVal(values)
	res, err := r.client.LRange(ctx, key, start, stop)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// LSet ...
func (r *RedisMock) LSet(ctx context.Context, key string, index int64, value interface{}) (string, error) {
	r.mock.ExpectLSet(key, int64(index), value)
	res, err := r.client.LSet(ctx, key, index, value)
	if err != nil {
		return "", err
	}

	return res, nil
}
