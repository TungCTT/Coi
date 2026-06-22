package middleware

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var devOriginPorts = map[string]bool{
	"3000": true,
	"5173": true,
}

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowed := map[string]bool{}
	allowAnyOrigin := false
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAnyOrigin = true
			continue
		}
		allowed[strings.TrimRight(origin, "/")] = true
	}

	return func(c *gin.Context) {
		origin := strings.TrimRight(c.GetHeader("Origin"), "/")
		if origin != "" {
			if allowAnyOrigin {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if allowed[origin] || isLocalDevOrigin(origin) {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		c.Header("Access-Control-Expose-Headers", "ETag")
		c.Header("Access-Control-Max-Age", "86400")
		if c.GetHeader("Access-Control-Request-Private-Network") == "true" {
			c.Header("Access-Control-Allow-Private-Network", "true")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isLocalDevOrigin(origin string) bool {
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" {
		return false
	}
	if !devOriginPorts[parsed.Port()] {
		return false
	}

	host := parsed.Hostname()
	if host == "localhost" || host == "0.0.0.0" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate()
}
