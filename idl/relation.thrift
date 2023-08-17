namespace go relation

include "user.thrift"

struct RelationActionRequest {
1: i32 user_id, // 当前登录用户
2: i32 to_user_id, // 对方用户id
3: i32 action_type, // 1-关注 2-取消关注
}

struct RelationActionResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
}

struct FollowListRequest {
1: i32 actor_id, // 当前登录用户id
2: i32 user_id, // 对方用户id
}

struct FollowListResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
3: list<user.User> user_list, // 用户信息列表
}

struct FollowerListRequest {
1: i32 actor_id, // 当前登录用户id
2: i32 user_id, // 对方用户id
}

struct FollowerListResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
3: list<user.User> user_list, // 用户列表
}

struct FollowListCountRequest {
1: i32 user_id, // 用户id
}

struct FollowListCountResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
3: i32 count, // 关注数
}

struct FollowerListCountRequest {
1: i32 user_id, // 用户id
}

struct FollowerListCountResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
3: i32 count, // 粉丝数
}

struct FriendListRequest {
1: i32 user_id, // 当前登录用户id
}

struct FriendListResponse {
1: i32 status_code, // 状态码，0-成功，其他值-失败
2: optional string status_msg, // 返回状态描述
3: list<user.User> user_list, // 用户列表
}

struct IsFollowingRequest {
1: i32 actor_id,
2: i32 user_id,
}

struct IsFollowingResponse {
1: bool result,
}

service RelationService {
RelationActionResponse RelationAction(1: RelationActionRequest Request),
FollowListResponse GetFollowList(1: FollowListRequest Request),
FollowerListResponse GetFollowerList(1: FollowerListRequest Request),
FollowListCountResponse GetFollowListCount(1: FollowListCountRequest Request),
FollowerListCountResponse GetFollowerListCount(1: FollowerListCountRequest Request),
FriendListResponse GetFriendList(1: FriendListRequest Request),
IsFollowingResponse IsFollowing(1: IsFollowingRequest Request),
}
