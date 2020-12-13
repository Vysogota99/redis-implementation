package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	REDIS_SERVER_URL = "http://127.0.0.1:3000"
)

func (r *router) ListHandler(c *gin.Context) {
	type payload struct {
		Key   interface{}   `json:"key" binding:"required"`
		Value []interface{} `json:"value"`
		TTL   int           `json:"ttl"`
	}

	type request struct {
		Method  string   `json:"method" binding:"required"`
		Payload *payload `json:"payload"`
	}
	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	switch req.Method {
	case "set":
		if req.Payload == nil || req.Payload.Value == nil {
			respond(c, http.StatusBadRequest, nil, "Поле payload или/и value не задано")
			return
		}

		payload, err := json.Marshal(req.Payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, nil, err.Error())
			return
		}

		result, err := post("/list/set", payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")

	case "get":
		result, err := get("/list/get", "key", req.Payload.Key)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")
	default:
		respond(c, http.StatusBadRequest, nil, "Метод не найден, доступны: get/set")
		return
	}
}

func (r *router) StringHandler(c *gin.Context) {
	type payload struct {
		Key   interface{} `json:"key" binding:"required"`
		Value string      `json:"value"`
		TTL   int         `json:"ttl"`
	}

	type request struct {
		Method  string   `json:"method" binding:"required"`
		Payload *payload `json:"payload"`
	}
	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	switch req.Method {
	case "set":
		if req.Payload == nil || req.Payload.Value == "" {
			respond(c, http.StatusBadRequest, nil, "Поле payload или/и value не задано")
			return
		}

		payload, err := json.Marshal(req.Payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, nil, err.Error())
			return
		}

		result, err := post("/string/set", payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")

	case "get":
		result, err := get("/string/get", "key", req.Payload.Key)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")
	default:
		respond(c, http.StatusBadRequest, nil, "Метод не найден, доступны: get/set")
		return
	}
}

func (r *router) MapHandler(c *gin.Context) {
	type payload struct {
		Key   interface{}            `json:"key" binding:"required"`
		Value map[string]interface{} `json:"value"`
		TTL   int                    `json:"ttl"`
	}

	type request struct {
		Method  string   `json:"method" binding:"required"`
		Payload *payload `json:"payload"`
	}
	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	switch req.Method {
	case "set":
		if req.Payload == nil || req.Payload.Value == nil {
			respond(c, http.StatusBadRequest, nil, "Поле payload или/и value не задано")
			return
		}

		payload, err := json.Marshal(req.Payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, nil, err.Error())
			return
		}

		result, err := post("/hash/set", payload)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")

	case "get":
		result, err := get("/hash/get", "key", req.Payload.Key)
		if err != nil {
			respond(c, http.StatusInternalServerError, "", err.Error())
			return
		}

		respond(c, http.StatusOK, result, "")
	default:
		respond(c, http.StatusBadRequest, nil, "Метод не найден, доступны: get/set")
		return
	}
}

// DelKey ...
func (r *router) DelKey(c *gin.Context) {
	type request struct {
		Key interface{} `json:"key" binding:"required"`
	}

	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	data := request{
		Key: req.Key,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	result, err := post("/del", payload)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

// GetKeys ...
func (r *router) GetKeys(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		respond(c, http.StatusBadRequest, nil, "Поле pattern не задано")
		return
	}

	result, err := get("/keys", "pattern", pattern)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func respond(c *gin.Context, code int, result interface{}, err string) {
	if err == "EOF" {
		result = "Неправильное тело запроса"
	}

	c.JSON(
		code,
		gin.H{
			"result": result,
			"error":  err,
		},
	)
}

func post(route string, requestData []byte) (interface{}, error) {
	response, err := http.Post(fmt.Sprintf("%s%s", REDIS_SERVER_URL, route), "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		response.Body.Close()
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var serialized interface{}
		if err := json.Unmarshal(body, &serialized); err != nil {
			return nil, err
		}

		return serialized, nil
	}

	return nil, fmt.Errorf("Сервер redis вернул ошибку %d", response.StatusCode)
}

func get(route, variable string, key interface{}) (interface{}, error) {
	var keyString string
	switch key.(type) {
	case float64:
		keyString = fmt.Sprintf("%f", key.(float64))
	case int64:
		keyString = fmt.Sprintf("%d", key.(int64))
	case map[string]interface{}:
		serializedData, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		keyString = string(serializedData)
	default:
		keyString = key.(string)
	}

	response, err := http.Get(fmt.Sprintf("%s%s?%s=%s", REDIS_SERVER_URL, route, variable, keyString))
	if err != nil {
		response.Body.Close()
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var serialized interface{}
		if err := json.Unmarshal(body, &serialized); err != nil {
			return nil, err
		}

		return serialized, nil
	}

	return nil, fmt.Errorf("Сервер redis вернул ошибку %d", response.StatusCode)
}
