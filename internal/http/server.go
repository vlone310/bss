package httpsrv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/vlone310/bss/internal/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	r := gin.Default()

	// adding routes
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccountByID)
	r.GET("/accounts", server.listAccounts)
	// r.Get("/accounts", server.listAccounts)
	server.router = r
	return server
}

func (s *Server) ServeHTTP(addr string) error {
	fmt.Println("Server is running on", addr)
	return s.router.Run(addr)
}
