package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"idea-collision-engine-api/internal/models"
)

type PostgresTestSuite struct {
	suite.Suite
	db     *sql.DB
	mock   sqlmock.Sqlmock
	pgdb   *PostgresDB
}

func (suite *PostgresTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New()
	assert.NoError(suite.T(), err)
	
	suite.pgdb = &PostgresDB{db: suite.db}
}

func (suite *PostgresTestSuite) TearDownTest() {
	suite.db.Close()
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

func (suite *PostgresTestSuite) TestCreateUser() {
	user := &models.User{
		ID:               uuid.New(),
		Email:            "test@example.com",
		PasswordHash:     "$2a$10$hashedpassword",
		SubscriptionTier: models.TierFree,
		Interests:        []string{"technology", "design"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	suite.mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.ID,
			user.Email,
			user.PasswordHash,
			user.SubscriptionTier,
			sqlmock.AnyArg(), // JSON interests
			user.CreatedAt,
			user.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := suite.pgdb.CreateUser(user)
	assert.NoError(suite.T(), err)
}

func (suite *PostgresTestSuite) TestGetUserByEmail() {
	userID := uuid.New()
	email := "test@example.com"
	interests := `["technology", "design"]`
	
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "subscription_tier", 
		"interests", "created_at", "updated_at",
	}).AddRow(
		userID,
		email,
		"$2a$10$hashedpassword",
		models.TierFree,
		interests,
		time.Now(),
		time.Now(),
	)
	
	suite.mock.ExpectQuery("SELECT .* FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnRows(rows)
	
	user, err := suite.pgdb.GetUserByEmail(email)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.ID)
	assert.Equal(suite.T(), email, user.Email)
	assert.Equal(suite.T(), models.TierFree, user.SubscriptionTier)
	assert.Equal(suite.T(), []string{"technology", "design"}, user.Interests)
}

func (suite *PostgresTestSuite) TestGetUserByEmailNotFound() {
	email := "nonexistent@example.com"
	
	suite.mock.ExpectQuery("SELECT .* FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)
	
	user, err := suite.pgdb.GetUserByEmail(email)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrNoRows, err)
	assert.Nil(suite.T(), user)
}

func (suite *PostgresTestSuite) TestGetUserByID() {
	userID := uuid.New()
	interests := `["technology", "design"]`
	
	rows := sqlmock.NewRows([]string{
		"id", "email", "password_hash", "subscription_tier",
		"interests", "created_at", "updated_at",
	}).AddRow(
		userID,
		"test@example.com",
		"$2a$10$hashedpassword",
		models.TierPro,
		interests,
		time.Now(),
		time.Now(),
	)
	
	suite.mock.ExpectQuery("SELECT .* FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)
	
	user, err := suite.pgdb.GetUserByID(userID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.ID)
	assert.Equal(suite.T(), models.TierPro, user.SubscriptionTier)
	assert.Equal(suite.T(), []string{"technology", "design"}, user.Interests)
}

func (suite *PostgresTestSuite) TestGetCollisionDomains() {
	domainID := uuid.New().String()
	tier := "basic"
	
	rows := sqlmock.NewRows([]string{
		"id", "name", "category", "description", "examples", 
		"keywords", "intensity", "tier", "created_at", "updated_at",
	}).AddRow(
		domainID,
		"Biomimicry",
		"Nature",
		"How nature solves problems",
		`["Velcro from burrs", "Bullet train from birds"]`,
		`["evolution", "adaptation"]`,
		`["gentle", "moderate"]`,
		"basic",
		time.Now(),
		time.Now(),
	)
	
	suite.mock.ExpectQuery("SELECT .* FROM collision_domains").
		WithArgs(tier).
		WillReturnRows(rows)
	
	domains, err := suite.pgdb.GetCollisionDomains(tier)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), domains, 1)
	assert.Equal(suite.T(), domainID, domains[0].ID)
	assert.Equal(suite.T(), "Biomimicry", domains[0].Name)
	assert.Equal(suite.T(), []string{"Velcro from burrs", "Bullet train from birds"}, domains[0].Examples)
	assert.Equal(suite.T(), []string{"evolution", "adaptation"}, domains[0].Keywords)
	assert.Equal(suite.T(), []string{"gentle", "moderate"}, domains[0].Intensity)
}

