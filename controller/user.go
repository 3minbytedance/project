package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	// todo 用户名密码检测是否合法且不冲突
	statusCode, msg := service.CheckUserRegisterInfo(username, password)
	if statusCode != 0 {
		c.JSON(http.StatusOK, models.UserLoginResponse{
			Response: models.Response{
				StatusCode: statusCode,
				StatusMsg:  msg,
			},
		})
		return
	}

	userId, err := service.RegisterUser(username, password)
	if err != nil {
		log.Println("Register user err:", err)
		c.JSON(http.StatusOK, models.UserLoginResponse{
			Response: models.Response{
				StatusCode: -1,
				StatusMsg:  "用户注册失败" + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.UserLoginResponse{
		Response: models.Response{
			StatusCode: statusCode,
			StatusMsg:  msg,
		},
		UserId: int64(userId),
		Token:  utils.GenerateToken(userId, username),
	})

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	user, exist := service.GetUserByName(username)
	if !exist {
		// 就是用户不存在
		c.JSON(http.StatusOK, models.UserLoginResponse{
			Response: models.Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
		return
	}

	// todo 判断是否重复登录
	// 判断密码是否正确
	turePassword := utils.CheckPassword(password, user.Salt, user.Password)
	if !turePassword {
		c.JSON(http.StatusOK, models.UserLoginResponse{
			Response: models.Response{
				StatusCode: 2,
				StatusMsg:  "密码错误",
			},
		})
		return
	}

	Token := utils.GenerateToken(user.ID, username)
	c.JSON(http.StatusOK, models.UserLoginResponse{
		Response: models.Response{StatusCode: 0, StatusMsg: "登录成功"},
		UserId:   int64(user.ID),
		Token:    Token,
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

	// 根据userId来寻找用户
	userId, err := utils.GetCurrentUserID(c)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusOK, models.UserDetailResponse{
			Response: models.Response{StatusCode: 2},
		})
		return
	}
	userInfo, b := service.GetUserInfoByUserId(uint(userId))
	if b {
		c.JSON(http.StatusOK, models.UserDetailResponse{
			Response: models.Response{StatusCode: 0},
			User:     userInfo,
		})
	} else {
		c.JSON(http.StatusOK, models.UserDetailResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
	}
}

//
//// UploadAvatar 上传头像（Apifox已测，不知道提供的apk里面有没有对应的接口）
//func UploadAvatar(c *gin.Context) {
//	token := c.Query("token")
//	if user, exist := mysql.FindUserByToken(token); exist {
//		url, err := UploadPic(token, c.Request)
//		if err != nil {
//			c.JSON(http.StatusOK, models.Response{
//				StatusCode: 1,
//				StatusMsg:  err.Error(),
//			})
//			return
//		}
//		mysql.DB.Model(&user).Update("avatar", url)
//		c.JSON(http.StatusOK, models.Response{
//			StatusCode: 0,
//			StatusMsg:  "上传头像成功",
//		})
//	} else {
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
//		return
//	}
//}
//
//// UploadBackGround 上传背景（Apifox已测，不知道提供的apk里面有没有对应的接口）
//func UploadBackGround(c *gin.Context) {
//	token := c.Query("token")
//	if user, exist := mysql.FindUserByToken(token); exist {
//		url, err := UploadPic(token, c.Request)
//		if err != nil {
//			c.JSON(http.StatusOK, models.Response{
//				StatusCode: 1,
//				StatusMsg:  err.Error(),
//			})
//			return
//		}
//		mysql.DB.Model(&user).Update("background_image", url)
//		c.JSON(http.StatusOK, models.Response{
//			StatusCode: 0,
//			StatusMsg:  "上传背景成功",
//		})
//	} else {
//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
//		return
//	}
//}
//
//func UploadPic(token string, r *http.Request) (string, error) {
//	file, head, err := r.FormFile("file")
//	if err != nil {
//		return "", err
//	}
//	oldName := head.Filename
//	fileName := fmt.Sprintf("%s_%s", token, oldName)
//	url := "./public/" + fileName
//	dstFile, err := os.Create(url)
//	if err != nil {
//		return "", err
//	}
//	_, err = io.Copy(dstFile, file)
//	if err != nil {
//		return "", err
//	}
//	return url, nil
//}
//
//// UpdateUserInfo 还有个更新用户信息的
//func UpdateUserInfo(c *gin.Context) {
//
//}
