package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/vlone310/bss/internal/db/sqlc"
	"github.com/vlone310/bss/util"
)

var errUserExists = errors.New("user already exists")

type createUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := s.store.CreateUser(c, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				c.JSON(http.StatusForbidden, errorResponse(errUserExists))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.Time.UTC(),
		CreatedAt:         user.CreatedAt.Time.UTC(),
	}

	c.JSON(http.StatusCreated, res)
}
