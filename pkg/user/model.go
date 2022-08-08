package user

import "time"

// User holds the user entity fields
type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
