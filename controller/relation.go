package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/dao/mysql"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
)

type UserListResponse struct {
	models.Response
	UserList []models.UserResponse `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	id, _ := strconv.ParseUint(c.Query("to_user_id"), 10, 64)
	to_user_id := uint(id)
	actionType := c.Query("action_type")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "token有错误"})
		return
	}
	user_id := claims.ID
	if _, exist := mysql.FindUserByID(to_user_id); !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "关注用户不存在"})
		return
	}
	// 关注
	if actionType == "1" {
		if err := service.AddFollow(user_id, to_user_id); err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "关注defeat " + err.Error()})
		}
	} else {
		if err := service.DeleteFollow(user_id, to_user_id); err != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "取注defeat " + err.Error()})
		}
	}
	// 这边的思路先是提供两个token来获取用户，可以改
	//token := c.Query("token_my")
	//token2 := c.Query("token_other")
	//actionType := c.Query("action_type")
	//if user, exist := mysql.FindUserByToken(token); !exist {
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "当前用户不存在"})
	//	return
	//} else {
	//	userOther, exist := mysql.FindUserByToken(token2)
	//	if !exist {
	//		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
	//		return
	//	}
	//	switch actionType {
	//	case "1":
	//		// 关注
	//		otherId := int64(userOther.ID)
	//		userId := int64(user.ID)
	//		mysql.DB.Create(&models.Relations{
	//			UserId:      userId,
	//			FollowingId: &otherId,
	//		})
	//		mysql.DB.Create(&models.Relations{
	//			UserId:     otherId,
	//			FollowedId: &userId,
	//		})
	//		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "关注成功"})
	//	case "2":
	//		// 取消关注
	//
	//	}
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "没操作"})
	//}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "token有错误"})
		return
	}
	result, err := service.GetFollowList(claims.ID)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "GetFollowList()错误"})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: result,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "token有错误"})
		return
	}
	result, err := service.GetFollowerList(claims.ID)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "GetFollowList()错误"})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: result,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "token有错误"})
		return
	}
	result, err := service.GetFriendList(claims.ID)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "GetFollowList()错误"})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: result,
	})
}
