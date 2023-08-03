package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"os"
	"project/dao/mysql"
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
	fmt.Println("<<<<<<< username: ", username)
	if len(username) == 0 || len(username) > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "用户名不合法",
			},
		})
		return
	}

	if len(password) <= 6 || len(password) > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{
				StatusCode: 2,
				StatusMsg:  "密码不合法",
			},
		})
		return
	}

	if _, ok := models.FindUserByName(mysql.DB, username); ok {
		c.JSON(http.StatusOK, UserResponse{
			Response: models.Response{
				StatusCode: 3,
				StatusMsg:  "用户已注册",
			},
		})
		return
	}

	// token 和加密密码
	user := models.User{}
	user.Name = username

	userStates := models.UserStates{}
	userStates.Name = username
	salt := fmt.Sprintf("%06d", rand.Int())
	userStates.Salt = salt
	userStates.Password = utils.MakePassword(password, salt)
	userStates.Token = utils.MakeToken()

	mysql.DB.Create(&userStates)
	mysql.DB.Create(&user)
	fmt.Println("<<<<<<<<<id: ", user.ID)
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: models.Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserId: int64(user.ID),
		Token:  userStates.Token,
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
	user, b := models.FindUserByName(mysql.DB, username)
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

	userState, _ := models.FindUserStateByName(mysql.DB, username)
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

	userState.Token = utils.MakeToken()
	userState.LoginTime = time.Now().Unix()
	userState.IsLogOut = false
	mysql.DB.Model(&userState).Updates(models.UserStates{})

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
	//userByID, b := utils.FindUserByID(mysql.DB, id)

	// 根据token来寻找用户
	token := c.Query("token")
	//fmt.Println("token: ", token)
	userByToken, b := models.FindUserByToken(mysql.DB, token)
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
}

// UploadAvatar 上传头像（Apifox已测，不知道提供的apk里面有没有对应的接口）
func UploadAvatar(c *gin.Context) {
	token := c.Query("token")
	if user, exist := models.FindUserByToken(mysql.DB, token); exist {
		url, err := UploadPic(token, c.Request)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		mysql.DB.Model(&user).Update("avatar", url)
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
	if user, exist := models.FindUserByToken(mysql.DB, token); exist {
		url, err := UploadPic(token, c.Request)
		if err != nil {
			c.JSON(http.StatusOK, models.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
		mysql.DB.Model(&user).Update("background_image", url)
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
