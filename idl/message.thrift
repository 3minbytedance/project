namespace go message

struct Message {
    1: i64 id, // 消息id
    2: i64 to_user_id, // 该消息接收者的id
    3: i64 from_user_id, // 该消息发送者的id
    4: string content, // 消息内容
    5: i64 create_time, // 消息创建时间
}

struct MessageChatRequest {
    1: i64 from_user_id, // 发送方id
    2: i64 to_user_id, // 对方用户id
    3: i64 pre_msg_time, // 上次最新消息的时间（新增字段-apk更新中）
}

struct MessageChatResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: string status_msg, // 返回状态描述
    3: list<Message> message_list, // 消息列表
}

struct MessageActionRequest {
    1: i64 from_user_id, // 用户id
    2: i64 to_user_id, // 对方用户id
    3: i32 action_type, // 1-发送消息
    4: string content, // 消息内容
}

struct MessageActionResponse {
    1: i32 status_code, // 状态码，0-成功，其他值-失败
    2: string status_msg, // 返回状态描述
}

service MessageService {
    MessageChatResponse MessageChat(1: MessageChatRequest Request),
    MessageActionResponse MessageAction(1: MessageActionRequest Request),
}
