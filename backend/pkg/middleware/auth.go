package middleware

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
	"yoga-backend/config"
	"yoga-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims defines the JWT claims structure
type Claims struct {
	Role models.UserRole
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a user.
func GenerateJWT(userID string, role models.UserRole, cfg *config.Config) (string, string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token valid for 24 hours
	accessTokenClaims := &Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %w", err)
	}

	expirationTime = time.Now().Add(7 * 24 * time.Hour) // Token valid for 24 hours
	refreshTokenClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// AuthMiddleware validates the JWT token from the request header.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

// AuthorizeRole creates a middleware that checks if the user has one of the allowed roles.
func AuthorizeRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleAny, exists := c.Get("userRole")
		if !exists {
			log.Println("userRole not found in context, AuthMiddleware might be missing or failed.")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		userRole, ok := userRoleAny.(models.UserRole)
		if !ok {
			log.Printf("Failed to cast userRole to models.UserRole, actual type: %T", userRoleAny)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user role type in context"})
			c.Abort()
			return
		}

		if slices.Contains(allowedRoles, userRole) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

// ValidateToken validates a JWT string and returns the claims.
func ValidateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
