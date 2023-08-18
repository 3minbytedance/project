package video

import (
	"context"
	"douyin/constant"
	"douyin/kitex_gen/video"
	"douyin/kitex_gen/video/videoservice"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net/http"
	"strconv"
)

var videoClient videoservice.Client

// 限制上传文件的最大大小 200MB 最小大小1MB
const maxFileSize = 200 * 1024 * 1024
const minFileSize = 1 * 1024 * 1024

func init() {

	// Etcd 服务发现
	r, err := etcd.NewEtcdResolver([]string{constant.EtcdAddr})
	if err != nil {
		log.Fatal(err)
	}

	videoClient, err = videoservice.NewClient(
		constant.VideoServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
		// Please keep the same as provider.WithServiceName
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.VideoServiceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func FeedList(ctx context.Context, c *app.RequestContext) {
	////todo 后面再改
	//token := c.Query("token") //TODO 视频流客户端传递这个参数，用处Token续签、未登录的情况下查询关注返回false
	//var userId uint
	//userToken, err := common.ParseToken(token)
	//if err != nil {
	//	//isLogged := false
	//	userId = 0
	//} else {
	//	//isLogged = true
	//	userId = userToken.ID
	//}
	//
	//latestTime := c.Query("latest_time")
	//unixTime, err := strconv.Atoi(latestTime)
	//if latestTime == "" || latestTime == "0" {
	//	latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	//} else if err != nil || unixTime < 0 {
	//	zap.L().Error("Parse videoIdStr err:", zap.Error(err))
	//	c.JSON(http.StatusOK, err.Error())
	//	return
	//}
	//req := &video.VideoFeedRequest{
	//	UserId:     int32(userId),
	//	LatestTime: &latestTime,
	//}
	//
	//resp, err := videoClient.VideoFeed(ctx, req)
	//
	//if err != nil {
	//	zap.L().Error("Get feed list from video client err.", zap.Error(err))
	//	c.JSON(http.StatusOK, video.VideoFeedResponse{
	//		StatusCode: 1,
	//		StatusMsg:  thrift.StringPtr("Server internal error."),
	//	})
	//	return
	//}
	//c.JSON(http.StatusOK, resp)
}

func GetPublishList(ctx context.Context, c *app.RequestContext) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, video.PublishVideoListResponse{
			StatusCode: 0,
		})
		return
	}
	req := &video.PublishVideoListRequest{
		FromUserId: 0,
		ToUserId:   int32(userId),
	}
	r, err := videoClient.GetPublishVideoList(ctx, req)

	c.JSON(http.StatusOK, r)
}

func Publish(ctx context.Context, c *app.RequestContext) {
	//// TODO 待改
	//token := c.PostForm("token")
	//title := c.PostForm("title")
	//file, err := c.FormFile("data")
	//
	//if token == "" || title == "" || err != nil || file.Size == 0 {
	//	c.JSON(http.StatusBadRequest, video.VideoFeedResponse{
	//		StatusCode: 1,
	//	})
	//	return
	//}
	//userToken, _ := common.ParseToken(token)
	////userId := userToken.ID
	//// 校验文件类型
	//ext := strings.ToLower(filepath.Ext(file.Filename))
	//if !isValidFileType(ext) {
	//	msg := "无效的文件类型"
	//	zap.L().Error(msg)
	//	c.JSON(http.StatusBadRequest, video.VideoFeedResponse{
	//		StatusCode: 400,
	//		StatusMsg:  &msg})
	//	return
	//}
	//
	//// 校验文件大小
	//if file.Size > maxFileSize || file.Size < minFileSize {
	//	c.JSON(http.StatusBadRequest, video.PublishVideoResponse{
	//		StatusCode: 1,
	//	})
	//	return
	//}
	//
	//// 生成 UUID
	//fileId := strings.Replace(uuid.New().String(), "-", "", -1)
	//
	//// 修改文件名
	//videoFileName := fileId + ".mp4"
	//
	////todo IO流优化待测，先用内置的
	//err = c.SaveUploadedFile(file, "./public/"+videoFileName)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, video.PublishVideoResponse{
	//		StatusCode: 1,
	//	})
	//	return
	//}

	//c.JSON(http.StatusOK, models.Response{
	//	StatusCode: int32(CodeSuccess),
	//	StatusMsg:  codeMsgMap[CodeSuccess]})

	// MQ 异步解耦,解决返回json阻塞 TODO

	//imgName := videoservice.GetVideoCover(videoFileName)
	//service.StoreVideoAndImg(videoFileName, imgName, userId, title)

	//req := &video.PublishVideoRequest{
	//	UserId:     int32(userId),
	//	Data: file,
	//}
	//
	//resp, err := videoClient.VideoFeed(ctx, req)

	//if err != nil {
	//	zap.L().Error("Get feed list from video client err.", zap.Error(err))
	//	c.JSON(http.StatusOK, video.VideoFeedResponse{
	//		StatusCode: 1,
	//		StatusMsg:  thrift.StringPtr("Server internal error."),
	//	})
	//	return
	//}
	//c.JSON(http.StatusOK, resp)
}

// 校验文件类型是否为视频类型
func isValidFileType(fileExt string) bool {
	validExts := []string{".mp4", ".avi", ".mov"}
	for _, ext := range validExts {
		if fileExt == ext {
			return true
		}
	}
	return false
}
