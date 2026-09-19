package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env, HTTPAddr, DatabaseURL, GeminiKey, GeminiModel, JWTSecret string
	AllowedOrigins                                                []string
	AccessTTL                                                     time.Duration
}

func Load() (Config, error) {
	c := Config{Env: get("APP_ENV", "development"), HTTPAddr: get("HTTP_ADDR", ":8080"), DatabaseURL: os.Getenv("DATABASE_URL"), GeminiKey: os.Getenv("GEMINI_API_KEY"), GeminiModel: get("GEMINI_MODEL", "gemini-2.5-flash"), JWTSecret: os.Getenv("JWT_SECRET"), AllowedOrigins: split(os.Getenv("ALLOWED_ORIGINS")), AccessTTL: duration("ACCESS_TOKEN_TTL_MINUTES", 15) * time.Minute}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL is required")
	}
	if c.Env == "production" {
		if len(c.JWTSecret) < 32 {
			return c, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
		}
		if c.GeminiKey == "" {
			return c, fmt.Errorf("GEMINI_API_KEY is required")
		}
		if len(c.AllowedOrigins) == 0 {
			return c, fmt.Errorf("ALLOWED_ORIGINS is required")
		}
	}
	if c.JWTSecret == "" {
		c.JWTSecret = "development-only-secret-change-before-production"
	}
	return c, nil
}
func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func split(v string) []string {
	var out []string
	for _, x := range strings.Split(v, ",") {
		if x = strings.TrimSpace(x); x != "" {
			out = append(out, x)
		}
	}
	return out
}
func duration(k string, d int) time.Duration {
	v, err := strconv.Atoi(os.Getenv(k))
	if err != nil || v < 1 {
		return time.Duration(d)
	}
	return time.Duration(v)
}
