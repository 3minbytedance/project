package common

import (
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
)

const ContextUserIDKey = "userId"
const TokenValid = "tokenValid"

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUserID(rc *app.RequestContext) (userID uint, err error) {
	uid, ok := rc.Get(ContextUserIDKey)
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
