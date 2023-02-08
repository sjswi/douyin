namespace go favorite

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

struct FavoriteActionRequest{
    1: required string video_id;
    2: required i64 auth_id;
    3: required string action_type;
}

struct FavoriteActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct FavoriteListRequest{
    1: required i64 auth_id;
    2: required string user_id;
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