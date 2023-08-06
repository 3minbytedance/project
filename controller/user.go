package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"project/models"
	"project/utils"
	"time"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
//// test data: username=zhanglei, password=douyin
//var usersLoginInfo = map[string]models.User{
//	"zhangleidouyin": {
//		Name:          "zhanglei",
//		FollowCount:   10,
//		FollowerCount: 5,
//		IsFollow:      true,
//	},
//}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	models.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	models.Response
	User models.User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// todo 用户名密码检测是否合法且不冲突
	statusCode, msg := models.CheckUserRegisterInfo(username, password)
	if statusCode != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{
				StatusCode: statusCode,
				StatusMsg:  msg,
			},
		})
		return
	}

	// todo token 和加密密码
	statusCode, msg, userId, token := models.RegisterUserInfo(username, password)

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{
			StatusCode: statusCode,
			StatusMsg:  msg,
		},
		UserId: userId,
		Token:  token,
	})

	//token := username + password
	//
	//if _, exist := usersLoginInfo[token]; exist {
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
	//	})
	//} else {
	//	atomic.AddInt64(&userIdSequence, 1)
	//	newUser := User{
	//		Id:   userIdSequence,
	//		Name: username,
	//	}
	//	usersLoginInfo[token] = newUser
	//	c.JSON(http.StatusOK, UserLoginResponse{
	//		Response: Response{StatusCode: 0},
	//		UserId:   userIdSequence,
	//		Token:    username + password,
	//	})
	//}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	fmt.Println("<<<<<<<<username: ", username)

	user, b := models.FindUserByName(utils.DB, username)
	if !b {
		// 就是用户不存在
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
		return
	}

	// todo 判断是否重复登录
	userState, _ := models.FindUserStateByName(utils.DB, username)
	// 判断密码是否正确
	checkPassword := utils.CheckPassword(password, userState.Salt, userState.Password)
	if !checkPassword {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{
				StatusCode: 2,
				StatusMsg:  "密码错误",
			},
		})
		return
	}

	userState.Token = utils.GenerateToken(int64(user.ID), username)
	userState.LoginTime = time.Now().Unix()
	userState.IsLogOut = false
	utils.DB.Model(&userState).Updates(models.UserStates{})

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{StatusCode: 0, StatusMsg: "登录成功"},
		UserId:   int64(user.ID),
		Token:    userState.Token,
	})
}

func UserInfo(c *gin.Context) {
	// 根据user_id来寻找用户
	//idStr := c.Query("user_id")
	//id, err := strconv.Atoi(idStr)
	//if err != nil {
	//	println(err)
	//}
	//userByID, b := utils.FindUserByID(utils.DB, id)

	// 根据token来寻找用户
	token := c.Query("token")
	//fmt.Println("token: ", token)
	userByToken, b := models.FindUserByToken(utils.DB, token)
	if b {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 0},
			User:     userByToken,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
	// todo 返回用户的个人信息
}

// UploadAvatar 上传头像（Apifox已测，不知道提供的apk里面有没有对应的接口）
func UploadAvatar(c *gin.Context) {
	token := c.Query("token")
	if user, exist := models.FindUserByToken(utils.DB, token); exist {
		url, err := UploadPic(token, c.Request)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		utils.DB.Model(&user).Update("avatar", url)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 0,
			StatusMsg:  "上传头像成功",
		})
	} else {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
		return
	}
}

// UploadBackGround 上传背景（Apifox已测，不知道提供的apk里面有没有对应的接口）
func UploadBackGround(c *gin.Context) {
	token := c.Query("token")
	if user, exist := models.FindUserByToken(utils.DB, token); exist {
		url, err := UploadPic(token, c.Request)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		utils.DB.Model(&user).Update("background_image", url)
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 0,
			StatusMsg:  "上传背景成功",
		})
	} else {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
		return
	}
}

func UploadPic(token string, r *http.Request) (string, error) {
	file, head, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	oldName := head.Filename
	fileName := fmt.Sprintf("%s_%s", token, oldName)
	url := "./public/" + fileName
	dstFile, err := os.Create(url)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dstFile, file)
	if err != nil {
		return "", err
	}
	return url, nil
}

// UpdateUserInfo 还有个更新用户信息的
func UpdateUserInfo(c *gin.Context) {

}
