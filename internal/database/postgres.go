package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"idea-collision-engine-api/internal/models"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(databaseURL string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

// User operations
func (p *PostgresDB) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, subscription_tier, interests, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	interestsJSON, _ := json.Marshal(user.Interests)
	
	_, err := p.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.SubscriptionTier,
		interestsJSON,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

func (p *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	var interestsJSON []byte
	
	query := `
		SELECT id, email, password_hash, subscription_tier, interests, created_at, updated_at
		FROM users WHERE email = $1
	`
	
	err := p.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.SubscriptionTier,
		&interestsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if len(interestsJSON) > 0 {
		json.Unmarshal(interestsJSON, &user.Interests)
	}
	
	return user, nil
}

func (p *PostgresDB) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	var interestsJSON []byte
	
	query := `
		SELECT id, email, password_hash, subscription_tier, interests, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	err := p.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.SubscriptionTier,
		&interestsJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if len(interestsJSON) > 0 {
		json.Unmarshal(interestsJSON, &user.Interests)
	}
	
	return user, nil
}

// Collision Domain operations
func (p *PostgresDB) GetCollisionDomains(tier string) ([]models.CollisionDomain, error) {
	query := `
		SELECT id, name, category, description, examples, keywords, intensity, tier, created_at, updated_at
		FROM collision_domains
		WHERE tier = $1 OR tier = 'basic'
		ORDER BY name
	`
	
	rows, err := p.db.Query(query, tier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var domains []models.CollisionDomain
	for rows.Next() {
		domain := models.CollisionDomain{}
		var examplesJSON, keywordsJSON, intensityJSON []byte
		
		err := rows.Scan(
			&domain.ID,
			&domain.Name,
			&domain.Category,
			&domain.Description,
			&examplesJSON,
			&keywordsJSON,
			&intensityJSON,
			&domain.Tier,
			&domain.CreatedAt,
			&domain.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		json.Unmarshal(examplesJSON, &domain.Examples)
		json.Unmarshal(keywordsJSON, &domain.Keywords)
		json.Unmarshal(intensityJSON, &domain.Intensity)
		
		domains = append(domains, domain)
	}
	
	return domains, nil
}

func (p *PostgresDB) CreateCollisionDomain(domain *models.CollisionDomain) error {
	query := `
		INSERT INTO collision_domains (id, name, category, description, examples, keywords, intensity, tier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	examplesJSON, _ := json.Marshal(domain.Examples)
	keywordsJSON, _ := json.Marshal(domain.Keywords)
	intensityJSON, _ := json.Marshal(domain.Intensity)
	
	_, err := p.db.Exec(query,
		domain.ID,
		domain.Name,
		domain.Category,
		domain.Description,
		examplesJSON,
		keywordsJSON,
		intensityJSON,
		domain.Tier,
		domain.CreatedAt,
		domain.UpdatedAt,
	)
	
	return err
}

// Collision Session operations
func (p *PostgresDB) CreateCollisionSession(session *models.CollisionSession) error {
	query := `
		INSERT INTO collision_sessions (id, user_id, input_data, collision_result, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	inputJSON, _ := json.Marshal(session.InputData)
	resultJSON, _ := json.Marshal(session.CollisionResult)
	
	_, err := p.db.Exec(query,
		session.ID,
		session.UserID,
		inputJSON,
		resultJSON,
		session.CreatedAt,
	)
	
	return err
}

func (p *PostgresDB) GetUserCollisionHistory(userID uuid.UUID, limit int) ([]models.CollisionSession, error) {
	query := `
		SELECT id, user_id, input_data, collision_result, user_rating, exploration_notes, created_at
		FROM collision_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := p.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sessions []models.CollisionSession
	for rows.Next() {
		session := models.CollisionSession{}
		var inputJSON, resultJSON []byte
		
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&inputJSON,
			&resultJSON,
			&session.UserRating,
			&session.ExplorationNotes,
			&session.CreatedAt,
		)
		
		if err != nil {
			return nil, err
		}
		
		json.Unmarshal(inputJSON, &session.InputData)
		json.Unmarshal(resultJSON, &session.CollisionResult)
		
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

func (p *PostgresDB) RateCollision(sessionID, userID uuid.UUID, rating int, notes *string) error {
	query := `
		UPDATE collision_sessions
		SET user_rating = $1, exploration_notes = $2
		WHERE id = $3 AND user_id = $4
	`
	
	_, err := p.db.Exec(query, rating, notes, sessionID, userID)
	return err
}

// Usage tracking operations
func (p *PostgresDB) GetUserUsage(userID uuid.UUID) (*models.UserUsage, error) {
	usage := &models.UserUsage{}
	
	query := `
		SELECT id, user_id, collision_count, reset_date, created_at, updated_at
		FROM user_usage
		WHERE user_id = $1 AND reset_date >= CURRENT_DATE - INTERVAL '7 days'
		ORDER BY reset_date DESC
		LIMIT 1
	`
	
	err := p.db.QueryRow(query, userID).Scan(
		&usage.ID,
		&usage.UserID,
		&usage.CollisionCount,
		&usage.ResetDate,
		&usage.CreatedAt,
		&usage.UpdatedAt,
	)
	
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	
	if err == sql.ErrNoRows {
		// Create new usage record
		usage = &models.UserUsage{
			ID:             uuid.New(),
			UserID:         userID,
			CollisionCount: 0,
			ResetDate:      time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		
		insertQuery := `
			INSERT INTO user_usage (id, user_id, collision_count, reset_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		
		_, err = p.db.Exec(insertQuery,
			usage.ID,
			usage.UserID,
			usage.CollisionCount,
			usage.ResetDate,
			usage.CreatedAt,
			usage.UpdatedAt,
		)
		
		if err != nil {
			return nil, err
		}
	}
	
	return usage, nil
}

func (p *PostgresDB) IncrementUserUsage(userID uuid.UUID) error {
	query := `
		UPDATE user_usage
		SET collision_count = collision_count + 1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND reset_date >= CURRENT_DATE - INTERVAL '7 days'
	`
	
	_, err := p.db.Exec(query, userID)
	return err
}