package user

type Repository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}