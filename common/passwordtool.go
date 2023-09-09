package common

import (
	"github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

func MakePassword(pwd string) (string, error) {
	hash, err := argon2id.CreateHash(pwd, argon2id.DefaultParams)
	if err != nil {
		zap.L().Error("generate passwordHash fail!", zap.Error(err))
		return "", err
	}

	return hash, err
}

func CheckPassword(pwd, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(pwd, hash)
	if err != nil {
		zap.L().Error("CheckPassword fail!", zap.Error(err))
		return false
	}
	return match
}
