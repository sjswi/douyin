namespace go video

enum Code {
    Success = 1;
    ParamInvalid = 2;
    DBError = 3;
    ServerError = 4;
}
struct Video {
    1: i64 id;
    2: User author;
    3: string play_url;
    4: string cover_url;
    5: i64 favorite_count;
    6: i64 comment_count;
    7: bool is_favorite;
    8: string title;
}
struct User {
    1: i64 id;
    2: string name;
    3: i64 follow_count;
    4: i64 follower_count;
    5: bool is_follow;
}

struct FeedRequest{
    1: i64 auth_id;
    2: i64 latest_time;
}

struct FeedResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
    4: i64 next_time;
}
struct PublishActionRequest{
    1: required binary data; //注意thrift没有定义multipart.FILE
    2: required i64 auth_id;
    3: required string title;
    4: required string filename;
}

struct PublishActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct PublishListRequest{
    1: required i64 auth_id;
    2: required i64 user_id;
}

struct PublishListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
}
struct Video1 {
    1: i64 id
    2: i64 author_id
    3: string title
    4: string play_url
    5: string cover_url
    6: i64 created_at
    7: i64 updated_at
}
/**
* query_type=1 根据视频id查询
* query_type=2 根据作者id查询
*
**/

struct GetVideoRequest {
    1: i64 video_id
    2: i64 author_id
    3: i32 query_type
}
struct GetVideoResponse {
    1: list<Video1> video
}
service FeedService {
    FeedResponse Feed(1: FeedRequest req);
    PublishActionResponse PublishAction(1: PublishActionRequest req);
    PublishListResponse PublishList(1: PublishListRequest req);
    GetVideoResponse GetVideo(1: GetVideoRequest req);
}