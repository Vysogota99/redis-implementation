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
	redis        store.RedisImpl
	sessionName  string
	sessionStore sessions.Store
}

// newRouter - helper for initialization http
func newRouter(serverPort, sessionName string, redis store.RedisImpl, sessionStore sessions.Store) *router {
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
		list.GET("/lrange", r.lRangeHandler)
		list.POST("/lset", r.keyToStringMiddleware(), r.lSetHandler)
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
		hash.POST("/hset", r.keyToStringMiddleware(), r.hSetHandler)
		hash.GET("/hget", r.hGetHandler)
	}

	r.router.GET("/keys", r.keysHandler)
	r.router.POST("/del", r.keyToStringMiddleware(), r.deleteHandler)

	r.router.POST("/login", r.loginHadler)
	r.router.POST("/signup", r.signupHandler)
	r.router.POST("/logout", r.authUserMiddleware(), r.logoutHandler)

	return r.router
}
