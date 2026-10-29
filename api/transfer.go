package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mayank-gulshan/Bank_System/db/sqlc"
)

type transferAccountRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"To_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) TransferAccount(ctx *gin.Context) {
	var req transferAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !server.ValidAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.ValidAccount(ctx, req.ToAccountID, req.Currency) {
		return

	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) ValidAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		err := fmt.Errorf("cuureny not match")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true

}
