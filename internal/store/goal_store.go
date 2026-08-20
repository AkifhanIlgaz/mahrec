package store

import (
	"database/sql"
	"time"

	"mahrec/internal/models"
)

type GoalStore struct {
	db *sql.DB
}

func NewGoalStore(db *sql.DB) *GoalStore {
	return &GoalStore{db: db}
}

func (s *GoalStore) List(clientID string) ([]models.Goal, error) {
	rows, err := s.db.Query(`
		SELECT id, client_id, title, description, target_count, current_count, created_at, completed_at
		FROM goals
		WHERE client_id = ?
		ORDER BY completed_at IS NOT NULL, created_at DESC
	`, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		g, err := scanGoal(rows)
		if err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, rows.Err()
}

func (s *GoalStore) Get(clientID string, id int64) (models.Goal, error) {
	row := s.db.QueryRow(`
		SELECT id, client_id, title, description, target_count, current_count, created_at, completed_at
		FROM goals
		WHERE id = ? AND client_id = ?
	`, id, clientID)
	return scanGoal(row)
}

func (s *GoalStore) Create(clientID, title, description string, targetCount int) (models.Goal, error) {
	res, err := s.db.Exec(`
		INSERT INTO goals (client_id, title, description, target_count) VALUES (?, ?, ?, ?)
	`, clientID, title, description, targetCount)
	if err != nil {
		return models.Goal{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Goal{}, err
	}
	return s.Get(clientID, id)
}

func (s *GoalStore) Increment(clientID string, id int64, by int) (models.Goal, error) {
	_, err := s.db.Exec(`
		UPDATE goals
		SET current_count = MIN(current_count + ?, target_count),
		    completed_at = CASE
		        WHEN completed_at IS NULL AND current_count + ? >= target_count
		            THEN datetime('now')
		        ELSE completed_at
		    END
		WHERE id = ? AND client_id = ?
	`, by, by, id, clientID)
	if err != nil {
		return models.Goal{}, err
	}
	return s.Get(clientID, id)
}

func (s *GoalStore) Reset(clientID string, id int64) (models.Goal, error) {
	_, err := s.db.Exec(`
		UPDATE goals
		SET current_count = 0, completed_at = NULL
		WHERE id = ? AND client_id = ?
	`, id, clientID)
	if err != nil {
		return models.Goal{}, err
	}
	return s.Get(clientID, id)
}

func (s *GoalStore) Delete(clientID string, id int64) error {
	_, err := s.db.Exec(`DELETE FROM goals WHERE id = ? AND client_id = ?`, id, clientID)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanGoal(row rowScanner) (models.Goal, error) {
	var g models.Goal
	var createdAt string
	var completedAt sql.NullString

	err := row.Scan(&g.ID, &g.ClientID, &g.Title, &g.Description, &g.TargetCount, &g.CurrentCount, &createdAt, &completedAt)
	if err != nil {
		return models.Goal{}, err
	}

	g.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return models.Goal{}, err
	}

	if completedAt.Valid {
		t, err := parseTime(completedAt.String)
		if err != nil {
			return models.Goal{}, err
		}
		g.CompletedAt = &t
	}

	return g, nil
}

// parseTime accepts both the sqlite text format ("YYYY-MM-DD HH:MM:SS")
// and RFC3339, since the driver's handling depends on column affinity.
func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02 15:04:05", s)
}
