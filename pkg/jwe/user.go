package jwe

type UserRepository interface {
	FindByUsername(username string) (User, error)
}

type User interface {
	ClearPassword()
	ComparePasswords(pwd string) bool
}
