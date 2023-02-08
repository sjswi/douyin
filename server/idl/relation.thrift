namespace go relation

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

struct RelationActionRequest{
    1: required i64 auth_id;
    2: required string to_user_id;
    3: required i32 action_type;
}

struct RelationActionResponse{
    1: Code status_code;
    2: string status_msg;
}

struct RelationFollowListRequest{
    1: required i64 auth_id;
    2: required string user_id;
}

struct RelationFollowListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}

struct RelationFollowerListRequest{
    1: required i64 auth_id;
    2: required string user_id;
}

struct RelationFollowerListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}

struct RelationFriendListRequest{
    1: required i64 auth_id;
    2: required string user_id;
}

struct RelationFriendListResponse{
    1: Code status_code;
    2: string status_msg;
    3: list<User> user_list;
}

struct Relation1{
    1: i64 id
    2: i64 user_id
    3: i64 target_id
    4: i32 type
}

struct GetRelationRequest {
    1: i64 id
    2: i64 user_id
    3: i64 target_id
    4: i64 relation_type
    5: i64 query_type

}

struct GetRelationResponse{
    1: list<Relation1> relations
}
struct GetCountRequest {
    1: i64 user_id
}
struct GetCountResponse{
    1: i64 FollowCount
    2: i64 FollowerCount
    3: i64 FriendCount
}

service RelationService {
    RelationActionResponse Action(1: RelationActionRequest req);
    RelationFollowListResponse FollowList(1: RelationFollowListRequest req);
    RelationFollowerListResponse FollowerList(1: RelationFollowerListRequest req);
    RelationFriendListResponse FriendList(1: RelationFriendListRequest req);
    GetRelationResponse GetRelation(1: GetRelationRequest req);
    GetCountResponse GetCount(1: GetCountRequest req);
}