// Code generated by Kitex v0.6.2. DO NOT EDIT.

package relationservice

import (
	"context"
	relation "douyin/kitex_gen/relation"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

func serviceInfo() *kitex.ServiceInfo {
	return relationServiceServiceInfo
}

var relationServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "RelationService"
	handlerType := (*relation.RelationService)(nil)
	methods := map[string]kitex.MethodInfo{
		"RelationAction":       kitex.NewMethodInfo(relationActionHandler, newRelationServiceRelationActionArgs, newRelationServiceRelationActionResult, false),
		"GetFollowList":        kitex.NewMethodInfo(getFollowListHandler, newRelationServiceGetFollowListArgs, newRelationServiceGetFollowListResult, false),
		"GetFollowerList":      kitex.NewMethodInfo(getFollowerListHandler, newRelationServiceGetFollowerListArgs, newRelationServiceGetFollowerListResult, false),
		"GetFriendList":        kitex.NewMethodInfo(getFriendListHandler, newRelationServiceGetFriendListArgs, newRelationServiceGetFriendListResult, false),
		"GetFollowListCount":   kitex.NewMethodInfo(getFollowListCountHandler, newRelationServiceGetFollowListCountArgs, newRelationServiceGetFollowListCountResult, false),
		"GetFollowerListCount": kitex.NewMethodInfo(getFollowerListCountHandler, newRelationServiceGetFollowerListCountArgs, newRelationServiceGetFollowerListCountResult, false),
		"IsFollowing":          kitex.NewMethodInfo(isFollowingHandler, newRelationServiceIsFollowingArgs, newRelationServiceIsFollowingResult, false),
		"IsFriend":             kitex.NewMethodInfo(isFriendHandler, newRelationServiceIsFriendArgs, newRelationServiceIsFriendResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "relation",
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

func relationActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceRelationActionArgs)
	realResult := result.(*relation.RelationServiceRelationActionResult)
	success, err := handler.(relation.RelationService).RelationAction(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRelationServiceRelationActionArgs() interface{} {
	return relation.NewRelationServiceRelationActionArgs()
}

func newRelationServiceRelationActionResult() interface{} {
	return relation.NewRelationServiceRelationActionResult()
}

func getFollowListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceGetFollowListArgs)
	realResult := result.(*relation.RelationServiceGetFollowListResult)
	success, err := handler.(relation.RelationService).GetFollowList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRelationServiceGetFollowListArgs() interface{} {
	return relation.NewRelationServiceGetFollowListArgs()
}

func newRelationServiceGetFollowListResult() interface{} {
	return relation.NewRelationServiceGetFollowListResult()
}

func getFollowerListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceGetFollowerListArgs)
	realResult := result.(*relation.RelationServiceGetFollowerListResult)
	success, err := handler.(relation.RelationService).GetFollowerList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRelationServiceGetFollowerListArgs() interface{} {
	return relation.NewRelationServiceGetFollowerListArgs()
}

func newRelationServiceGetFollowerListResult() interface{} {
	return relation.NewRelationServiceGetFollowerListResult()
}

func getFriendListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceGetFriendListArgs)
	realResult := result.(*relation.RelationServiceGetFriendListResult)
	success, err := handler.(relation.RelationService).GetFriendList(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newRelationServiceGetFriendListArgs() interface{} {
	return relation.NewRelationServiceGetFriendListArgs()
}

func newRelationServiceGetFriendListResult() interface{} {
	return relation.NewRelationServiceGetFriendListResult()
}

func getFollowListCountHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceGetFollowListCountArgs)
	realResult := result.(*relation.RelationServiceGetFollowListCountResult)
	success, err := handler.(relation.RelationService).GetFollowListCount(ctx, realArg.UserId)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newRelationServiceGetFollowListCountArgs() interface{} {
	return relation.NewRelationServiceGetFollowListCountArgs()
}

