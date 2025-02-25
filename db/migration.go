package db

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

// CreateTables creates database tables for the given models using GORM's AutoMigrate function.
func CreateTables(db *gorm.DB, models ...interface{}) error {

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}

// DropTables drops database tables for the given models using GORM's DropTable function.
func DropTables(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		if err := db.Migrator().DropTable(model); err != nil {
			return err
		}
	}
	return nil
}

func CreateEnumTypes(db *gorm.DB, enumStatements ...string) error {
	// List of ENUM creation statements
	for _, stmt := range enumStatements {
		// Execute the statement if the type does not exist
		// PostgreSQL doesn't have an IF NOT EXISTS for CREATE TYPE, so we need to check first
		var exists bool
		checkStmt := `
			SELECT EXISTS (
				SELECT 1
				FROM pg_type
				WHERE typname = $1
			);
		`
		typeName := getTypeNameFromCreateStmt(stmt)
		if err := db.Raw(checkStmt, typeName).Scan(&exists).Error; err != nil {
			return err
		}

		if !exists {
			if err := db.Exec(stmt).Error; err != nil {
				return err
			}
			log.Printf("Created ENUM type: %s", typeName)
		} else {
			log.Printf("ENUM type already exists: %s", typeName)
		}
	}

	return nil
}

// Helper function to extract type name from CREATE TYPE statement
func getTypeNameFromCreateStmt(stmt string) string {
	var typeName string
	_, err := fmt.Sscanf(stmt, "CREATE TYPE %s AS ENUM", &typeName)
	if err != nil {
		log.Fatalf("Failed to parse ENUM type name from statement: %s", stmt)
	}
	// Remove any trailing characters like space or quotes
	typeName = strings.Trim(typeName, " '\"")
	return typeName
}
