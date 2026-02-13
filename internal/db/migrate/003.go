package migrate

import (
	"fmt"

	"gorm.io/gorm"
)

func init() {
	RegisterAfterAutoMigration(Migration{
		Version: 3,
		Up:      migrateAPIKeyResetUnit,
	})
}

func migrateAPIKeyResetUnit(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	// Ensure all have a default if they are empty
	if err := db.Exec("UPDATE api_keys SET reset_unit = 'minute' WHERE reset_unit IS NULL OR reset_unit = ''").Error; err != nil {
		return fmt.Errorf("failed to set default reset_unit: %w", err)
	}

	// For existing api_keys, if reset_duration is a multiple of 86400, set reset_unit to 'day'
	if err := db.Exec("UPDATE api_keys SET reset_unit = 'day' WHERE auto_reset_quota = ? AND reset_duration > 0 AND reset_duration % 86400 = 0", true).Error; err != nil {
		return fmt.Errorf("failed to migrate api_keys reset_unit to day: %w", err)
	}

	// If it's a multiple of 3600 and still 'minute', set it to 'hour'
	if err := db.Exec("UPDATE api_keys SET reset_unit = 'hour' WHERE auto_reset_quota = ? AND reset_duration > 0 AND reset_duration % 3600 = 0 AND reset_unit = 'minute'", true).Error; err != nil {
		return fmt.Errorf("failed to migrate api_keys reset_unit to hour: %w", err)
	}

	return nil
}
