package hash

import (
	"github.com/ysfgrl/gerror"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, *gerror.Error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", gerror.GetError(err)
	}
	return string(hashed), nil
}

func VerifyPassword(hashed string, password string) *gerror.Error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}

func VerifyPassword2(hashed string, password string) *gerror.Error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}
