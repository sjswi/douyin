namespace go user

struct User {
    1: i64 id;
    2: string name;
    3: i64 follow_count;
    4: i64 follower_count;
    5: bool is_follow;
}
struct User1 {
    1: i64 id
    2: string name
    3: string password
    4: i64 created_at
    5: i64 updated_at
    6: string salt
}
enum Code {
    Success = 0;
    ParamInvalid = 1;
    DBError = 2;
    ServerError = 3;
    ErrorRequest = 4;
}
struct LoginRequest{
    1: required string username;
    2: required string password;
}

struct LoginResponse {
  1: Code status_code;
  2: string status_msg;
  3: i64 user_id;
}
struct RegisterRequest{
    1: required string username;
    2: required string password;
}

struct RegisterResponse {
  1: Code status_code;
  2: string status_msg;
  3: i64 user_id;
}
struct UserRequest{
    1: required i64 user_id;
    2: required i64 auth_id;
}

struct UserResponse {
  1: Code status_code;
  2: string status_msg;
  3: User user;
}
/*
* query_type=1 根据user_id查询用户
* query_type=2 根据username查询用户
* */
struct GetUserRequest {
    1: i64 user_id
    2: string username
    3: i32 query_type
}
struct GetUserResponse {
    1: Code status_code
    2: User1 user
}

service UserService{
    LoginResponse Login(1: LoginRequest req);
    RegisterResponse Register(1: RegisterRequest req);
    GetUserResponse GetUser(1: GetUserRequest req);
    UserResponse User(1: UserRequest req);
}