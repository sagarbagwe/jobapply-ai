package model

import "time"

type Profile struct {
	FullName           string   `json:"full_name"`
	Email              string   `json:"email"`
	Phone              string   `json:"phone"`
	Location           string   `json:"location"`
	YearsExperience    float64  `json:"years_experience"`
	Skills             []string `json:"skills"`
	Education          string   `json:"education,omitempty"`
	Degree             string   `json:"degree,omitempty"`
	GraduationYear     string   `json:"graduation_year,omitempty"`
	Companies          []string `json:"companies,omitempty"`
	JobTitles          []string `json:"job_titles,omitempty"`
	Projects           []string `json:"projects,omitempty"`
	Certifications     []string `json:"certifications,omitempty"`
	CurrentCompany     string   `json:"current_company,omitempty"`
	CurrentCTC         string   `json:"current_ctc,omitempty"`
	ExpectedCTC        string   `json:"expected_ctc,omitempty"`
	NoticePeriod       string   `json:"notice_period,omitempty"`
	PreferredLocations []string `json:"preferred_locations,omitempty"`
	WorkAuthorization  string   `json:"work_authorization,omitempty"`
}
type Job struct {
	Company     string   `json:"company"`
	Title       string   `json:"title"`
	Location    string   `json:"location"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Source      string   `json:"source"`
	Skills      []string `json:"skills,omitempty"`
}
type MatchRequest struct {
	Profile Profile `json:"profile"`
	Job     Job     `json:"job"`
}
type MatchResult struct {
	MatchScore      int      `json:"match_score"`
	Category        string   `json:"category"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	ExperienceMatch bool     `json:"experience_match"`
	LocationMatch   bool     `json:"location_match"`
	Recommendation  string   `json:"recommendation"`
	Reasons         []string `json:"reasons"`
}
type AnswerRequest struct {
	Profile  Profile `json:"profile"`
	Job      Job     `json:"job"`
	Question string  `json:"question"`
}
type AnswerResult struct {
	Answer         string `json:"answer"`
	RequiresReview bool   `json:"requires_review"`
}
type ParseResumeRequest struct {
	Name     string `json:"name"`
	MIMEType string `json:"mime_type"`
	Data     string `json:"data"`
}
type ParseResumeResult struct {
	Profile       Profile  `json:"profile"`
	Skills        []string `json:"skills"`
	MissingFields []string `json:"missing_fields"`
}
type ResumeSummary struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Profile Profile  `json:"profile"`
	Tags    []string `json:"tags"`
}
type RecommendRequest struct {
	Job     Job             `json:"job"`
	Resumes []ResumeSummary `json:"resumes"`
}
type RecommendResult struct {
	ResumeID string `json:"resume_id"`
	Reason   string `json:"reason"`
}
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
type ApplicationInput struct {
	Job        Job    `json:"job"`
	Status     string `json:"status"`
	MatchScore int    `json:"match_score"`
	Notes      string `json:"notes"`
}
type Application struct {
	ID         string     `json:"id"`
	Job        Job        `json:"job"`
	Status     string     `json:"status"`
	MatchScore int        `json:"match_score"`
	Notes      string     `json:"notes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	AppliedAt  *time.Time `json:"applied_at,omitempty"`
}
