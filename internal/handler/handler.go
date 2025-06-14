package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-mcp/pkg/jwt"
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(
	logger *log.Logger,
) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.MyCustomClaims).UserId
}
