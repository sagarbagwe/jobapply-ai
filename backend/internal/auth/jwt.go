package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}
type Claims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func New(secret string, ttl time.Duration) *Manager { return &Manager{[]byte(secret), ttl} }
func (m *Manager) Issue(userID string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(m.ttl)
	c := Claims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{Subject: userID, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(exp), Issuer: "jobapply-ai", Audience: jwt.ClaimStrings{"jobapply-extension"}}}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, e := t.SignedString(m.secret)
	return s, exp, e
}
func (m *Manager) Parse(raw string) (string, error) {
	t, e := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer("jobapply-ai"), jwt.WithAudience("jobapply-extension"))
	if e != nil || !t.Valid {
		return "", errors.New("invalid token")
	}
	c, ok := t.Claims.(*Claims)
	if !ok || c.UserID == "" {
		return "", errors.New("invalid claims")
	}
	return c.UserID, nil
}
