package store

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestSetList(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "ivan"
	valuesExp := []interface{}{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}
	values := []interface{}{
		"1", "2", "qwerty",
	}

	mock.ExpectRPush(key, valuesExp...).SetVal(3)
	err := client.SetList(context.Background(), key, values, 0)
	assert.NoError(t, err)
}

func TestSetHash(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "ivan"
	valuesExp := map[string]interface{}{}
	valuesExp["name"] = "Ivan"
	valuesExp["age"] = 21

	values := map[string]interface{}{}
	values["name"] = "Ivan"
	values["age"] = 21

	mock.ExpectHMSet(key, valuesExp).SetVal(true)
	err := client.SetHash(context.Background(), key, values, 0)
	assert.NoError(t, err)
}

func TestSetString(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "ivan"
	valuesExp := "Lapshin"

	mock.ExpectSet(key, valuesExp, 0).SetVal("OK")
	_, err := client.SetString(context.Background(), key, valuesExp, 0)
	assert.NoError(t, err)
}

func TestGetKeys(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	pattern := "*"
	mock.ExpectKeys(pattern).SetVal([]string{"one", "two", "three"})
	_, err := client.Getkeys(context.Background(), pattern)
	assert.NoError(t, err)
}

func TestDel(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "Ivan"
	mock.ExpectDel(key).SetVal(1)
	_, err := client.Delete(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetString(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "Ivan"
	value := "Lapshin"
	mock.ExpectGet(key).SetVal(value)
	_, err := client.GetString(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetHash(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	value := map[string]string{
		"name":     "Ivan",
		"lastname": "lapshin",
	}
	mock.ExpectHGetAll(key).SetVal(value)
	_, err := client.GetHash(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetList(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	values := []string{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}

	mock.ExpectLRange(key, 0, -1).SetVal(values)
	_, err := client.GetList(context.Background(), key)
	assert.NoError(t, err)
}
