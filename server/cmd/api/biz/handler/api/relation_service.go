// Code generated by hertz generator.

package api

import (
	"context"
	"douyin_rpc/server/cmd/api/global"
	"douyin_rpc/server/cmd/api/kitex_gen/relation"
	"strconv"

	api "douyin_rpc/server/cmd/api/biz/model/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Action
// @Summary 关注和取消关注操作
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Param action_type query int true "操作类型"
// @Success 200 object relation.RelationActionResponse 成功后返回值
// @Failure 409 object relation.RelationActionResponse 失败后返回值
// @Router /douyin/relation/action/ [post]
// @router /douyin/relation/action/ [POST]
func Action(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	response := new(api.RelationActionResponse)
	value, exist := c.Get("accountID")
	if !exist {
		return
	}
	toUserID, err := strconv.ParseInt(req.ToUserID, 0, 64)
	if err != nil {
		return
	}
	resp, err := global.RelationClient.Action(ctx, &relation.RelationActionRequest{
		AuthId:     value.(int64),
		ToUserId:   toUserID,
		ActionType: req.ActionType,
	})
	if err != nil {

		response.StatusCode = 4

		response.StatusMsg = err.Error()
		c.JSON(consts.StatusConflict, response)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// FollowList
// @Summary 获取关注列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object api.RelationFollowListResponse 成功后返回值
// @Failure 409 object api.RelationFollowListResponse 失败后返回值
// @Router /douyin/relation/follow/list/ [get]
// @router /douyin/relation/follow/list/ [GET]
func FollowList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RelationFollowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	value, exist := c.Get("accountID")
	if !exist {
		return
	}
	userID, err := strconv.ParseInt(req.UserID, 0, 64)
	if err != nil {
		return
	}
	resp, err := global.RelationClient.FollowList(ctx, &relation.RelationFollowListRequest{
		AuthId: value.(int64),
		UserId: userID,
	})
	if err != nil {

		resp.StatusCode = 4

		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusConflict, resp)
		return
	}
	//resp := new(api.PublishActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// FollowerList
// @Summary 获取关注者列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object api.RelationFollowerListResponse 成功后返回值
// @Failure 409 object api.RelationFollowerListResponse 失败后返回值
// @Router /douyin/relation/follower/list/ [get]
// @router /douyin/relation/follower/list/ [GET]
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RelationFollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	value, exist := c.Get("accountID")
	if !exist {
		return
	}
	userID, err := strconv.ParseInt(req.UserID, 0, 64)
	if err != nil {
		return
	}
	resp, err := global.RelationClient.FollowerList(ctx, &relation.RelationFollowerListRequest{
		AuthId: value.(int64),
		UserId: userID,
	})
	if err != nil {

		resp.StatusCode = 4

		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusConflict, resp)
		return
	}
	//resp := new(api.PublishActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// FriendList
// @Summary 获取好友列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object api.RelationFriendListResponse 成功后返回值
// @Failure 409 object api.RelationFriendListResponse 失败后返回值
// @Router /douyin/relation/friend/list/ [get]
// @router /douyin/relation/friend/list/ [GET]
func FriendList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RelationFriendListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	value, exist := c.Get("accountID")
	if !exist {
		return
	}
	userID, err := strconv.ParseInt(req.UserID, 0, 64)
	if err != nil {
		return
	}
	resp, err := global.RelationClient.FriendList(ctx, &relation.RelationFriendListRequest{
		AuthId: value.(int64),
		UserId: userID,
	})
	if err != nil {

		resp.StatusCode = 4

		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusConflict, resp)
		return
	}
	//resp := new(api.PublishActionResponse)

	c.JSON(consts.StatusOK, resp)
}
