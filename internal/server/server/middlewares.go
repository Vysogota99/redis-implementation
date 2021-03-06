package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *router) keyToStringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		type requestKey struct {
			Key interface{} `json:"key" binding:"required"`
		}
		rk := &requestKey{}

		body, err := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := json.Unmarshal(body, rk); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		if rk.Key == nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		var keyString string
		switch rk.Key.(type) {
		case float64:
			keyString = fmt.Sprintf("%f", rk.Key.(float64))
		case int64:
			keyString = fmt.Sprintf("%d", rk.Key.(int64))
		case map[string]interface{}:
			serializedData, err := json.Marshal(rk.Key)
			if err != nil {
				c.AbortWithError(http.StatusUnprocessableEntity, err)
				return
			}
			keyString = string(serializedData)
		default:
			keyString = rk.Key.(string)
		}

		c.Set("key", keyString)
		c.Next()
	}
}

// AuthUserMiddleware - ...
func (r *router) authUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := r.sessionStore.Get(c.Request, r.sessionName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err":     err.Error(),
				"meesage": "",
			})
			return
		}

		login, ok := session.Values["user_login"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"meesage": "Пользователь не авторизован",
			})
			return
		}

		userKey := fmt.Sprintf("user:%s", login)
		user, err := r.redis.GetHash(context.Background(), userKey)
		if len(user) == 0 {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Пользователь не найден"))
			return
		}

		c.Set("user", user)
		c.Next()

	}
}
