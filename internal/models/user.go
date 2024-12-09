package models

type User struct {
	Username string
	Active   bool
	Password []byte
	Hash     []byte
	Salt     []byte
	token    []byte
}
