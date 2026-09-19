package model

type Profile struct { FullName string `json:"full_name"`; Skills []string `json:"skills"`; YearsExperience float64 `json:"years_experience"`; Locations []string `json:"preferred_locations"` }
type Job struct { Company string `json:"company"`; Title string `json:"title"`; Location string `json:"location"`; Description string `json:"description"`; URL string `json:"url"` }
type MatchRequest struct { Profile Profile `json:"profile"`; Job Job `json:"job"` }
type MatchResult struct { MatchScore int `json:"match_score"`; MatchedSkills []string `json:"matched_skills"`; MissingSkills []string `json:"missing_skills"`; ExperienceMatch bool `json:"experience_match"`; LocationMatch bool `json:"location_match"`; Recommendation string `json:"recommendation"`; Rationale []string `json:"rationale"` }
type AnswerRequest struct { Profile Profile `json:"profile"`; Job Job `json:"job"`; Question string `json:"question"` }
type AnswerResult struct { Answer string `json:"answer"`; RequiresReview bool `json:"requires_review"` }
