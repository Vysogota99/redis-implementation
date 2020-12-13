package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Vysogota99/redis-implementation/internal/server/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
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
		return
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
		return
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
		return
	}
	result, err := r.redis.Getkeys(c, pattern)
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

func (r *router) LoginHadler(c *gin.Context) {
	req := models.User{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "", err.Error())
		return
	}

	userKey := fmt.Sprintf("user:%s", req.Login)
	user, err := r.redis.GetHash(context.Background(), userKey)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if len(user) == 0 {
		respond(c, http.StatusUnauthorized, "Пользователь с таким логином еще не зарегистрирован", "")
		return
	}

	if !comparePasswords(user["password"], req.Password) {
		respond(c, http.StatusUnauthorized, "Неправильный пароль", "")
		return
	}

	if err = login(&req, r.sessionStore, c.Writer, c.Request, r.sessionName); err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusCreated, req.Login, "")
}

func (r *router) SignupHandler(c *gin.Context) {
	req := models.User{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "", err.Error())
		return
	}

	userKey := fmt.Sprintf("user:%s", req.Login)
	user, err := r.redis.GetHash(context.Background(), userKey)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if len(user) != 0 {
		respond(c, http.StatusBadRequest, "Пользователь с таким логином уже зарегистрирован", "")
		return

	}

	userCreate := map[string]interface{}{
		"login":    req.Login,
		"password": hashPassword(req.Password),
	}

	err = r.redis.SetHash(context.Background(), userKey, userCreate, 0)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	if err = login(&req, r.sessionStore, c.Writer, c.Request, r.sessionName); err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusCreated, userCreate, "")
}

func (r *router) LogoutHandler(c *gin.Context) {
	session, err := r.sessionStore.Get(c.Request, r.sessionName)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	session.Options.MaxAge = -1
	if err = sessions.Save(c.Request, c.Writer); err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, "loggedout", "")
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

func login(user *models.User, sessionStore sessions.Store, w http.ResponseWriter, r *http.Request, sessionName string) error {
	log.Println(sessionName)
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values["user_login"] = user.Login
	if err := sessionStore.Save(r, w, session); err != nil {
		return err
	}

	return nil
}

func hashPassword(value string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(value), 8)
	if err != nil {
		return ""
	}

	return string(result)
}

// ComparePasswords - ...
func comparePasswords(hash, password string) bool {
	result := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return result == nil
}
