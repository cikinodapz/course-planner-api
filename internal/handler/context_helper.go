package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// getUserIDFromContext extracts user_id from JWT claims stored by the middleware
func getUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return uuid.Nil, errors.New("Token tidak ditemukan")
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("Claims token tidak valid")
	}

	idValue, ok := claims["user_id"]
	if !ok {
		return uuid.Nil, errors.New("ID pengguna ('user_id') tidak ditemukan di token")
	}

	var idStr string
	switch v := idValue.(type) {
	case string:
		idStr = v
	case uuid.UUID:
		idStr = v.String()
	default:
		if s, ok := idValue.(string); ok {
			idStr = s
		} else {
			return uuid.Nil, errors.New("Tipe ID pengguna tidak valid")
		}
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("ID pengguna tidak valid")
	}

	return userID, nil
}