package models

import (
	"fmt"
	"gorm.io/gorm"
)

// SetupModels initializes database tables based on the defined models
func SetupModels(db *gorm.DB) error {
	// Auto migrate will create or update tables according to model structures
	err := db.AutoMigrate(&User{}, &Task{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate models: %v", err)
	}

	fmt.Println("Database migration completed successfully")
	return nil
}