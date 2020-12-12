package server

import (
	"github.com/Vysogota99/redis-implementation/internal/server/store"
	"github.com/gin-gonic/gin"
)

// router ...
type router struct {
	router     *gin.Engine
	serverPort string
	redis      *store.Redis
}

// newRouter - helper for initialization http
func newRouter(serverPort string, redis *store.Redis) *router {
	return &router{
		router:     gin.Default(),
		serverPort: serverPort,
		redis:      redis,
	}
}

// Setup - найстройка роутера
func (r *router) setup() *gin.Engine {
	list := r.router.Group("/list")
	{
		list.POST("/set", r.setListHandler)
		list.GET("/get", r.getListHandler)
	}

	str := r.router.Group("/string")
	{
		str.POST("/set", r.setStringHandler)
		str.GET("/get", r.getStringHandler)
	}

	hash := r.router.Group("/hash")
	{
		hash.POST("/set", r.setHashHandler)
		hash.GET("/get", r.getHashHandler)
	}

	r.router.GET("/keys", r.keysHandler)
	r.router.DELETE("/del", r.deleteHandler)
	return r.router
}
