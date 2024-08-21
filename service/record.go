package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/rainbowmga/timetravel/entity"
)

// Define custom error types for better error handling
var (
	ErrRecordDoesNotExist  = errors.New("record with that id does not exist")
	ErrRecordIDInvalid     = errors.New("record id must be >= 0")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrVersionNotFound     = errors.New("version not found")
)

// RecordService defines the interface for record operations
type RecordService interface {
	GetRecord(ctx context.Context, id int) (entity.Record, error)
	GetRecordVersion(ctx context.Context, id, version int) (entity.Record, error)
	CreateRecord(ctx context.Context, record entity.Record) error
	UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error)
	GetRecordVersions(ctx context.Context, id int) ([]int, error)
}

// SQLiteRecordService implements RecordService using SQLite
type SQLiteRecordService struct {
	db *sql.DB
}

// NewSQLiteRecordService creates a new SQLiteRecordService
func NewSQLiteRecordService(dbPath string) (*SQLiteRecordService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Ensure table exists
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteRecordService{db: db}, nil
}

// createTable ensures the required table exists in the database
func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER,
			version INTEGER,
			data TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id, version)
		)
	`)
	return err
}

// GetRecord retrieves the latest version of a record by its ID
func (s *SQLiteRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	var record entity.Record
	var dataJSON string

	err := s.db.QueryRowContext(ctx, `
		SELECT id, version, data, created_at, updated_at 
		FROM records 
		WHERE id = ? 
		ORDER BY version DESC 
		LIMIT 1
	`, id).Scan(&record.ID, &record.Version, &dataJSON, &record.CreatedAt, &record.UpdatedAt)

	if err == sql.ErrNoRows {
		return entity.Record{}, ErrRecordDoesNotExist
	} else if err != nil {
		return entity.Record{}, fmt.Errorf("failed to get record: %w", err)
	}

	if err := json.Unmarshal([]byte(dataJSON), &record.Data); err != nil {
		return entity.Record{}, fmt.Errorf("failed to unmarshal record data: %w", err)
	}

	return record, nil
}

// GetRecordVersion retrieves a specific version of a record
func (s *SQLiteRecordService) GetRecordVersion(ctx context.Context, id, version int) (entity.Record, error) {
	var record entity.Record
	var dataJSON string

	err := s.db.QueryRowContext(ctx, `
		SELECT id, version, data, created_at, updated_at 
		FROM records 
		WHERE id = ? AND version = ?
	`, id, version).Scan(&record.ID, &record.Version, &dataJSON, &record.CreatedAt, &record.UpdatedAt)

	if err == sql.ErrNoRows {
		return entity.Record{}, ErrVersionNotFound
	} else if err != nil {
		return entity.Record{}, fmt.Errorf("failed to get record version: %w", err)
	}

	if err := json.Unmarshal([]byte(dataJSON), &record.Data); err != nil {
		return entity.Record{}, fmt.Errorf("failed to unmarshal record data: %w", err)
	}

	return record, nil
}

// CreateRecord inserts a new record into the database
func (s *SQLiteRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	if record.ID <= 0 {
		return ErrRecordIDInvalid
	}

	dataJSON, err := json.Marshal(record.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal record data: %w", err)
	}

	now := time.Now()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO records (id, version, data, created_at, updated_at) 
		VALUES (?, 1, ?, ?, ?)
	`, record.ID, string(dataJSON), now, now)
	if err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}

	return nil
}

// UpdateRecord modifies an existing record in the database
func (s *SQLiteRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	record, err := s.GetRecord(ctx, id)
	if err != nil {
		return entity.Record{}, err
	}

	for key, value := range updates {
		if value == nil {
			delete(record.Data, key)
		} else {
			record.Data[key] = *value
		}
	}

	dataJSON, err := json.Marshal(record.Data)
	if err != nil {
		return entity.Record{}, fmt.Errorf("failed to marshal updated record data: %w", err)
	}

	now := time.Now()
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO records (id, version, data, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`, id, record.Version+1, string(dataJSON), record.CreatedAt, now)
	if err != nil {
		return entity.Record{}, fmt.Errorf("failed to update record: %w", err)
	}

	record.Version++
	record.UpdatedAt = now
	return record, nil
}

// GetRecordVersions retrieves all version numbers for a given record ID
func (s *SQLiteRecordService) GetRecordVersions(ctx context.Context, id int) ([]int, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT version FROM records WHERE id = ? ORDER BY version", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get record versions: %w", err)
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan version: %w", err)
		}
		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over versions: %w", err)
	}

	if len(versions) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	return versions, nil
}
