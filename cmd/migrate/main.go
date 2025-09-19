package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"idea-collision-engine-api/internal/database"
	"idea-collision-engine-api/pkg/config"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		printHelp()
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("🗄️  Running database migrations...")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Seed collision domains
	pgDB, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create PostgresDB: %v", err)
	}
	defer pgDB.Close()

	if err := seedCollisionDomains(pgDB); err != nil {
		log.Fatalf("Failed to seed collision domains: %v", err)
	}

	fmt.Println("✅ Database setup completed successfully!")
}

func printHelp() {
	fmt.Println(`Database Migration Utility

This utility sets up the database schema and seeds initial data for the Idea Collision Engine API.

Usage:
  ./migrate                Run all migrations and seed data
  ./migrate --help         Show this help message

Environment Variables:
  DATABASE_URL            PostgreSQL connection string (required)

Examples:
  DATABASE_URL="postgresql://user:pass@localhost/db" ./migrate
`)
}

func runMigrations(db *sql.DB) error {
	// Read migration file
	migrationPath := "migrations/001_initial_schema.sql"
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		// Try relative path from cmd/migrate
		migrationPath = "../../migrations/001_initial_schema.sql"
	}

	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
	}

	// Execute migration
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	fmt.Println("📋 Applied initial schema migration")
	return nil
}

func seedCollisionDomains(db *database.PostgresDB) error {
	// Check if domains already exist
	domains, err := db.GetCollisionDomains("basic")
	if err != nil {
		return err
	}

	if len(domains) > 0 {
		fmt.Printf("✅ Collision domains already seeded (%d domains)\n", len(domains))
		return nil
	}

	// Seed the domains
	seedDomains := database.GetCollisionDomainSeeds()
	
	fmt.Printf("🌱 Seeding %d collision domains...\n", len(seedDomains))
	
	for i, domain := range seedDomains {
		if err := db.CreateCollisionDomain(&domain); err != nil {
			return fmt.Errorf("failed to seed domain %s: %w", domain.Name, err)
		}
		
		if i%10 == 0 {
			fmt.Printf("   Seeded %d/%d domains...\n", i, len(seedDomains))
		}
	}

	fmt.Printf("✅ Successfully seeded %d collision domains\n", len(seedDomains))
	return nil
}