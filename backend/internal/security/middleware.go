package security

import (
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type bucket struct {
	n     int
	reset time.Time
}

var mu sync.Mutex
var buckets = map[string]bucket{}

func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Cache-Control", "no-store")
		c.Next()
	}
}
func BodyLimit(n int64) gin.HandlerFunc {
	return func(c *gin.Context) { c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n); c.Next() }
}
func RateLimit(max int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		now := time.Now()
		mu.Lock()
		b := buckets[key]
		if now.After(b.reset) {
			b = bucket{reset: now.Add(window)}
		}
		b.n++
		buckets[key] = b
		mu.Unlock()
		if b.n > max {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
func APIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		want := os.Getenv("API_AUTH_TOKEN")
		if want == "" && os.Getenv("APP_ENV") != "production" {
			c.Next()
			return
		}
		got := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if want == "" || subtle.ConstantTimeCompare([]byte(got), []byte(want)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
