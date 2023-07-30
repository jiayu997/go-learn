package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// 生成token
func GenerateJwtToken() (string, error) {
	c := MyClaims{
		Username: "admin",
		Password: "admin@123",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 120,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "hnkc",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte("admin"))
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析JWT TOKEN
func ParseJwtToken(token string) error {
	tokenClaims, err := jwt.ParseWithClaims(token, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("admin"), nil
	})
	if err != nil {
		return err
	}
	//fmt.Println(tokenClaims.Claims.(*MyClaims))
	if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid {
		if claims.Username == "admin" && claims.Password == "admin@123" {
			return nil
		} else {
			return errors.New("token verify failed")
		}
	}
	return errors.New("token verify failed")
}
