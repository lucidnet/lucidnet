package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/lucidnet/lucidnet/internal/app/vertex/config"
	"strings"
	"time"
)

type Authenticator struct {
	configManager *config.Manager
}

func NewAuthenticator(configManager *config.Manager) *Authenticator {
	return &Authenticator{
		configManager: configManager,
	}
}

const (
	BearerToken    = "Bearer"
	TokenTypeAdmin = "admin"
)

func (a *Authenticator) GenerateAdminToken(identifier string) (string, error) {
	vertexEndpoint, err := a.configManager.GetVertexEndpoint()

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  identifier,
		"iss":  vertexEndpoint,
		"aud":  vertexEndpoint,
		"type": TokenTypeAdmin,
		"iat":  int(time.Now().Unix()),
		"exp":  int(time.Now().Add(time.Hour * 24).Unix()),
	})

	jwtSigningKey, err := a.configManager.GetJwtSigningKey("secret")

	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(jwtSigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authenticator) ValidateAdminContext(c *gin.Context) (*AuthenticatedAdmin, error) {
	tokenType, token, err := a.extractToken(c)

	if err != nil {
		return nil, err
	}

	switch tokenType {
	case BearerToken:
		return a.validateAdminToken(token)
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func (a *Authenticator) validateAdminToken(tokenString string) (*AuthenticatedAdmin, error) {
	jwtSigningKey, err := a.configManager.GetJwtSigningKey("secret")

	if err != nil {
		return nil, err
	}

	vertexEndpoint, err := a.configManager.GetVertexEndpoint()

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"]
		if !ok {
			return nil, fmt.Errorf("invalid token type")
		}

		tokenTypeString, ok := tokenType.(string)
		if !ok {
			return nil, fmt.Errorf("invalid token type")
		}

		if tokenTypeString != TokenTypeAdmin {
			return nil, fmt.Errorf("invalid token type")
		}

		issuer, ok := claims["iss"]

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		issuerString, ok := issuer.(string)

		if !ok {
			return nil, fmt.Errorf("invalid token issuer")
		}

		if issuerString != vertexEndpoint {
			return nil, fmt.Errorf("invalid token issuer")
		}

		aud, ok := claims["aud"]

		if !ok {
			return nil, fmt.Errorf("invalid token audience")
		}

		audString, ok := aud.(string)

		if !ok {
			return nil, fmt.Errorf("invalid token audience")
		}

		if audString != vertexEndpoint {
			return nil, fmt.Errorf("invalid token audience")
		}

		identifier, ok := claims["sub"]

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		identifierString, ok := identifier.(string)

		if !ok {
			return nil, fmt.Errorf("invalid access token")
		}

		return &AuthenticatedAdmin{
			Identifier: identifierString,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid access token")
	}
}

func (a *Authenticator) extractToken(c *gin.Context) (string, string, error) {
	authorizationHeader := c.GetHeader("Authorization")

	if len(authorizationHeader) == 0 {
		return "", "", fmt.Errorf("authorization header is not set")
	}

	components := strings.Split(authorizationHeader, " ")

	if len(components) != 2 {
		return "", "", fmt.Errorf("invalid access token")
	}

	return components[0], components[1], nil
}
