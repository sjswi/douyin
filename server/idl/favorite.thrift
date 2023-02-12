namespace go favorite

enum Code {
    Success = 0;
    ParamInvalid = 1;
    DBError = 2;
    ServerError = 3;
    ErrorRequest = 4;
}
struct User {
    1: string id;
    2: string name;
    3: i64 follow_count;
    4: i64 follower_count;
    5: bool is_follow;
}

struct Video {
    1: string id;
    2: User author;
    3: string play_url;
    4: string cover_url;
    5: i64 favorite_count;
    6: i64 comment_count;
    7: bool is_favorite;
    8: string title;
}

struct FavoriteActionRequest{
    1: required i64 video_id;
    2: required i64 auth_id;
    3: required i32 action_type;
}

struct FavoriteActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct FavoriteListRequest{
    1: required i64 auth_id;
    2: required i64 user_id;
}

struct FavoriteListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
}

struct Favorite1{
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: i64 created_at
    5: i64 updated_at
}
/*
    query_type=1  根据id查询
    query_type=2  根据user_id查询
    query_type=3  根据video_id查询
    query_type=4  根据video_id和user_id查询
*/
struct GetFavoriteRequest {
    1: i64 id
    2: i64 user_id
    3: i64 video_id
    4: i64 query_type

}

struct GetFavoriteResponse{
    1: list<Favorite1> favorites
}
struct GetFavoriteCountRequest {
    1: i64 video_id
}
struct GetFavoriteCountResponse{
    1: i64 count
}
service FavoriteService {
    FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest req);
    FavoriteListResponse FavoriteList(1: FavoriteListRequest req);
    GetFavoriteResponse GetFavorite(1: GetFavoriteRequest req);
    GetFavoriteCountResponse GetFavoriteCount(1: GetFavoriteCountRequest req);
}