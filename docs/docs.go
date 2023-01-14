// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "余晓兵",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/douyin/comment/action": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "评论"
                ],
                "summary": "评论操作，删除或增加",
                "parameters": [
                    {
                        "description": "文章",
                        "name": "commentAction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PostCommentActionForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentActionResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentActionResponse"
                        }
                    }
                }
            }
        },
        "/douyin/comment/list": {
            "get": {
                "tags": [
                    "评论"
                ],
                "summary": "获取所有评论",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "视频id",
                        "name": "video_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentListResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentListResponse"
                        }
                    }
                }
            }
        },
        "/douyin/favorite/action": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "点赞"
                ],
                "summary": "点赞操作，点赞或取消点赞",
                "parameters": [
                    {
                        "description": "文章",
                        "name": "favoriteAction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PostFavoriteActionForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentActionResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.CommentActionResponse"
                        }
                    }
                }
            }
        },
        "/douyin/favorite/list": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "点赞"
                ],
                "summary": "获取所有点赞过的视频",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    }
                }
            }
        },
        "/douyin/feed": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "feed"
                ],
                "summary": "获取所有点赞过的视频",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "latest_time",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    }
                }
            }
        },
        "/douyin/message/action": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "消息"
                ],
                "summary": "发送消息",
                "parameters": [
                    {
                        "description": "消息",
                        "name": "messageAction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PostMessageActionForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/douyin/message/chat": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "消息"
                ],
                "summary": "获取聊天记录",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    }
                }
            }
        },
        "/douyin/publish/action": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "发布"
                ],
                "summary": "获取聊天记录",
                "parameters": [
                    {
                        "description": "视频信息",
                        "name": "publishAction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.PostPublishActionForm"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/douyin/publish/list": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "发布"
                ],
                "summary": "获取聊天记录",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.FavoriteListResponse"
                        }
                    }
                }
            }
        },
        "/douyin/relation/action": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "关系"
                ],
                "summary": "关注和取消关注操作",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "to_user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "操作类型",
                        "name": "action_type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/douyin/relation/follow/list": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "关系"
                ],
                "summary": "获取关注列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    }
                }
            }
        },
        "/douyin/relation/follower/list": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "关系"
                ],
                "summary": "获取关注者列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    }
                }
            }
        },
        "/douyin/relation/friend/list": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "关系"
                ],
                "summary": "获取聊天记录",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.RelationList"
                        }
                    }
                }
            }
        },
        "/douyin/user": {
            "get": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户信息",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "用户id",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "token",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.UserInfoResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/controllers.UserInfoResponse"
                        }
                    }
                }
            }
        },
        "/douyin/user/login": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户登录",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/douyin/user/register": {
            "post": {
                "consumes": [
                    "application/x-json-stream"
                ],
                "tags": [
                    "用户"
                ],
                "summary": "用户注册",
                "parameters": [
                    {
                        "type": "string",
                        "description": "用户名",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "密码",
                        "name": "password",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Comment": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "create_date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "user": {
                    "$ref": "#/definitions/controllers.User"
                }
            }
        },
        "controllers.CommentActionResponse": {
            "type": "object",
            "properties": {
                "comment": {
                    "$ref": "#/definitions/controllers.Comment"
                },
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                }
            }
        },
        "controllers.CommentListResponse": {
            "type": "object",
            "properties": {
                "comment_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.Comment"
                    }
                },
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                }
            }
        },
        "controllers.FavoriteListResponse": {
            "type": "object",
            "properties": {
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                },
                "video_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.Video"
                    }
                }
            }
        },
        "controllers.PostCommentActionForm": {
            "type": "object",
            "properties": {
                "action_type": {
                    "type": "integer"
                },
                "comment_id": {
                    "type": "string"
                },
                "comment_text": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "video_id": {
                    "type": "integer"
                }
            }
        },
        "controllers.PostFavoriteActionForm": {
            "type": "object",
            "properties": {
                "action_type": {
                    "description": "1: 点赞，2:取消点赞",
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                },
                "video_id": {
                    "type": "integer"
                }
            }
        },
        "controllers.PostMessageActionForm": {
            "type": "object",
            "properties": {
                "action_type": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "to_user_id": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "controllers.PostPublishActionForm": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "data": {
                    "type": "integer"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "controllers.RelationList": {
            "type": "object",
            "properties": {
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                },
                "user_list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.User"
                    }
                }
            }
        },
        "controllers.User": {
            "type": "object",
            "properties": {
                "follow_count": {
                    "type": "integer"
                },
                "follower_count": {
                    "type": "integer"
                },
                "is_follow": {
                    "type": "boolean"
                },
                "user_id": {
                    "type": "integer"
                },
                "user_name": {
                    "type": "string"
                }
            }
        },
        "controllers.UserInfoResponse": {
            "type": "object",
            "properties": {
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/controllers.User"
                }
            }
        },
        "controllers.Video": {
            "type": "object",
            "properties": {
                "comment_count": {
                    "type": "integer"
                },
                "cover_url": {
                    "type": "string"
                },
                "favorite_count": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "is_favorite": {
                    "type": "boolean"
                },
                "play_url": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/controllers.User"
                }
            }
        },
        "utils.Response": {
            "type": "object",
            "properties": {
                "status_code": {
                    "type": "integer"
                },
                "status_msg": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:8080",
	BasePath:         "/douyin",
	Schemes:          []string{},
	Title:            "douyin",
	Description:      "青训营抖音大项目",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
