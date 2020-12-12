package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
	"github.com/Vysogota99/redis-implementation/internal/server/store"
	"github.com/stretchr/testify/assert"
)

func TestSetlistHandler(t *testing.T) {
	redis, _ := store.New("localhost:6379")
	router := newRouter(":3000", redis)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	mapValue := make(map[string]interface{})
	mapValue["key"] = 1

	reqBody := models.SetListRequest{
		Key:   "key",
		Value: []interface{}{3, "123", mapValue},
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/list/set", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSetHashHandler(t *testing.T) {
	redis, _ := store.New("localhost:6379")
	router := newRouter(":3000", redis)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.SetHashRequest{
		Key: "user:ivan",
		Value: map[string]interface{}{
			"name": "Ivan",
			"age":  21,
		},
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/hash/set", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSetStringHandler(t *testing.T) {
	redis, err := store.New("localhost:6379")
	assert.NoError(t, err)

	router := newRouter(":3000", redis)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.SetStringRequest{
		Key:   "user:ivan",
		Value: "lapshin",
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/string/set", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetHashHandler(t *testing.T) {
	redis, err := store.New("localhost:6379")
	assert.NoError(t, err)

	router := newRouter(":3000", redis)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "user:ivan"
	resp, err := http.Get(fmt.Sprintf("%s/hash/get?key=%s", ts.URL, key))
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetListHandler(t *testing.T) {
	redis, err := store.New("localhost:6379")
	assert.NoError(t, err)

	router := newRouter(":3000", redis)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "tresh"
	resp, err := http.Get(fmt.Sprintf("%s/list/get?key=%s", ts.URL, key))
	assert.Equal(t, 200, resp.StatusCode)
}
