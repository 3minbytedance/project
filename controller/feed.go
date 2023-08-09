package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"project/models"
//	"time"
//)
//
//type FeedResponse struct {
//	models.Response
//	VideoList []models.VideoRes `json:"video_list,omitempty"`
//	NextTime  int64             `json:"next_time,omitempty"`
//}
//
//// Feed same demo video list for every request
//func Feed(c *gin.Context) {
//	c.JSON(http.StatusOK, FeedResponse{
//		Response:  models.Response{StatusCode: 0},
//		VideoList: DemoVideos,
//		NextTime:  time.Now().Unix(),
//	})
//}
