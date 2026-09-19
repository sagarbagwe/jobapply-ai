package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/model"
	"io"
	"net/http"
	"strings"
	"time"
)

type Gemini struct {
	apiKey, model string
	client        *http.Client
}

func NewGemini(k, m string) *Gemini {
	return &Gemini{apiKey: k, model: m, client: &http.Client{Timeout: 40 * time.Second}}
}

type part struct {
	Text       string      `json:"text,omitempty"`
	InlineData *inlineData `json:"inlineData,omitempty"`
}
type inlineData struct {
	MIMEType string `json:"mimeType"`
	Data     string `json:"data"`
}
type content struct {
	Parts []part `json:"parts"`
}
type request struct {
	Contents         []content `json:"contents"`
	GenerationConfig config    `json:"generationConfig"`
}
type config struct {
	ResponseMIMEType string `json:"responseMimeType"`
}
type response struct {
	Candidates []struct {
		Content content `json:"content"`
	} `json:"candidates"`
}

func (g *Gemini) generate(ctx context.Context, parts []part, out any) error {
	if g.apiKey == "" {
		return ErrNotConfigured
	}
	body, _ := json.Marshal(request{Contents: []content{{Parts: parts}}, GenerationConfig: config{ResponseMIMEType: "application/json"}})
	u := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return fmt.Errorf("gemini status %d: %s", res.StatusCode, string(b))
	}
	var r response
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	if len(r.Candidates) == 0 || len(r.Candidates[0].Content.Parts) == 0 {
		return fmt.Errorf("empty AI response")
	}
	return json.Unmarshal([]byte(strings.TrimSpace(r.Candidates[0].Content.Parts[0].Text)), out)
}
func textPrompt(v any, instruction string) []part {
	b, _ := json.Marshal(v)
	return []part{{Text: instruction + "\nInput JSON: " + string(b)}}
}
func (g *Gemini) Match(ctx context.Context, in model.MatchRequest) (out model.MatchResult, err error) {
	err = g.generate(ctx, textPrompt(in, "Return strict JSON. Internal relevance only, not interview probability. Use only input facts. Fields: match_score 0-100, category HIGH MATCH|GOOD MATCH|LOW MATCH, matched_skills, missing_skills, experience_match, location_match, recommendation apply|review|ignore, reasons."), &out)
	return
}
func (g *Gemini) Answer(ctx context.Context, in model.AnswerRequest) (out model.AnswerResult, err error) {
	err = g.generate(ctx, textPrompt(in, "Draft a concise truthful application answer using only supplied facts. Never invent experience, employment, education, achievements, salary, projects, certifications, or work authorization. JSON fields: answer, requires_review=true."), &out)
	out.RequiresReview = true
	return
}
func (g *Gemini) ParseResume(ctx context.Context, in model.ParseResumeRequest) (out model.ParseResumeResult, err error) {
	if in.MIMEType != "application/pdf" {
		return out, fmt.Errorf("unsupported resume type")
	}
	prompt := "Extract only facts explicitly present in this resume. Use empty values for unknowns; never infer. JSON: profile with full_name,email,phone,location,years_experience,skills,education,degree,graduation_year,companies,job_titles,projects,certifications,current_company,current_ctc,expected_ctc,notice_period,preferred_locations,work_authorization; skills; missing_fields."
	err = g.generate(ctx, []part{{Text: prompt}, {InlineData: &inlineData{MIMEType: in.MIMEType, Data: in.Data}}}, &out)
	return
}
func (g *Gemini) RecommendResume(ctx context.Context, in model.RecommendRequest) (out model.RecommendResult, err error) {
	err = g.generate(ctx, textPrompt(in, "Choose the resume best grounded in the job description. Return JSON fields resume_id and concise reason. Do not invent qualifications."), &out)
	return
}
