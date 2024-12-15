package crypto

import "errors"

var (
	ErrGeneratingRandomSalt  = errors.New("generateSalt: error generating random salt")
	ErrHashingPassword       = errors.New("GeneratePasswordWithSaltHash: error in password hash generation")
	ErrCheckPassword         = errors.New("CheckPassword: error in password check")
	ErrJWTSignEmptyLogin     = errors.New("CreateJWT: empty login")
	ErrJWTSignEmptySecretKey = errors.New("CreateJWT: empty secret key")
)
