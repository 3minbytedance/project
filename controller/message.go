package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project/models"
	"project/service"
	"project/utils"
	"strconv"
)

var tempChat = map[string][]models.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	models.Response
	MessageList []models.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(c *gin.Context) {
	toUserIdStr := c.Query("to_user_id")
	toUserId, err := strconv.ParseUint(toUserIdStr, 10, 64)
	content := c.Query("content")
	fromUserId, err := utils.GetCurrentUserID(c)
	if err != nil {
		log.Println("Get user id from ctx err:", err)
		return
	}
	err = service.SendMessage(uint(fromUserId), uint(toUserId), content)
	if err != nil {
		c.JSON(http.StatusOK,
			models.Response{
				StatusCode: int32(CodeServerBusy),
				StatusMsg:  CodeServerBusy.Msg(),
			})
	}
	c.JSON(http.StatusOK,
		models.Response{
			StatusCode: int32(CodeSuccess),
			StatusMsg:  CodeSuccess.Msg(),
		})
	//

	//if user, exist := mysql.FindUserByToken(token); exist {
	//	userIdB, _ := strconv.Atoi(toUserId)
	//	chatKey := genChatKey(int64(user.ID), int64(userIdB))
	//
	//	atomic.AddInt64(&messageIdSequence, 1)
	//	curMessage := models.Message{
	//		Id:         messageIdSequence,
	//		Content:    content,
	//		CreateTime: time.Now().Format(time.Kitchen),
	//	}
	//
	//	if messages, exist := tempChat[chatKey]; exist {
	//		tempChat[chatKey] = append(messages, curMessage)
	//	} else {
	//		tempChat[chatKey] = []models.Message{curMessage}
	//	}
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// MessageChat all users have same follow list
func MessageChat(c *gin.Context) {
	toUserIdStr := c.Query("to_user_id")
	preMsgTimeStr := c.Query("pre_msg_time")
	preMsgTime, err := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if err != nil {
		return
	}
	toUserId, _ := strconv.ParseInt(toUserIdStr, 10, 64)
	fromUserId, _ := utils.GetCurrentUserID(c)

	msgList, err := service.GetMessageList(fromUserId, uint(toUserId), preMsgTime)
	if err != nil {
		c.JSON(http.StatusOK,
			models.MessageChatResponse{
				Response: models.Response{
					StatusCode: -1,
					StatusMsg:  "Found message chat failed:" + err.Error(),
				},
				MessageList: nil,
			})
		return
	}
	c.JSON(http.StatusOK,
		models.MessageChatResponse{
			Response: models.Response{
				StatusCode: 0,
				StatusMsg:  "Found comments success",
			},
			MessageList: msgList,
		})

}

// 用redis就需要，不用就可以删掉
func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
