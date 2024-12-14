package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	e "github.com/itaraxa/turbo-waddle/internal/errors"
)

const (
	MIN_SALT_LENGTH = 8
	MAX_SALT_LENGTH = 128
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

/*
generateSalt generates random bytes slice of random bytes of specified length

Args:

	size int: size of slice. In range MIN_SALT_LENGTH <= size <= MAX_SALT_LENGTH

Returns:

	[]byte: slice of random bytes
	error
*/
func generateSalt(size int) ([]byte, error) {
	if size < 8 || size > 128 {
		return nil, errors.Join(ErrGeneratingRandomSalt,
			fmt.Errorf("incorrect size = %d. Should be in range %d <= size <= %d", size, MIN_SALT_LENGTH, MAX_SALT_LENGTH),
		)
	}
	s := make([]byte, size)
	_, err := rand.Reader.Read(s)
	if err != nil {
		return nil, errors.Join(ErrGeneratingRandomSalt, err)
	}
	return s, nil
}

/*
GeneratePasswordWithSaltHash generates SHA-256 hash for <salt + password> string

Args:

	salt []byte: salt as slice of bytes
	paawsord []byte: password as slice of bytes

Returns:

	hash [32]byte: sha256 hash
	err error: nil or error if salt is too short or long, if password is empty
*/
func GeneratePasswordWithSaltHash(salt []byte, password []byte) (hash [32]byte, err error) {
	if len(salt) < MIN_SALT_LENGTH || len(salt) > MAX_SALT_LENGTH {
		return hash, errors.Join(ErrHashingPassword,
			e.ErrInternalServerError,
			fmt.Errorf("incorrect salt size = %d. Should be in range %d <= size <= %d", len(salt), MIN_SALT_LENGTH, MAX_SALT_LENGTH),
		)
	}
	if len(password) == 0 {
		return hash, errors.Join(ErrHashingPassword,
			e.ErrInvalidRequestFormat,
			fmt.Errorf("password is empty"),
		)
	}
	hash = sha256.Sum256(append(salt, password...))

	return
}

/*
CheckPassword compares generated hash with stored hash for password verification

Args:

	salt []byte
	password []byte: password received during authorization
	storedHash [32]byte: hash generated during user registration

Returns:

	result bool: true - if hashes are equal
	err error: nil or error, occured while checking the hash
*/
func CheckPassword(salt []byte, password []byte, storedHash [32]byte) (result bool, err error) {
	if len(password) == 0 {
		return false, errors.Join(ErrCheckPassword,
			e.ErrInvalidRequestFormat,
			errors.New("password is empty"),
		)
	}
	checkedHash, err := GeneratePasswordWithSaltHash(salt, password)
	if err != nil {
		return false, errors.Join(ErrCheckPassword,
			e.ErrInternalServerError,
			err,
		)
	}
	result = checkedHash == storedHash
	return
}

/*
GenerateToken64 generates 64-byte random token

Returns:

	token64 string
	err error
*/
func GenerateToken64() (token64 string, err error) {
	s := make([]byte, 64)
	_, err = rand.Reader.Read(s)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s), nil
}
