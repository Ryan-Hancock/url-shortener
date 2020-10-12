package main

import (
	"testing"
)

func GetTestDB() urlDB {
	return newDB(":memory:")
}

func TestSetAndGet(t *testing.T) {
	db := GetTestDB()

	t.Run("should not allow empty values", func(t *testing.T) {
		err := db.Set("", "")
		if err == nil {
			t.Errorf("Set() returned err: got %v want error", nil)
		}
	})

	t.Run("should be able to retrieve record", func(t *testing.T) {
		tKey := "testkey"
		tValue := "testvalue"
		db.Set(tKey, tValue)
		if value := db.Get(tKey); value != tValue {
			t.Errorf("Get() returned err: got %s want %s", value, tValue)
		}
	})
}
