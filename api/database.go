package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func initDatabase() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost/carbonapi?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("Database connection warning: %v (using in-memory fallback)", err)
		// In production, this would fail. For demo, we continue with limited functionality
	}

	// Create tables if they don't exist
	createTables(db)

	return db
}

func initRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password
		DB:       0,  // default DB
	})

	return rdb
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS emission_factors (
			id SERIAL PRIMARY KEY,
			activity VARCHAR(100) NOT NULL,
			transport_mode VARCHAR(50),
			factor DECIMAL(10,6) NOT NULL,
			unit VARCHAR(20) NOT NULL,
			source VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS calculations (
			id SERIAL PRIMARY KEY,
			activity VARCHAR(100) NOT NULL,
			input_data JSONB NOT NULL,
			carbon_footprint DECIMAL(10,6) NOT NULL,
			unit VARCHAR(20) NOT NULL,
			user_id VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS api_usage (
			id SERIAL PRIMARY KEY,
			endpoint VARCHAR(100) NOT NULL,
			user_id VARCHAR(100),
			response_time_ms INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			log.Printf("Failed to create table: %v", err)
		}
	}

	// Insert sample emission factors
	insertSampleData(db)
}

func insertSampleData(db *sql.DB) {
	// Check if data already exists
	var count int
	db.QueryRow("SELECT COUNT(*) FROM emission_factors").Scan(&count)
	if count > 0 {
		return // Data already exists
	}

	// Insert sample emission factors (kg CO2e per unit)
	sampleFactors := []struct {
		activity      string
		transportMode string
		factor        float64
		unit          string
		source        string
	}{
		{"shipping", "air", 0.996, "kg_co2e_per_tonne_km", "IPCC 2023"},
		{"shipping", "sea", 0.015, "kg_co2e_per_tonne_km", "IPCC 2023"},
		{"shipping", "road", 0.209, "kg_co2e_per_tonne_km", "IPCC 2023"},
		{"shipping", "rail", 0.028, "kg_co2e_per_tonne_km", "IPCC 2023"},
		{"electricity", "grid", 0.525, "kg_co2e_per_kwh", "IEA 2023"},
		{"electricity", "solar", 0.041, "kg_co2e_per_kwh", "IPCC 2023"},
		{"electricity", "wind", 0.011, "kg_co2e_per_kwh", "IPCC 2023"},
		{"fuel", "gasoline", 2.31, "kg_co2e_per_liter", "EPA 2023"},
		{"fuel", "diesel", 2.68, "kg_co2e_per_liter", "EPA 2023"},
		{"fuel", "natural_gas", 0.202, "kg_co2e_per_kwh", "EPA 2023"},
	}

	for _, factor := range sampleFactors {
		_, err := db.Exec(`
			INSERT INTO emission_factors (activity, transport_mode, factor, unit, source)
			VALUES ($1, $2, $3, $4, $5)
		`, factor.activity, factor.transportMode, factor.factor, factor.unit, factor.source)

		if err != nil {
			log.Printf("Failed to insert sample data: %v", err)
		}
	}

	log.Println("✅ Sample emission factors inserted successfully")
}