func (suite *PostgresTestSuite) TestCreateCollisionDomain() {
	domain := &models.CollisionDomain{
		ID:          uuid.New().String(),
		Name:        "Test Domain",
		Category:    "Test Category",
		Description: "Test description",
		Examples:    []string{"example1", "example2"},
		Keywords:    []string{"keyword1", "keyword2"},
		Intensity:   []string{"gentle", "moderate"},
		Tier:        "basic",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	suite.mock.ExpectExec("INSERT INTO collision_domains").
		WithArgs(
			domain.ID,
			domain.Name,
			domain.Category,
			domain.Description,
			sqlmock.AnyArg(), // JSON examples
			sqlmock.AnyArg(), // JSON keywords
			sqlmock.AnyArg(), // JSON intensity
			domain.Tier,
			domain.CreatedAt,
			domain.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := suite.pgdb.CreateCollisionDomain(domain)
	assert.NoError(suite.T(), err)
}

func (suite *PostgresTestSuite) TestCreateCollisionSession() {
	session := &models.CollisionSession{
		ID:     uuid.New(),
		UserID: uuid.New(),
		InputData: models.CollisionInput{
			UserInterests:      []string{"tech", "design"},
			CurrentProject:     "mobile app",
			ProjectType:        "product",
			CollisionIntensity: "moderate",
		},
		CollisionResult: models.CollisionResult{
			ID:              uuid.New().String(),
			PrimaryDomain:   "Technology",
			CollisionDomain: "Jazz",
			Connection:      "Test connection",
			QualityScore:    85.5,
		},
		CreatedAt: time.Now(),
	}
	
	suite.mock.ExpectExec("INSERT INTO collision_sessions").
		WithArgs(
			session.ID,
			session.UserID,
			sqlmock.AnyArg(), // JSON input_data
			sqlmock.AnyArg(), // JSON collision_result
			session.CreatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := suite.pgdb.CreateCollisionSession(session)
	assert.NoError(suite.T(), err)
}

func (suite *PostgresTestSuite) TestGetUserCollisionHistory() {
	userID := uuid.New()
	sessionID := uuid.New()
	limit := 10
	
	inputData := `{"user_interests":["tech"],"current_project":"app","project_type":"product","collision_intensity":"moderate"}`
	resultData := `{"id":"123","primary_domain":"Tech","collision_domain":"Jazz","connection":"Test","quality_score":85.5,"timestamp":"2024-01-01T00:00:00Z","spark_questions":[],"examples":[],"next_steps":[]}`
	
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "input_data", "collision_result",
		"user_rating", "exploration_notes", "created_at",
	}).AddRow(
		sessionID,
		userID,
		inputData,
		resultData,
		nil,
		nil,
		time.Now(),
	)
	
	suite.mock.ExpectQuery("SELECT .* FROM collision_sessions").
		WithArgs(userID, limit).
		WillReturnRows(rows)
	
	sessions, err := suite.pgdb.GetUserCollisionHistory(userID, limit)
	
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), sessions, 1)
	assert.Equal(suite.T(), sessionID, sessions[0].ID)
	assert.Equal(suite.T(), userID, sessions[0].UserID)
	assert.Equal(suite.T(), []string{"tech"}, sessions[0].InputData.UserInterests)
	assert.Equal(suite.T(), "Jazz", sessions[0].CollisionResult.CollisionDomain)
}

func (suite *PostgresTestSuite) TestRateCollision() {
	sessionID := uuid.New()
	userID := uuid.New()
	rating := 5
	notes := "Great collision!"
	
	suite.mock.ExpectExec("UPDATE collision_sessions").
		WithArgs(rating, &notes, sessionID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	err := suite.pgdb.RateCollision(sessionID, userID, rating, &notes)
	assert.NoError(suite.T(), err)
}

func (suite *PostgresTestSuite) TestGetUserUsage() {
	userID := uuid.New()
	usageID := uuid.New()
	
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "collision_count", "reset_date",
		"created_at", "updated_at",
	}).AddRow(
		usageID,
		userID,
		3,
		time.Now().Format("2006-01-02"),
		time.Now(),
		time.Now(),
	)
	
	suite.mock.ExpectQuery("SELECT .* FROM user_usage").
		WithArgs(userID).
		WillReturnRows(rows)
	
	usage, err := suite.pgdb.GetUserUsage(userID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), usage)
	assert.Equal(suite.T(), usageID, usage.ID)
	assert.Equal(suite.T(), userID, usage.UserID)
	assert.Equal(suite.T(), 3, usage.CollisionCount)
}

func (suite *PostgresTestSuite) TestGetUserUsageNotFound() {
	userID := uuid.New()
	
	// First query returns no rows (user has no usage record)
	suite.mock.ExpectQuery("SELECT .* FROM user_usage").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
	
	// Should create new usage record
	suite.mock.ExpectExec("INSERT INTO user_usage").
		WithArgs(sqlmock.AnyArg(), userID, 0, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	usage, err := suite.pgdb.GetUserUsage(userID)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), usage)
	assert.Equal(suite.T(), userID, usage.UserID)
	assert.Equal(suite.T(), 0, usage.CollisionCount)
}

func (suite *PostgresTestSuite) TestIncrementUserUsage() {
	userID := uuid.New()
	
	suite.mock.ExpectExec("UPDATE user_usage").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	
	err := suite.pgdb.IncrementUserUsage(userID)
	assert.NoError(suite.T(), err)
}

func (suite *PostgresTestSuite) TestDatabaseConnectionError() {
	// Test error handling when database operations fail
	email := "test@example.com"
	
	suite.mock.ExpectQuery("SELECT .* FROM users WHERE email = \\$1").
		WithArgs(email).
		WillReturnError(sql.ErrConnDone)
	
	user, err := suite.pgdb.GetUserByEmail(email)
	
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), sql.ErrConnDone, err)
	assert.Nil(suite.T(), user)
}

func (suite *PostgresTestSuite) TestJSONMarshaling() {
	// Test that JSON marshaling works correctly for interests
	user := &models.User{
		ID:               uuid.New(),
		Email:            "json@example.com",
		PasswordHash:     "hash",
		SubscriptionTier: models.TierPro,
		Interests:        []string{"complex interest with spaces", "unicode: 🚀"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	suite.mock.ExpectExec("INSERT INTO users").
		WithArgs(
			user.ID,
			user.Email,
			user.PasswordHash,
			user.SubscriptionTier,
			sqlmock.AnyArg(),
			user.CreatedAt,
			user.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	
	err := suite.pgdb.CreateUser(user)
	assert.NoError(suite.T(), err)
}

// Benchmark tests for database operations
func BenchmarkCreateUser(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	
	pgdb := &PostgresDB{db: db}
	
	user := &models.User{
		ID:               uuid.New(),
		Email:            "benchmark@example.com",
		PasswordHash:     "$2a$10$hashedpassword",
		SubscriptionTier: models.TierFree,
		Interests:        []string{"technology"},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	
	// Setup expectations for all benchmark iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO users").
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := pgdb.CreateUser(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}