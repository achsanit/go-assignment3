package handler

import (
	"net/http"

	"github.com/achsanit/go-assignment2/main/service"
	"github.com/gin-gonic/gin"
)

type PublicHandler interface {
	GeneratePublicAccessToken(ctx *gin.Context)
}

type publicHandlerImpl struct {
	svc service.PublicService
}

func NewPublicHandler(svc service.PublicService) PublicHandler {
	return &publicHandlerImpl{
		svc: svc,
	}
}

func (p *publicHandlerImpl) GeneratePublicAccessToken(ctx *gin.Context)  {
	accessToken, err := p.svc.GeneratePublicAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message":err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":"access token generated successfully..",
		"token": accessToken,
	})
}