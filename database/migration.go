package database

// CreateTables creates database tables for the given models using GORM's AutoMigrate function.
func CreateTables(models ...interface{}) error {

	for _, model := range models {
		if err := DefaultDataSource.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}

// DropTables drops database tables for the given models using GORM's DropTable function.
func DropTables(models ...interface{}) error {
	for _, model := range models {
		if err := DefaultDataSource.Migrator().DropTable(model); err != nil {
			return err
		}
	}
	return nil
}
