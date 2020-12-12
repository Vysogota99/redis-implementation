package server

import (
	"log"
	"net/http"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func (r *router) setHashHandler(c *gin.Context) {
	data := &models.SetHashRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	result, err := r.redis.SetHash(c, data.Key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) setStringHandler(c *gin.Context) {
	data := &models.SetStringRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	result, err := r.redis.SetString(c, data.Key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	respond(c, http.StatusOK, result, "")
}

func (r *router) setListHandler(c *gin.Context) {
	data := &models.SetListRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusUnprocessableEntity, nil, err.Error())
		return
	}

	result, err := r.redis.SetList(c, data.Key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) getHashHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, nil, "No field key in get query")
	}

	result, err := r.redis.GetHash(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, nil, err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) getStringHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, nil, "No field key in get query")
	}

	result, err := r.redis.GetString(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, nil, err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) keysHandler(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		respond(c, http.StatusBadRequest, nil, "No field pattern in get query")
	}
	result, err := r.redis.Getkeys(c, pattern)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, nil, err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, nil, err.Error())
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) getListHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, nil, "No field key in get query")
	}

	result, err := r.redis.GetList(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, nil, err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) deleteHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, nil, "No field key in get query")
	}

	result, err := r.redis.Delete(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, nil, err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, nil, err.Error())
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
