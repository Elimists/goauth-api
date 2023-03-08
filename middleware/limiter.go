package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	limiter "github.com/gofiber/fiber/v2/middleware/limiter"
)

//Limits the # of requests a ip can send.
//
//Takes in the maximum number of connections (# of request allowed as integer) and expirate time (retry after in seconds)
func Limiter(maximumNumOfConnections int, expirationTimeInSeconds time.Duration) func(*fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        maximumNumOfConnections,
		Expiration: expirationTimeInSeconds * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}
