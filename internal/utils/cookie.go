package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// RefreshTokenCookieName is the cookie key used to transport refresh tokens
const RefreshTokenCookieName = "refreshToken"

// SetRefreshTokenCookie writes the refresh token as an HTTP-only cookie
func SetRefreshTokenCookie(c *fiber.Ctx, token string, maxAge time.Duration, secure bool) {
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(maxAge),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	})
}

// ClearRefreshTokenCookie removes the refresh token cookie
func ClearRefreshTokenCookie(c *fiber.Ctx, secure bool) {
	c.Cookie(&fiber.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	})
}
