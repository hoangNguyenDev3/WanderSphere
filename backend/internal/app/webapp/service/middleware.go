package service

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangNguyenDev3/WanderSphere/backend/internal/pkg/types"
)

// AuthRequired is a middleware that checks if the user is authenticated
// using session-based authentication
func (svc *WebService) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check session authentication
		_, userId, err := svc.checkSessionAuthentication(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, types.ErrorResponse{
				Error:   "unauthorized",
				Code:    http.StatusUnauthorized,
				Message: "Session is invalid or expired",
			})
			return
		}

		// Set user ID in context for later use
		c.Set("user_id", userId)
		c.Next()
	}
}

// RefreshSession is a middleware that refreshes the session expiration time
func (svc *WebService) RefreshSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, _, err := svc.checkSessionAuthentication(c)
		if err == nil {
			// Get session configuration from config
			var expirationTime time.Duration

			// Use config value if available, otherwise use a reasonable default
			if svc.Config != nil && svc.Config.Auth.Session.ExpirationMinutes > 0 {
				expirationTime = time.Minute * time.Duration(svc.Config.Auth.Session.ExpirationMinutes)
			} else {
				expirationTime = time.Hour * 24 // Default to 24 hours
			}

			// Refresh session in Redis
			_ = svc.RedisClient.Expire(svc.RedisClient.Context(), sessionId, expirationTime).Err()

			// Get cookie settings
			cookieName := "session_id" // Default
			if svc.Config != nil && svc.Config.Auth.Session.CookieName != "" {
				cookieName = svc.Config.Auth.Session.CookieName
			}

			// Default to secure settings unless explicitly configured otherwise
			secure := true
			httpOnly := true
			sameSite := http.SameSiteStrictMode

			if svc.Config != nil {
				// Only override defaults if explicitly set in config
				if svc.Config.Auth.Session.Secure == false {
					secure = false
				}
				if svc.Config.Auth.Session.HTTPOnly == false {
					httpOnly = false
				}

				if svc.Config.Auth.Session.SameSite == "lax" {
					sameSite = http.SameSiteLaxMode
				} else if svc.Config.Auth.Session.SameSite == "none" {
					sameSite = http.SameSiteNoneMode
				}
			}

			// Refresh cookie
			maxAge := int(expirationTime.Seconds())
			c.SetSameSite(sameSite)
			c.SetCookie(cookieName, sessionId, maxAge, "/", "", secure, httpOnly)
		}

		c.Next()
	}
}
