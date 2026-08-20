package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"mahrec/internal/store"
)

const (
	CookieName   = "client_id"
	localsKey    = "clientID"
	cookieMaxAge = 10 * 365 * 24 * time.Hour // effectively persists until cleared
)

// ClientID ensures every visitor has a stable, httponly-cookie-backed id.
// No login, no signup: first visit mints a uuid, later visits reuse it.
func ClientID(clients *store.ClientStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Cookies(CookieName)

		if id == "" {
			id = uuid.NewString()
			c.Cookie(&fiber.Cookie{
				Name:     CookieName,
				Value:    id,
				HTTPOnly: true,
				SameSite: fiber.CookieSameSiteLaxMode,
				Path:     "/",
				MaxAge:   int(cookieMaxAge.Seconds()),
			})
		}

		if err := clients.Ensure(id); err != nil {
			return err
		}

		c.Locals(localsKey, id)
		return c.Next()
	}
}

// FromCtx reads the client id previously set by the ClientID middleware.
func FromCtx(c fiber.Ctx) string {
	id, _ := c.Locals(localsKey).(string)
	return id
}
