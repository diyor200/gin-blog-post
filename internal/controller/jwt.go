package controller

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

//	type TokenClaims struct {
//		jwt.StandardClaims
//		UserID int
//	}
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

var jwtSect = os.Getenv("PASSWORD_SALT")

// Generate jwt-token
func GenerateToken(userID int) (string, error) {
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
	//	StandardClaims: jwt.StandardClaims{
	//		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	//		IssuedAt:  time.Now().Unix(),
	//	},
	//	UserID: int(userID),
	//})
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
	})
	log.Println("config.JwtSecret=", jwtSect)
	tokenString, err := token.SignedString([]byte(jwtSect))
	return tokenString, err
}

// ParseToken - parses token
func ParseToken(accessToken string) (int, error) {
	//token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		log.Println("ParseToken - jwt.ParseWithClaims - token.Method.(*jwt.SigningMethodHMAC) == ", ok)
	//		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	//	}
	//	return []byte(config.JwtAccessKey), nil
	//})
	//log.Println("ParseToken - token== ", token)
	//if err != nil {
	//	log.Println("ParseToken - cannot parse token")
	//	return 0, errors.New("cannot parse token")
	//}
	//claims, ok := token.Claims.(*TokenClaims)
	//if !ok {
	//	log.Println("ParseToken - token.Claims.(*TokenClaims) == ", ok)
	//	return 0, errors.New("cannot parse token")
	//}
	//log.Println("ParseToken - UserID== ", claims.UserID)
	//return claims.UserID, nil

	//t, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return 0, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	//
	//	}
	//	return os.Getenv("PASSWORD_SALT"), nil
	//})

	t, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSect), nil
	})
	if err != nil {
		log.Println("err := jwt.ParseWithClaims ==== ", err)
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	userId, err := strconv.Atoi(subject)
	log.Println(userId)
	if err != nil {
		return 0, errors.New("invalid subject")
	}
	log.Println("parseToken - userId=", userId)
	return userId, nil

}
