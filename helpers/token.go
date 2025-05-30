package helpers

import (
	"fmt"
	"go-jwt/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("mysecretkey")

type MyCustomClaims struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func CreateToken(user *models.User) (string, error) {
	claims := MyCustomClaims{
		user.ID,
		user.Name,
		user.Email,
		user.Role,
	jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	},
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
ss, err := token.SignedString(mySigningKey)

return ss, err
}

func VerifyToken(tokenString string) (any, error) {

token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	return mySigningKey, nil
})
if err != nil {
	return nil, fmt.Errorf("unauthorized")
} 

claims, ok := token.Claims.(*MyCustomClaims)

if !ok || !token.Valid {
	return nil, fmt.Errorf("unauthorized")
} 

return claims, nil

}