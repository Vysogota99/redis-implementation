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

func TestSetListHandler(t *testing.T) {
	redis := store.NewMock()
	router := newRouter(":3000", "auth", redis, nil)

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
	resp.Body.Close()
}

func TestSetHashHandler(t *testing.T) {
	redis := store.NewMock()
	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.SetHashRequest{
		Key: "user:1",
		Value: map[string]interface{}{
			"name": "Ivan",
			"age":  21,
		},
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/hash/set", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestSetStringHandler(t *testing.T) {
	redis := store.NewMock()
	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.SetStringRequest{
		Key:   "user:1",
		Value: "lapshin",
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/string/set", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestGetHashHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "user:1"
	resp, _ := http.Get(fmt.Sprintf("%s/hash/get?key=%s", ts.URL, key))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestGetListHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "tresh"
	resp, _ := http.Get(fmt.Sprintf("%s/list/get?key=%s", ts.URL, key))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestSignUP(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.User{
		Login:    "user:ivan2",
		Password: "qwerty",
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/signup", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()
}

func TestLoginHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.User{
		Login:    "Ivan",
		Password: "lapshin",
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/login", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

func TestHGetHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "user:1"
	field := "name"

	resp, _ := http.Get(fmt.Sprintf("%s/hash/hget?key=%s&&field=%s", ts.URL, key, field))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestHSetHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := models.SetHashRequest{
		Key: "user:1",
		Value: map[string]interface{}{
			"role": "admin",
		},
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/hash/hset", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestLRangeHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "user:1"
	start := "0"
	stop := "1"

	resp, _ := http.Get(fmt.Sprintf("%s/list/lrange?key=%s&&start=%s&&stop=%s", ts.URL, key, start, stop))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestLSetHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := map[string]interface{}{
		"key":   "user:1",
		"value": -1,
		"index": 1,
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/list/lset", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestKeysHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	pattern := "*"
	resp, _ := http.Get(fmt.Sprintf("%s/keys?pattern=%s", ts.URL, pattern))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestGetStringHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	key := "user:1"
	resp, _ := http.Get(fmt.Sprintf("%s/string/get?key=%s", ts.URL, key))
	assert.Equal(t, 200, resp.StatusCode)
	resp.Body.Close()
}

func TestDeleteHandler(t *testing.T) {
	redis := store.NewMock()

	router := newRouter(":3000", "auth", redis, nil)

	ts := httptest.NewServer(router.setup())
	defer ts.Close()

	reqBody := map[string]interface{}{
		"key":   "user:1",
	}

	data, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	resp, err := http.Post(fmt.Sprintf("%s/del", ts.URL), "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
