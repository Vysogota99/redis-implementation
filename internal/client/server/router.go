package server

import "github.com/gin-gonic/gin"

// router ...
type router struct {
	router     *gin.Engine
	serverPort string
}

// newRouter - helper for initialization http
func newRouter(serverPort string) *router {
	return &router{
		router:     gin.Default(),
		serverPort: serverPort,
	}
}

// Setup - найстройка роутера
func (r *router) setup() *gin.Engine {
	r.router.POST("/list", r.ListHandler)
	r.router.POST("/string", r.StringHandler)
	r.router.POST("/map", r.MapHandler)
	r.router.GET("/keys", r.GetKeys)
	r.router.POST("/delete", r.DelKey)

	return r.router
}
