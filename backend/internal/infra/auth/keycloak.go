package auth

import (
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type KeycloakAuth struct {
	jwks *keyfunc.JWKS
}

func NewKeycloakAuth(jwksURL string) (*KeycloakAuth, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}
	return &KeycloakAuth{jwks: jwks}, nil
}

func (k *KeycloakAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, k.jwks.Keyfunc)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidas"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["sub"])
		c.Set("email", claims["email"])
		if familyID, exists := claims["family_id"]; exists {
			c.Set("family_id", familyID)
		}

		c.Next()
	}
}