func newRelationServiceGetFollowListCountResult() interface{} {
	return relation.NewRelationServiceGetFollowListCountResult()
}

func getFollowerListCountHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceGetFollowerListCountArgs)
	realResult := result.(*relation.RelationServiceGetFollowerListCountResult)
	success, err := handler.(relation.RelationService).GetFollowerListCount(ctx, realArg.UserId)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newRelationServiceGetFollowerListCountArgs() interface{} {
	return relation.NewRelationServiceGetFollowerListCountArgs()
}

func newRelationServiceGetFollowerListCountResult() interface{} {
	return relation.NewRelationServiceGetFollowerListCountResult()
}

func isFollowingHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceIsFollowingArgs)
	realResult := result.(*relation.RelationServiceIsFollowingResult)
	success, err := handler.(relation.RelationService).IsFollowing(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newRelationServiceIsFollowingArgs() interface{} {
	return relation.NewRelationServiceIsFollowingArgs()
}

func newRelationServiceIsFollowingResult() interface{} {
	return relation.NewRelationServiceIsFollowingResult()
}

func isFriendHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*relation.RelationServiceIsFriendArgs)
	realResult := result.(*relation.RelationServiceIsFriendResult)
	success, err := handler.(relation.RelationService).IsFriend(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = &success
	return nil
}
func newRelationServiceIsFriendArgs() interface{} {
	return relation.NewRelationServiceIsFriendArgs()
}

func newRelationServiceIsFriendResult() interface{} {
	return relation.NewRelationServiceIsFriendResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) RelationAction(ctx context.Context, request *relation.RelationActionRequest) (r *relation.RelationActionResponse, err error) {
	var _args relation.RelationServiceRelationActionArgs
	_args.Request = request
	var _result relation.RelationServiceRelationActionResult
	if err = p.c.Call(ctx, "RelationAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowList(ctx context.Context, request *relation.FollowListRequest) (r *relation.FollowListResponse, err error) {
	var _args relation.RelationServiceGetFollowListArgs
	_args.Request = request
	var _result relation.RelationServiceGetFollowListResult
	if err = p.c.Call(ctx, "GetFollowList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowerList(ctx context.Context, request *relation.FollowerListRequest) (r *relation.FollowerListResponse, err error) {
	var _args relation.RelationServiceGetFollowerListArgs
	_args.Request = request
	var _result relation.RelationServiceGetFollowerListResult
	if err = p.c.Call(ctx, "GetFollowerList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFriendList(ctx context.Context, request *relation.FriendListRequest) (r *relation.FriendListResponse, err error) {
	var _args relation.RelationServiceGetFriendListArgs
	_args.Request = request
	var _result relation.RelationServiceGetFriendListResult
	if err = p.c.Call(ctx, "GetFriendList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowListCount(ctx context.Context, userId int64) (r int32, err error) {
	var _args relation.RelationServiceGetFollowListCountArgs
	_args.UserId = userId
	var _result relation.RelationServiceGetFollowListCountResult
	if err = p.c.Call(ctx, "GetFollowListCount", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFollowerListCount(ctx context.Context, userId int64) (r int32, err error) {
	var _args relation.RelationServiceGetFollowerListCountArgs
	_args.UserId = userId
	var _result relation.RelationServiceGetFollowerListCountResult
	if err = p.c.Call(ctx, "GetFollowerListCount", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) IsFollowing(ctx context.Context, request *relation.IsFollowingRequest) (r bool, err error) {
	var _args relation.RelationServiceIsFollowingArgs
	_args.Request = request
	var _result relation.RelationServiceIsFollowingResult
	if err = p.c.Call(ctx, "IsFollowing", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) IsFriend(ctx context.Context, request *relation.IsFriendRequest) (r bool, err error) {
	var _args relation.RelationServiceIsFriendArgs
	_args.Request = request
	var _result relation.RelationServiceIsFriendResult
	if err = p.c.Call(ctx, "IsFriend", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
