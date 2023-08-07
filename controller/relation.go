package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project/dao/mysql"
	"project/models"
)

type UserListResponse struct {
	models.Response
	UserList []models.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	// 这边的思路先是提供两个token来获取用户，可以改
	token := c.Query("token_my")
	token2 := c.Query("token_other")
	actionType := c.Query("action_type")
	if user, exist := mysql.FindUserByToken(token); !exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "当前用户不存在"})
		return
	} else {
		userOther, exist := mysql.FindUserByToken(token2)
		if !exist {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "用户不存在"})
			return
		}
		switch actionType {
		case "1":
			// 关注
			otherId := int64(userOther.ID)
			userId := int64(user.ID)
			mysql.DB.Create(&models.Relations{
				UserId:      userId,
				FollowingId: &otherId,
			})
			mysql.DB.Create(&models.Relations{
				UserId:     otherId,
				FollowedId: &userId,
			})
			c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "关注成功"})
		case "2":
			// 取消关注

		}
		c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "没操作"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}
