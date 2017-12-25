package models

// Model is the base model all db tables are based off
type Model struct {
	// ID is the unique identifier value of the model
	ID uint `gorm:"primary_key"`
}
