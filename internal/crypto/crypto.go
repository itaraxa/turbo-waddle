package crypto

import (
	"crypto/rand"
	"crypto/sha256"
)

/*
Объект для работы с криптографией:
- генерация соли и хэша пароля
- проверка пароля
*/
type Crypt struct {
}

func NewCrypt() *Crypt {
	return &Crypt{}
}

func generateSalt(size int) ([]byte, error) {
	s := make([]byte, size)
	_, err := rand.Reader.Read(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func GeneratePasswordWithSaltHash(salt []byte, password []byte) [32]byte {
	return sha256.Sum256(append(salt, password...))
}

func CheckPassword(password string, salt []byte, storedHash [32]byte) bool {
	return GeneratePasswordWithSaltHash(salt, []byte(password)) == storedHash
}
