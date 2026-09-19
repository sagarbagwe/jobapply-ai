package store

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sagarbagwe/jobapply-ai/backend/internal/model"
	"time"
)

var ErrNotFound = errors.New("not found")

type Store struct{ DB *pgxpool.Pool }

func Open(ctx context.Context, url string) (*Store, error) {
	p, e := pgxpool.New(ctx, url)
	if e != nil {
		return nil, e
	}
	if e = p.Ping(ctx); e != nil {
		p.Close()
		return nil, e
	}
	return &Store{p}, nil
}
func (s *Store) Close() { s.DB.Close() }
func (s *Store) CreateUser(ctx context.Context, email, hash string) (string, error) {
	var id string
	e := s.DB.QueryRow(ctx, "INSERT INTO users(email,password_hash) VALUES(lower($1),$2) RETURNING id", email, hash).Scan(&id)
	return id, e
}
func (s *Store) UserForLogin(ctx context.Context, email string) (string, string, error) {
	var id, h string
	e := s.DB.QueryRow(ctx, "SELECT id,password_hash FROM users WHERE email=lower($1) AND deleted_at IS NULL", email).Scan(&id, &h)
	if errors.Is(e, pgx.ErrNoRows) {
		e = ErrNotFound
	}
	return id, h, e
}
func (s *Store) SaveProfile(ctx context.Context, user string, p model.Profile) error {
	b, _ := json.Marshal(p)
	_, e := s.DB.Exec(ctx, "INSERT INTO profiles(user_id,data) VALUES($1,$2) ON CONFLICT(user_id) DO UPDATE SET data=excluded.data,updated_at=now()", user, b)
	return e
}
func (s *Store) Profile(ctx context.Context, user string) (model.Profile, error) {
	var b []byte
	e := s.DB.QueryRow(ctx, "SELECT data FROM profiles WHERE user_id=$1", user).Scan(&b)
	if errors.Is(e, pgx.ErrNoRows) {
		e = ErrNotFound
	}
	var p model.Profile
	if e == nil {
		e = json.Unmarshal(b, &p)
	}
	return p, e
}
func (s *Store) UpsertApplication(ctx context.Context, user string, in model.ApplicationInput) (model.Application, error) {
	tx, e := s.DB.Begin(ctx)
	if e != nil {
		return model.Application{}, e
	}
	defer tx.Rollback(ctx)
	var jobID string
	e = tx.QueryRow(ctx, `INSERT INTO jobs(canonical_url,company,title,location,description,source) VALUES($1,$2,$3,$4,$5,$6) ON CONFLICT(canonical_url) DO UPDATE SET company=excluded.company,title=excluded.title,location=excluded.location,description=excluded.description RETURNING id`, in.Job.URL, in.Job.Company, in.Job.Title, in.Job.Location, in.Job.Description, in.Job.Source).Scan(&jobID)
	if e != nil {
		return model.Application{}, e
	}
	var a model.Application
	e = tx.QueryRow(ctx, `INSERT INTO applications(user_id,job_id,status,match_score,notes,applied_at) VALUES($1,$2,$3,$4,$5,CASE WHEN $3='applied' THEN now() END) ON CONFLICT(user_id,job_id) DO UPDATE SET status=excluded.status,match_score=excluded.match_score,notes=excluded.notes,updated_at=now(),applied_at=COALESCE(applications.applied_at,excluded.applied_at) RETURNING id,status,match_score,notes,created_at,updated_at,applied_at`, user, jobID, in.Status, in.MatchScore, in.Notes).Scan(&a.ID, &a.Status, &a.MatchScore, &a.Notes, &a.CreatedAt, &a.UpdatedAt, &a.AppliedAt)
	if e != nil {
		return a, e
	}
	a.Job = in.Job
	if e = tx.Commit(ctx); e != nil {
		return a, e
	}
	return a, nil
}
func (s *Store) Applications(ctx context.Context, user string) ([]model.Application, error) {
	rows, e := s.DB.Query(ctx, `SELECT a.id,a.status,a.match_score,a.notes,a.created_at,a.updated_at,a.applied_at,j.company,j.title,j.location,j.canonical_url,j.description,j.source FROM applications a JOIN jobs j ON j.id=a.job_id WHERE a.user_id=$1 ORDER BY a.updated_at DESC LIMIT 500`, user)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []model.Application{}
	for rows.Next() {
		var a model.Application
		if e = rows.Scan(&a.ID, &a.Status, &a.MatchScore, &a.Notes, &a.CreatedAt, &a.UpdatedAt, &a.AppliedAt, &a.Job.Company, &a.Job.Title, &a.Job.Location, &a.Job.URL, &a.Job.Description, &a.Job.Source); e != nil {
			return nil, e
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
func (s *Store) Export(ctx context.Context, user string) (map[string]any, error) {
	p, _ := s.Profile(ctx, user)
	a, e := s.Applications(ctx, user)
	return map[string]any{"profile": p, "applications": a, "exported_at": time.Now().UTC()}, e
}
func (s *Store) DeleteUser(ctx context.Context, user string) error {
	_, e := s.DB.Exec(ctx, "DELETE FROM users WHERE id=$1", user)
	return e
}
