package dstore

import (
	"fmt"
)

// SaveIfNot will query for a model and save it to the database only if it does
// not exist. Returns the saved, or queried model. Along with an error if one
// occurs, or nil on success.
func SaveIfNot(m interface{}, res interface{}) (interface{}, error) {
	// Make db
	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	// Query
	if err := db.Where(m).First(res).Error; err != nil {
		return nil, fmt.Errorf("error querying for model: %s", err.Error())
	}

	// Check
	if res == nil {
		// Save
		if err := db.Save(m).Error; err != nil {
			return nil, fmt.Errorf("error saving model: %s", err.Error())
		}

		return m, nil
	}

	return m, nil
}