package model

// Migrate returns the models that need to be migrated.
func Migrate() []interface{} {
	return []interface{}{
		&User{},
		&JWT{},
	}
}
