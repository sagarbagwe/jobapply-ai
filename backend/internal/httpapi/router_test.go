package httpapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeAI struct{}

func (fakeAI) Match(context.Context, model.MatchRequest) (model.MatchResult, error) {
	return model.MatchResult{MatchScore: 80}, nil
}
func (fakeAI) Answer(context.Context, model.AnswerRequest) (model.AnswerResult, error) {
	return model.AnswerResult{}, nil
}
func (fakeAI) ParseResume(context.Context, model.ParseResumeRequest) (model.ParseResumeResult, error) {
	return model.ParseResumeResult{}, nil
}
func (fakeAI) RecommendResume(context.Context, model.RecommendRequest) (model.RecommendResult, error) {
	return model.RecommendResult{}, nil
}
func TestValidCredentials(t *testing.T) {
	if !validCredentials(model.AuthRequest{Email: "user@example.com", Password: "correct-horse-battery"}) {
		t.Fatal("expected valid")
	}
	if validCredentials(model.AuthRequest{Email: "bad", Password: "short"}) {
		t.Fatal("expected invalid")
	}
}
func TestCORSRejectsUnknownPreflight(t *testing.T) {
	r := ginTestRouter()
	w := httptest.NewRecorder()
	q := httptest.NewRequest(http.MethodOptions, "/", nil)
	q.Header.Set("Origin", "https://evil.example")
	r.ServeHTTP(w, q)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status=%d", w.Code)
	}
}
func ginTestRouter() http.Handler {
	imported := struct{}{}
	_ = imported
	r := gin.New()
	r.Use(cors(map[string]bool{"chrome-extension://ok": true}))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	return r
}
