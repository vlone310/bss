package httpsrv

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/vlone310/bss/internal/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	AmountCents   int64  `json:"amount_cents" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req createTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !s.validAccount(c, req.FromAccountID, req.Currency, -req.AmountCents) {
		return
	}

	if !s.validAccount(c, req.ToAccountID, req.Currency, req.AmountCents) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		AmountCents:   req.AmountCents,
	}

	transferResult, err := s.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, transferResult)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string, amountCents int64) bool {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(errAccountNotFound))
			return false
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account [%d] currency missmatch: %s and %s", account.ID, account.Currency, currency)))
		return false
	}

	if account.Balance+amountCents < 0 {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account [%d] balance not enough: %d", account.ID, account.Balance)))
		return false
	}

	return true
}
