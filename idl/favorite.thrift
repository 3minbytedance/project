namespace go favorite

include "video.thrift"

struct FavoriteActionRequest {
    1: i32 user_id, // 用户id
    2: i32 video_id, // 视频id
    3: i32 action_type, // 1-点赞，2-取消点赞
}

struct FavoriteActionResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
}

struct FavoriteListRequest {
    1: i32 user_id, // 用户id
}

struct FavoriteListResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: list<video.Video> video_list, // 用户点赞视频列表
}

struct VideoFavoriteCountRequest {
    1: i32 video_id, // 视频id
}

struct VideoFavoriteCountResponse {
    1: i32 status_code,
    2: optional string status_msg,
    3: i32 count, // 点赞数
}

struct UserFavoriteCountRequest {
    1: i32 user_id,
}

struct UserFavoriteCountResponse {
    1: i32 status_code,
    2: optional string status_msg,
    3: i32 count, // 点赞数
}

service FavoriteService {
    FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest Request),
    FavoriteListResponse GetFavoriteList(1: FavoriteListRequest Request),
    VideoFavoriteCountResponse GetVideoFavoriteCount(1: VideoFavoriteCountRequest Request),
    UserFavoriteCountResponse GetUserFavoriteCount(1: UserFavoriteCountRequest Request),
}
