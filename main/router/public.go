package router

import (
	"github.com/achsanit/go-assignment2/main/handler"
	"github.com/gin-gonic/gin"
)

type PublicRouter interface {
	Mount()
}

type publicRouterImpl struct {
	v *gin.RouterGroup
	h handler.PublicHandler
}

func NewPublicRouter(v *gin.RouterGroup, h handler.PublicHandler) PublicRouter {
	return &publicRouterImpl{
		v: v,
		h: h,
	}
}

func (p *publicRouterImpl) Mount() {
	p.v.GET("generate_token", p.h.GeneratePublicAccessToken)
}