package api

import (
	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/EmilioCliff/e-commerce/db/worker"
	"github.com/gin-gonic/gin"
)

// Server struct
type Server struct {
	store       *db.Store
	router      *gin.Engine
	maker       token.Maker
	distributor *worker.RedisTaskDistributor
}

// Create a server
func NewServer(store *db.Store, maker token.Maker, distributor *worker.RedisTaskDistributor) *Server {
	server := Server{
		store:       store,
		maker:       maker,
		distributor: distributor,
	}

	server.setRoutes()
	return &server
}

// Set our routes and endpoints
func (server *Server) setRoutes() {
	r := gin.Default()

	// logger middleware and authentication middleware
	r.Use(loggerMiddleware())
	auth := r.Group("/").Use(authMiddleware(server.maker))

	auth.POST("/api/blogs/:id", server.addBlog)       // id == user_id creating blog
	auth.POST("/api/blogs/edit/:id", server.editBlog) // id == blog_id being edited
	r.GET("/api/blogs", server.listBlogs)
	auth.DELETE("api/blogs/delete/:id", server.deleteBlog)
	r.GET("/api/blogs/:id", server.getAdminsBlogs) // id == user_id to retrieve blogs of

	auth.POST("/api/orders", server.createOrder)

	server.router = r
}

// Start the gin engine
func (server *Server) Start(address string) error {
	return server.router.Run()
}

// Structure the errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
