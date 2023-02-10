namespace go comment

enum Code {
    Success = 1;
    ParamInvalid = 2;
    DBError = 3;
    ServerError = 4;
}
struct User {
    1: i64 id;
    2: string name;
    3: i64 follow_count;
    4: i64 follower_count;
    5: bool is_follow;
}

struct Comment {
  1: i64 id;
  2: User user;
  3: string content;
  4: string create_date;
}
struct CommentActionRequest{
    1: required i64 video_id;
    2: required i64 auth_id;
    3: required i32 action_type;
    4: string comment_text;
    5: i64 comment_id;
}

struct CommentActionResponse{
    1: Code status_code;
    2: string status_msg;
    3: Comment comment;
}
struct CommentListRequest{
    1: required i64 auth_id;
    2: required i64 video_id;
}

struct CommentListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Comment> comment_list;
}

struct Comment1{
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: i64 created_at
    5: i64 updated_at
    6: string content
}
/*
    query_type=1  根据id查询
    query_type=2  根据user_id查询
    query_type=3  根据video_id查询
    query_type=4  根据video_id和user_id查询
*/
struct GetCommentRequest {
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: i64 query_type

}
struct GetCommentCountRequest {
    1: i64 video_id
}
struct GetCommentCountResponse{
    1: i64 count
}
struct GetCommentResponse{
    1: list<Comment1> comments
}
service CommentService {
    CommentActionResponse CommentAction(1: CommentActionRequest req);
    CommentListResponse CommentList(1: CommentListRequest req);
    GetCommentResponse GetComment(1: GetCommentRequest req);
    GetCommentCountResponse GetCommentCount(1: GetCommentCountRequest req)
}
