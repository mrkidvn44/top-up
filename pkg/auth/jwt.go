package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"top-up-api/internal/schema"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Interface interface {
	CreateToken(user schema.UserLoginDetail) (string, error)
	AuthenticateService(c *gin.Context) (*jwt.Token, error)
	GetUserFromToken(token *jwt.Token) *schema.UserAuthDetail
	ValidateUserIDFromToken(token *jwt.Token, idStr string) (uint, error)
}

type authService struct {
	jwtSecret []byte
}

var _ Interface = (*authService)(nil)

func NewAuthService(jwtSecret []byte) *authService {
	return &authService{jwtSecret: jwtSecret}
}

func (a *authService) CreateToken(user schema.UserLoginDetail) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.PhoneNumber,
		"id":  user.ID,
		"iss": "top-up-api",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenString, err := claims.SignedString(a.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (a *authService) AuthenticateService(c *gin.Context) (*jwt.Token, error) {

	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.Abort()
		return nil, errors.New("auth: Unauthorized")
	}
	token, err := a.verifyToken(tokenString)
	if err != nil {
		c.Abort()
		return nil, fmt.Errorf("auth: %s", err.Error())
	}
	c.Next()
	return token, nil
}

func (a *authService) verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return a.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (a *authService) GetUserFromToken(token *jwt.Token) *schema.UserAuthDetail {
	user := &schema.UserAuthDetail{}
	claims := token.Claims.(jwt.MapClaims)
	user.PhoneNumber = claims["sub"].(string)
	user.ID = uint(claims["id"].(float64))
	return user
}

func (a *authService) ValidateUserIDFromToken(token *jwt.Token, idStr string) (uint, error) {
	user := a.GetUserFromToken(token)
	idInt, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %w", err)
	}
	if user.ID != uint(idInt) {
		return 0, fmt.Errorf("unauthorized: token user ID does not match request ID")
	}
	return uint(idInt), nil
}
