package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userId"

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUserID(c *gin.Context) (userID uint, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	fmt.Println(uid, ok)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(uint)
	fmt.Println(userID, ok)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
