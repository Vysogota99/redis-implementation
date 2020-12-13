package server

import (
	"github.com/Vysogota99/redis-implementation/internal/server/store"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// router ...
type router struct {
	router       *gin.Engine
	serverPort   string
	redis        *store.Redis
	sessionName  string
	sessionStore sessions.Store
}

// newRouter - helper for initialization http
func newRouter(serverPort, sessionName string, redis *store.Redis, sessionStore sessions.Store) *router {
	return &router{
		router:       gin.Default(),
		serverPort:   serverPort,
		redis:        redis,
		sessionName:  sessionName,
		sessionStore: sessionStore,
	}
}

// Setup - найстройка роутера
func (r *router) setup() *gin.Engine {
	list := r.router.Group("/list")
	{
		list.POST("/set", r.keyToStringMiddleware(), r.setListHandler)
		list.GET("/get", r.getListHandler)
	}

	str := r.router.Group("/string")
	{
		str.POST("/set", r.keyToStringMiddleware(), r.setStringHandler)
		str.GET("/get", r.getStringHandler)
	}

	hash := r.router.Group("/hash")
	{
		hash.POST("/set", r.keyToStringMiddleware(), r.setHashHandler)
		hash.GET("/get", r.getHashHandler)
	}

	r.router.GET("/keys", r.keysHandler)
	r.router.POST("/del", r.keyToStringMiddleware(), r.deleteHandler)

	r.router.POST("/login", r.LoginHadler)
	r.router.POST("/signup", r.SignupHandler)
	r.router.POST("/logout", r.AuthUserMiddleware(), r.LogoutHandler)

	return r.router
}
