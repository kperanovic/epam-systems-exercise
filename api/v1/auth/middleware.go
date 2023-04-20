package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kperanovic/epam-systems/internal/token"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "authorization_payload"
)

// AuthMiddleware is responsabile for request authentication.
// It accepts JWT token. Checks if the header is provided and is the header in the right format.
// Checks the validation type, and then validates the token sent in the header.
func AuthMiddleware(t token.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": err,
				},
			)

			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": err,
				},
			)

			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("unsupported authorization type %v", authType)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": err,
				},
			)

			return
		}

		accessToken := fields[1]
		payload, err := t.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Set(authPayloadKey, payload)
		c.Next()
	}
}
