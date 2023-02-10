namespace go message
enum Code {
    Success = 1;
    ParamInvalid = 2;
    DBError = 3;
    ServerError = 4;
}
struct Message{
  1: i64 id;
  2: string content;
  3: string create_time;
}


struct MessageActionRequest{
    1: required i64 to_user_id;
    2: required i64 auth_id;
    3: required i32 action_type;
    4: required string content;
}

struct MessageActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct MessageListRequest{
    1: required i64 auth_id;
    2: required i64 to_user_id;
}

struct MessageListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Message> message_list;
}
struct Message1{
    1: i64 id
    2: i64 user_id
    3: i64 target_id
    4: string content
    5: i64 create_time
    6: i64 created_at
    7: i64 updated_at
}
/*
    query_type=1  根据id查询
    query_type=2  根据user_id查询
    query_type=3  根据target_id查询
    query_type=4  根据user_id和target_id查询
   如果给定了根据create_time也需要根据这个查询，默认-1不给定
*/
struct GetMessageRequest {
    1: i64 id
    2: i64 user_id
    3: i64 target_id
    4: i64 create_time // 获取消息需要的
    5: i64 query_type

}

struct GetMessageResponse{
    1: list<Message1> messages
}

service MessageService {
    MessageActionResponse MessageAction(1: MessageActionRequest req);
    MessageListResponse MessageList(1: MessageListRequest req);
    GetMessageResponse GetMessage(1: GetMessageRequest req)
}