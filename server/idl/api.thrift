namespace go api


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

struct LoginRequest{
    1: required string username (api.query="username", api.vd="len($)>0 && len($)<=32");
    2: required string password (api.query="password", api.vd="len($)>0 && len($)<=32");
}

struct LoginResponse {
  1: Code status_code;
  2: string status_msg;
  3: i64 user_id;
  4: string token;
}
struct RegisterRequest{
    1: required string username (api.query="username", api.vd="len($)>0 && len($)<=32");
    2: required string password (api.query="password", api.vd="len($)>0 && len($)<=32");
}

struct RegisterResponse {
  1: Code status_code;
  2: string status_msg;
  3: i64 user_id;
  4: string token;
}
struct UserRequest{
    1: required string user_id (api.query="user_id", api.vd="len($)>0 && len($)<=32");
    2: required string token (api.query="token");
}

struct UserResponse {
  1: Code status_code;
  2: string status_msg;
  3: User user;
}
struct Comment {
  1: i64 id;
  2: User user;
  3: string content;
  4: string create_date;
}
struct Message{
  1: i64 id;
  2: string content;
  3: string create_time;
}






service UserService {
    LoginResponse Login(1: LoginRequest req) (api.post="/douyin/user/login/");
    RegisterResponse Register(1: RegisterRequest req) (api.post="/douyin/user/register/");
    UserResponse GetUser(1: UserRequest req) (api.get="/douyin/user/");

}




struct RelationActionRequest{
    1: required string token (api.query="token");
    2: required string to_user_id (api.query="to_user_id");
    3: required i32 action_type (api.query="action_type");
}

struct RelationActionResponse{
    1: Code status_code;
    2: string status_msg;
}

struct RelationFollowListRequest{
    1: required string token (api.post="token");
    2: required string user_id (api.post="user_id");
}

struct RelationFollowListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}

struct RelationFollowerListRequest{
    1: required string token (api.query="token");
    2: required string user_id (api.query="user_id");
}

struct RelationFollowerListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}

struct RelationFriendListRequest{
    1: required string token (api.query="token");
    2: required string user_id (api.query="user_id");
}

struct RelationFriendListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}







service RelationService {
    RelationActionResponse Action(1: RelationActionRequest req) (api.post="/douyin/relation/action/");
    RelationFollowListResponse FollowList(1: RelationFollowListRequest req) (api.get="/douyin/relation/follow/list/");
    RelationFollowerListResponse FollowerList(1: RelationFollowerListRequest req) (api.get="/douyin/relation/follower/list/");
    RelationFriendListResponse FriendList(1: RelationFriendListRequest req) (api.get="/douyin/relation/friend/list/");
}

struct FeedRequest{
    1: string token (api.query="token");
    2: string latest_time (api.query="latest_time");
}

struct FeedResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
    4: i64 next_time;
}


service FeedService {
    FeedResponse Feed(1: FeedRequest req) (api.get="/douyin/feed/");
}
struct PublishActionRequest{
    1: required binary data (api.form="data");
    2: required string token (api.form="token");
    3: required string title (api.form="title");
}

struct PublishActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct PublishListRequest{
    1: required string token (api.query="token");
    2: required string user_id (api.query="user_id");
}

struct PublishListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
}
service PublishService {
    PublishActionResponse PublishAction(1: PublishActionRequest req) (api.post="/douyin/publish/action/");
    PublishListResponse PublishList(1: PublishListRequest req) (api.get="/douyin/publish/list/");
}


struct FavoriteActionRequest{
    1: required string video_id (api.query="video_id");
    2: required string token (api.query="token");
    3: required string action_type (api.query="action_type");
}

struct FavoriteActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct FavoriteListRequest{
    1: required string token (api.query="token");
    2: required string user_id (api.query="user_id");
}

struct FavoriteListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Video> video_list;
}



service FavoriteService {
    FavoriteActionResponse FavoriteAction(1: FavoriteActionRequest req) (api.post="/douyin/favorite/action/");
    FavoriteListResponse FavoriteList(1: FavoriteListRequest req) (api.get="/douyin/favorite/list/");
}



struct CommentActionRequest{
    1: required string video_id (api.query="video_id");
    2: required string token (api.query="token");
    3: required string action_type (api.query="action_type");
    4: string comment_text (api.query="comment_text");
    5: string comment_id (api.query="comment_id");
}

struct CommentActionResponse{
    1: Code status_code;
    2: string status_msg;
    3: Comment comment;
}
struct CommentListRequest{
    1: required string token (api.query="token");
    2: required string video_id (api.query="video_id");
}

struct CommentListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Comment> comment_list;
}


service CommentService {
    CommentActionResponse CommentAction(1: CommentActionRequest req) (api.post="/douyin/comment/action/");
    CommentListResponse CommentList(1: CommentListRequest req) (api.get="/douyin/comment/list/");
}

struct MessageActionRequest{
    1: required string to_user_id (api.query="to_user_id");
    2: required string token (api.query="token");
    3: required string action_type (api.query="action_type");
    4: required string content (api.query="content");
}

struct MessageActionResponse{
    1: Code status_code;
    2: string status_msg;
}
struct MessageListRequest{
    1: required string token (api.query="token");
    2: required string to_user_id (api.query="to_user_id");
}

struct MessageListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<Message> message_list;
}

service MessageService {
    MessageActionResponse MessageAction(1: MessageActionRequest req) (api.post="/douyin/message/action/");
    MessageListResponse MessageList(1: MessageListRequest req) (api.get="/douyin/message/list/");
}