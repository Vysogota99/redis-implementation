package server

import (
	"log"
	"net/http"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func (r *router) setHashHandler(c *gin.Context) {
	keyChan, exists := c.Get("key")
	if !exists {
		respond(c, http.StatusInternalServerError, "", "No key chan in context")
		return
	}
	key := <-keyChan.(chan string)

	data := &models.SetHashRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusBadRequest, "", err.Error())
		return
	}

	err := r.redis.SetHash(c, key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, "success", "")
}

func (r *router) setStringHandler(c *gin.Context) {
	keyChan, exists := c.Get("key")
	if !exists {
		respond(c, http.StatusInternalServerError, "", "No key chan in context")
		return
	}
	key := <-keyChan.(chan string)

	data := &models.SetStringRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusUnprocessableEntity, "", err.Error())
		return
	}

	result, err := r.redis.SetString(c, key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}
	respond(c, http.StatusOK, result, "")
}

func (r *router) setListHandler(c *gin.Context) {
	keyChan, exists := c.Get("key")
	if !exists {
		respond(c, http.StatusInternalServerError, "", "No key chan in context")
		return
	}
	key := <-keyChan.(chan string)

	data := &models.SetListRequest{}
	if err := c.ShouldBindJSON(data); err != nil {
		respond(c, http.StatusUnprocessableEntity, "", err.Error())
		return
	}

	err := r.redis.SetList(c, key, data.Value, data.TTL)
	if err != nil {
		log.Println(err)
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, "success", "")
}

func (r *router) getHashHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, "", "No field key in get query")
	}

	result, err := r.redis.GetHash(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, "", err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) getStringHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, "", "No field key in get query")
	}

	result, err := r.redis.GetString(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, "", err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) keysHandler(c *gin.Context) {
	pattern := c.Query("pattern")
	if pattern == "" {
		respond(c, http.StatusBadRequest, "", "No field pattern in get query")
	}
	result, err := r.redis.Getkeys(c, pattern)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, "", err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, "", err.Error())
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) getListHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		respond(c, http.StatusBadRequest, "", "No field key in get query")
		return
	}

	result, err := r.redis.GetList(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, "", err.Error())
			return
		}

		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *router) deleteHandler(c *gin.Context) {
	keyChan, exists := c.Get("key")
	if !exists {
		respond(c, http.StatusInternalServerError, "", "No key chan in context")
		return
	}
	key := <-keyChan.(chan string)

	result, err := r.redis.Delete(c, key)
	if err != nil {
		if err == redis.Nil {
			respond(c, http.StatusNoContent, "", err.Error())
			return
		}

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
