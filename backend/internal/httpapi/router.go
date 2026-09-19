package httpapi

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/ai"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/auth"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/model"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/security"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

type API struct {
	AI      ai.Provider
	Store   *store.Store
	Tokens  *auth.Manager
	Origins map[string]bool
}

func Router(a API) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), security.Headers(), security.BodyLimit(15<<20), security.RateLimit(120, time.Minute), cors(a.Origins))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.POST("/v1/auth/register", a.register)
	r.POST("/v1/auth/login", a.login)
	v := r.Group("/v1", a.requireAuth())
	v.POST("/match", a.match)
	v.POST("/answers", a.answer)
	v.POST("/resumes/parse", a.parseResume)
	v.POST("/resumes/recommend", a.recommend)
	v.GET("/profile", a.getProfile)
	v.PUT("/profile", a.putProfile)
	v.GET("/applications", a.listApplications)
	v.POST("/applications", a.upsertApplication)
	v.GET("/me/export", a.export)
	v.DELETE("/me", a.deleteMe)
	return r
}
func (a API) register(c *gin.Context) {
	var x model.AuthRequest
	if c.ShouldBindJSON(&x) != nil || !validCredentials(x) {
		bad(c)
		return
	}
	h, e := bcrypt.GenerateFromPassword([]byte(x.Password), bcrypt.DefaultCost)
	if e != nil {
		fail(c)
		return
	}
	id, e := a.Store.CreateUser(c, x.Email, string(h))
	if e != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "account already exists"})
		return
	}
	a.token(c, id)
}
func (a API) login(c *gin.Context) {
	var x model.AuthRequest
	if c.ShouldBindJSON(&x) != nil || !validCredentials(x) {
		bad(c)
		return
	}
	id, h, e := a.Store.UserForLogin(c, x.Email)
	if e != nil || bcrypt.CompareHashAndPassword([]byte(h), []byte(x.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	a.token(c, id)
}
func (a API) token(c *gin.Context, id string) {
	t, exp, e := a.Tokens.Issue(id)
	if e != nil {
		fail(c)
		return
	}
	c.JSON(http.StatusOK, model.AuthResponse{AccessToken: t, ExpiresAt: exp})
}
func (a API) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		id, e := a.Tokens.Parse(raw)
		if e != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("userID", id)
		c.Next()
	}
}
func (a API) match(c *gin.Context) {
	var x model.MatchRequest
	if c.ShouldBindJSON(&x) != nil || x.Job.Description == "" {
		bad(c)
		return
	}
	o, e := a.AI.Match(c, x)
	respond(c, o, e)
}
func (a API) answer(c *gin.Context) {
	var x model.AnswerRequest
	if c.ShouldBindJSON(&x) != nil || x.Question == "" {
		bad(c)
		return
	}
	o, e := a.AI.Answer(c, x)
	respond(c, o, e)
}
func (a API) parseResume(c *gin.Context) {
	var x model.ParseResumeRequest
	if c.ShouldBindJSON(&x) != nil || x.Data == "" || len(x.Data) > 14_000_000 {
		bad(c)
		return
	}
	o, e := a.AI.ParseResume(c, x)
	respond(c, o, e)
}
func (a API) recommend(c *gin.Context) {
	var x model.RecommendRequest
	if c.ShouldBindJSON(&x) != nil || len(x.Resumes) == 0 {
		bad(c)
		return
	}
	o, e := a.AI.RecommendResume(c, x)
	respond(c, o, e)
}
func (a API) getProfile(c *gin.Context) {
	o, e := a.Store.Profile(c, user(c))
	if errors.Is(e, store.ErrNotFound) {
		c.JSON(http.StatusOK, model.Profile{})
		return
	}
	respond(c, o, e)
}
func (a API) putProfile(c *gin.Context) {
	var x model.Profile
	if c.ShouldBindJSON(&x) != nil {
		bad(c)
		return
	}
	if e := a.Store.SaveProfile(c, user(c), x); e != nil {
		fail(c)
		return
	}
	c.Status(http.StatusNoContent)
}
func (a API) listApplications(c *gin.Context) {
	o, e := a.Store.Applications(c, user(c))
	respond(c, o, e)
}
func (a API) upsertApplication(c *gin.Context) {
	var x model.ApplicationInput
	if c.ShouldBindJSON(&x) != nil || x.Job.URL == "" || x.MatchScore < 0 || x.MatchScore > 100 || !validStatus(x.Status) {
		bad(c)
		return
	}
	o, e := a.Store.UpsertApplication(c, user(c), x)
	respond(c, o, e)
}
func (a API) export(c *gin.Context) { o, e := a.Store.Export(c, user(c)); respond(c, o, e) }
func (a API) deleteMe(c *gin.Context) {
	if c.GetHeader("X-Confirm-Delete") != "DELETE MY DATA" {
		c.JSON(http.StatusPreconditionRequired, gin.H{"error": "confirmation header required"})
		return
	}
	if e := a.Store.DeleteUser(c, user(c)); e != nil {
		fail(c)
		return
	}
	c.Status(http.StatusNoContent)
}
func cors(origins map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		o := c.GetHeader("Origin")
		if origins[o] {
			c.Header("Access-Control-Allow-Origin", o)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Confirm-Delete")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			if !origins[o] {
				c.AbortWithStatus(http.StatusForbidden)
			} else {
				c.AbortWithStatus(http.StatusNoContent)
			}
			return
		}
		c.Next()
	}
}
func validCredentials(x model.AuthRequest) bool {
	_, e := mail.ParseAddress(x.Email)
	return e == nil && len(x.Email) <= 254 && len(x.Password) >= 12 && len(x.Password) <= 128
}
func validStatus(s string) bool {
	switch s {
	case "discovered", "prepared", "applied", "in_review", "interview", "rejected", "offer", "withdrawn":
		return true
	}
	return false
}
func user(c *gin.Context) string { return c.GetString("userID") }
func bad(c *gin.Context)         { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"}) }
func fail(c *gin.Context)        { c.JSON(http.StatusInternalServerError, gin.H{"error": "request failed"}) }
func respond(c *gin.Context, o any, e error) {
	if e == nil {
		c.JSON(http.StatusOK, o)
		return
	}
	if errors.Is(e, ai.ErrNotConfigured) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI provider not configured"})
		return
	}
	fail(c)
}
