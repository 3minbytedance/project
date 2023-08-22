package common

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func MakePassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("generate passwordHash fail!", zap.Error(err))
		return "", err
	}

	return string(bytes), err
}

func CheckPassword(pwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))

	return err == nil
}
