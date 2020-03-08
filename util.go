package hecho

import "github.com/google/uuid"

// uuidGenerator generate new UUID
func uuidGenerator() string {
	return uuid.New().String()
}
