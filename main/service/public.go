package service

import (
	"fmt"
	"time"

	"github.com/achsanit/go-assignment2/main/helper"
	"github.com/achsanit/go-assignment2/main/model"
)

type PublicService interface {
	GeneratePublicAccessToken() (accessToken string, err error)
}

type publicServiceImpl struct{}

// GeneratePublicAccessToken implements PublicService.
func (*publicServiceImpl) GeneratePublicAccessToken() (accessToken string, err error) {
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-assignment",
		Aud: "go-assignment",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	accessToken, err = helper.GenerateToken(claim)
	return
}

func NewPublicService() PublicService {
	return &publicServiceImpl{}
}
