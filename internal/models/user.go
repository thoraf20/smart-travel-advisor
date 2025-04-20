package models

type User struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Name         string `db:"name"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}
