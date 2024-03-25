package helper

import (
	"encoding/json"
	"log"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRET_JWT = "secret.jwt.go.assignment.achsanit"
)

func GenerateToken(claim any) (token string, err error) {
	jwtClaim := jwt.MapClaims{}
	b, err := json.Marshal(claim)
	if err != nil {
		log.Println("cannot marshal claim payload")
		return
	}
	err = json.Unmarshal(b, &jwtClaim)
	if err != nil {
		log.Println("cannot mapping claim to jwt claim")
		return
	}
	// prepare
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaim)
	// generate token
	token, err = parseToken.SignedString([]byte(SECRET_JWT))
	if err != nil {
		log.Println("cannot generate token", err.Error())
		return
	}
	return
}

func ValidateToken(token string) (claim jwt.MapClaims, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(SECRET_JWT), nil
	})
	if err != nil {
		log.Println("error validating jwt token", err.Error())
		return
	}

	// translate claim
	claim, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("error translate claim")
		return
	}
	return
}