// Code generated by Kitex v0.6.2. DO NOT EDIT.

package videoservice

import (
	"context"
	video "douyin/kitex_gen/video"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return videoServiceServiceInfo
}

var videoServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "VideoService"
	handlerType := (*video.VideoService)(nil)
	methods := map[string]kitex.MethodInfo{
		"VideoFeed":           kitex.NewMethodInfo(videoFeedHandler, newVideoServiceVideoFeedArgs, newVideoServiceVideoFeedResult, false),
		"PublishVideo":        kitex.NewMethodInfo(publishVideoHandler, newVideoServicePublishVideoArgs, newVideoServicePublishVideoResult, false),
		"GetPublishVideoList": kitex.NewMethodInfo(getPublishVideoListHandler, newVideoServiceGetPublishVideoListArgs, newVideoServiceGetPublishVideoListResult, false),
		"GetWorkCount":        kitex.NewMethodInfo(getWorkCountHandler, newVideoServiceGetWorkCountArgs, newVideoServiceGetWorkCountResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "video",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.6.2",
		Extra:           extra,
	}
	return svcInfo
}

func videoFeedHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceVideoFeedArgs)
	realResult := result.(*video.VideoServiceVideoFeedResult)
	success, err := handler.(video.VideoService).VideoFeed(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceVideoFeedArgs() interface{} {
	return video.NewVideoServiceVideoFeedArgs()
}

func newVideoServiceVideoFeedResult() interface{} {
	return video.NewVideoServiceVideoFeedResult()
}

func publishVideoHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServicePublishVideoArgs)
	realResult := result.(*video.VideoServicePublishVideoResult)
	success, err := handler.(video.VideoService).PublishVideo(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServicePublishVideoArgs() interface{} {
	return video.NewVideoServicePublishVideoArgs()
}

func newVideoServicePublishVideoResult() interface{} {
	return video.NewVideoServicePublishVideoResult()
}

func getPublishVideoListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetPublishVideoListArgs)
	realResult := result.(*video.VideoServiceGetPublishVideoListResult)
	success, err := handler.(video.VideoService).GetPublishVideoList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newVideoServiceGetPublishVideoListArgs() interface{} {
	return video.NewVideoServiceGetPublishVideoListArgs()
}

func newVideoServiceGetPublishVideoListResult() interface{} {
	return video.NewVideoServiceGetPublishVideoListResult()
}

func getWorkCountHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*video.VideoServiceGetWorkCountArgs)
	realResult := result.(*video.VideoServiceGetWorkCountResult)
	success, err := handler.(video.VideoService).GetWorkCount(ctx, realArg.UserId)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newVideoServiceGetWorkCountArgs() interface{} {
	return video.NewVideoServiceGetWorkCountArgs()
}

func newVideoServiceGetWorkCountResult() interface{} {
	return video.NewVideoServiceGetWorkCountResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) VideoFeed(ctx context.Context, request *video.VideoFeedRequest) (r *video.VideoFeedResponse, err error) {
	var _args video.VideoServiceVideoFeedArgs
	_args.Request = request
	var _result video.VideoServiceVideoFeedResult
	if err = p.c.Call(ctx, "VideoFeed", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) PublishVideo(ctx context.Context, request *video.PublishVideoRequest) (r *video.PublishVideoResponse, err error) {
	var _args video.VideoServicePublishVideoArgs
	_args.Request = request
	var _result video.VideoServicePublishVideoResult
	if err = p.c.Call(ctx, "PublishVideo", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetPublishVideoList(ctx context.Context, request *video.PublishVideoListRequest) (r *video.PublishVideoListResponse, err error) {
	var _args video.VideoServiceGetPublishVideoListArgs
	_args.Request = request
	var _result video.VideoServiceGetPublishVideoListResult
	if err = p.c.Call(ctx, "GetPublishVideoList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetWorkCount(ctx context.Context, userId int64) (r int32, err error) {
	var _args video.VideoServiceGetWorkCountArgs
	_args.UserId = userId
	var _result video.VideoServiceGetWorkCountResult
	if err = p.c.Call(ctx, "GetWorkCount", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
