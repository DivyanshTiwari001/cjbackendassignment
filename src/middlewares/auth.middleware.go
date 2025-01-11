package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is the middleware to authenticate requests using JWT
func AuthMiddleware(c *fiber.Ctx) error {
	// Retrieve access token from cookies
	accessToken := c.Cookies("accessToken")

	// If no access token is found, return unauthorized error
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No access token found",
		})
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		// Return the secret key for token validation
		return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
	})

	// If error in parsing the token or token is invalid
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired access token",
		})
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["id"] == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Add user ID to the request context
	c.Locals("userId", claims["id"])

	// Continue to the next handler
	return c.Next()
}
