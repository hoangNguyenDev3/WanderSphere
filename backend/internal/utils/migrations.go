package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hoangNguyenDev3/WanderSphere/backend/configs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MigrationManager handles database migrations
type MigrationManager struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Config *configs.PostgresConfig
}

// Migration represents a database migration
type Migration struct {
	ID        string    `gorm:"primaryKey"`
	Applied   bool      `gorm:"default:false"`
	AppliedAt time.Time `gorm:"default:null"`
	Filename  string
	UpSQL     string
	DownSQL   string
}

// MigrationStatus represents the status of migrations
type MigrationStatus struct {
	TotalMigrations   int
	AppliedMigrations int
	PendingMigrations []string
	LastApplied       string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB, cfg *configs.PostgresConfig, logger *zap.Logger) *MigrationManager {
	return &MigrationManager{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}
}

// Initialize sets up the migration tracking table
func (mm *MigrationManager) Initialize() error {
	mm.Logger.Info("Initializing migration manager")

	// Create migrations table if it doesn't exist
	if err := mm.DB.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	mm.Logger.Info("Migration manager initialized successfully")
	return nil
}

// LoadMigrationsFromDir loads migration files from a directory
func (mm *MigrationManager) LoadMigrationsFromDir(migrationDir string) ([]Migration, error) {
	mm.Logger.Info("Loading migrations from directory", zap.String("dir", migrationDir))

	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("migration directory does not exist: %s", migrationDir)
	}

	var migrations []Migration

	err := filepath.Walk(migrationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".up.sql") {
			migration, err := mm.loadMigrationFile(path)
			if err != nil {
				mm.Logger.Error("Failed to load migration file",
					zap.String("file", path),
					zap.Error(err))
				return err
			}
			migrations = append(migrations, *migration)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk migration directory: %w", err)
	}

	// Sort migrations by ID
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	mm.Logger.Info("Loaded migrations", zap.Int("count", len(migrations)))
	return migrations, nil
}

// loadMigrationFile loads a single migration file
func (mm *MigrationManager) loadMigrationFile(upFilePath string) (*Migration, error) {
	// Extract migration ID from filename (e.g., "0001_init_schema.up.sql" -> "0001")
	filename := filepath.Base(upFilePath)
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid migration filename format: %s", filename)
	}
	migrationID := parts[0]

	// Read up migration
	upContent, err := os.ReadFile(upFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read up migration file: %w", err)
	}

	// Try to read corresponding down migration
	downFilePath := strings.Replace(upFilePath, ".up.sql", ".down.sql", 1)
	var downContent []byte
	if _, err := os.Stat(downFilePath); err == nil {
		downContent, err = os.ReadFile(downFilePath)
		if err != nil {
			mm.Logger.Warn("Failed to read down migration file",
				zap.String("file", downFilePath),
				zap.Error(err))
		}
	}

	return &Migration{
		ID:       migrationID,
		Filename: filename,
		UpSQL:    string(upContent),
		DownSQL:  string(downContent),
	}, nil
}

// GetStatus returns the current migration status
func (mm *MigrationManager) GetStatus(migrationDir string) (*MigrationStatus, error) {
	migrations, err := mm.LoadMigrationsFromDir(migrationDir)
	if err != nil {
		return nil, err
	}

	// Get applied migrations from database
	var appliedMigrations []Migration
	if err := mm.DB.Where("applied = ?", true).Find(&appliedMigrations).Error; err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	var lastApplied string
	for _, migration := range appliedMigrations {
		appliedMap[migration.ID] = true
		if migration.ID > lastApplied {
			lastApplied = migration.ID
		}
	}

	var pendingMigrations []string
	for _, migration := range migrations {
		if !appliedMap[migration.ID] {
			pendingMigrations = append(pendingMigrations, migration.ID)
		}
	}

	return &MigrationStatus{
		TotalMigrations:   len(migrations),
		AppliedMigrations: len(appliedMigrations),
		PendingMigrations: pendingMigrations,
		LastApplied:       lastApplied,
	}, nil
}

// Migrate runs pending migrations
func (mm *MigrationManager) Migrate(migrationDir string) error {
	mm.Logger.Info("Starting database migration", zap.String("dir", migrationDir))

	if err := mm.Initialize(); err != nil {
		return err
	}

	migrations, err := mm.LoadMigrationsFromDir(migrationDir)
	if err != nil {
		return err
	}

	// Get already applied migrations
	appliedMap := make(map[string]bool)
	var appliedMigrations []Migration
	if err := mm.DB.Find(&appliedMigrations).Error; err == nil {
		for _, migration := range appliedMigrations {
			appliedMap[migration.ID] = true
		}
	}

	var appliedCount int
	for _, migration := range migrations {
		if appliedMap[migration.ID] {
			mm.Logger.Debug("Migration already applied", zap.String("id", migration.ID))
			continue
		}

		mm.Logger.Info("Applying migration",
			zap.String("id", migration.ID),
			zap.String("filename", migration.Filename))

		if err := mm.applyMigration(&migration); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.ID, err)
		}

		appliedCount++
	}

	mm.Logger.Info("Database migration completed",
		zap.Int("applied_migrations", appliedCount),
		zap.Int("total_migrations", len(migrations)))

	return nil
}

// applyMigration applies a single migration
func (mm *MigrationManager) applyMigration(migration *Migration) error {
	// Start transaction
	tx := mm.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			mm.Logger.Error("Migration rolled back due to panic",
				zap.String("migration", migration.ID),
				zap.Any("panic", r))
		}
	}()

	// Execute the migration SQL
	if err := tx.Exec(migration.UpSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record the migration as applied
	migration.Applied = true
	migration.AppliedAt = time.Now()

	if err := tx.Create(migration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	mm.Logger.Info("Migration applied successfully",
		zap.String("id", migration.ID),
		zap.String("filename", migration.Filename))

	return nil
}

// Rollback rolls back the last migration
func (mm *MigrationManager) Rollback() error {
	mm.Logger.Info("Rolling back last migration")

	// Get the last applied migration
	var lastMigration Migration
	if err := mm.DB.Where("applied = ?", true).Order("id DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			mm.Logger.Info("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to find last migration: %w", err)
	}

	if lastMigration.DownSQL == "" {
		return fmt.Errorf("no down migration available for %s", lastMigration.ID)
	}

	mm.Logger.Info("Rolling back migration",
		zap.String("id", lastMigration.ID),
		zap.String("filename", lastMigration.Filename))

	// Start transaction
	tx := mm.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Execute the down migration
	if err := tx.Exec(lastMigration.DownSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback SQL: %w", err)
	}

	// Remove the migration record
	if err := tx.Delete(&lastMigration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	mm.Logger.Info("Migration rolled back successfully",
		zap.String("id", lastMigration.ID))

	return nil
}

// CheckDatabaseConnection verifies database connectivity
func (mm *MigrationManager) CheckDatabaseConnection() error {
	sqlDB, err := mm.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Reset drops all tables and re-runs all migrations (DANGER: Use only in development)
func (mm *MigrationManager) Reset(migrationDir string) error {
	mm.Logger.Warn("RESETTING DATABASE - ALL DATA WILL BE LOST")

	// Get all table names
	var tables []string
	if err := mm.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	// Drop all tables
	for _, table := range tables {
		if err := mm.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			mm.Logger.Error("Failed to drop table", zap.String("table", table), zap.Error(err))
		}
	}

	mm.Logger.Info("All tables dropped")

	// Re-run migrations
	return mm.Migrate(migrationDir)
}
