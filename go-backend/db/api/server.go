package api

import (
	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
)

// Server struct
type Server struct {
	store  *db.Store
	router *gin.Engine
	maker  token.Maker
}

// Create a server
func NewServer(store *db.Store, maker token.Maker) *Server {
	server := Server{
		store: store,
		maker: maker,
	}

	server.setRoutes()
	return &server
}

// Set our routes and endpoints
func (server *Server) setRoutes() {
	r := gin.Default()

	auth := r.Group("/").Use(authMiddleware(server.maker))

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
