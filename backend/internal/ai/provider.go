package ai

import("context";"errors";"github.com/sagarbagwe/jobapply-ai/backend/internal/model")
var ErrNotConfigured=errors.New("AI provider is not configured")
type Provider interface { Match(context.Context,model.MatchRequest)(model.MatchResult,error); Answer(context.Context,model.AnswerRequest)(model.AnswerResult,error) }
