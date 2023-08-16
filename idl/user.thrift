namespace go user

struct User {
    1: i32 id, // 用户id
    2: string name, // 用户名称
    3: i32 follow_count, // 关注总数
    4: i32 follower_count, // 粉丝总数
    5: bool is_follow, // true-已关注，false-未关注
    6: optional string avatar, // 用户头像
    7: optional string background_image, // 用户个人顶部大图
    8: optional string signature, // 个人简介
    9: optional i32 total_favorited, // 获赞数量
    10: optional i32 work_count, // 作品数量
    11: optional i32 favorite_count, // 点赞数量
}

struct UserRegisterRequest {
    1: string username, // 注册用户名，最长32个字符
    2: string password, // 密码，最长32个字符
}

struct UserRegisterResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: i32 user_id, // 用户id
    4: string token, // 用户鉴权token
}

struct UserLoginRequest {
    1: string username, // 注册用户名，最长32个字符
    2: string password, // 密码，最长32个字符
}

struct UserLoginResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: i32 user_id, // 用户id
    4: string token, // 用户鉴权token
}

struct UserInfoByIdRequest {
    1: i32 user_id, // 用户id
}

struct UserInfoByIdResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: User user, // 用户信息
}

struct UserInfoByNameRequest {
    1: string username, // 用户名
}

struct UserInfoByNameResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: optional string status_msg, // 返回状态描述
    3: i32 user_id, // 用户信息
    4: string password,
    5: string salt,
}

struct UserExistsRequest {
    1: string username,
}

struct UserExistsResponse {
    1: bool exist,
}

service UserService {
    UserRegisterResponse Register(1: UserRegisterRequest Request),
    UserLoginResponse Login(1: UserLoginRequest Request),
    UserInfoByIdResponse GetUserInfoById(1: UserInfoByIdRequest Request),
    UserInfoByNameResponse GetUserInfoByName(1: UserInfoByNameRequest Request),
    UserExistsResponse CheckUserExists(1: UserExistsRequest Request),
}
