package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vlone310/bss/internal/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	r.POST("/users", server.createUser)

	// adding routes
	r.POST("/accounts", server.createAccount)
	r.GET("/accounts/:id", server.getAccountByID)
	r.GET("/accounts", server.listAccounts)

	r.POST("/transfers", server.createTransfer)
	server.router = r
	return server
}

func (s *Server) ServeHTTP(addr string) error {
	fmt.Println("Server is running on", addr)
	return s.router.Run(addr)
}
