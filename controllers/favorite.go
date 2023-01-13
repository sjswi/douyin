package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostFavoriteActionForm struct {
	Token      string `json:"token"`
	VideoID    uint   `json:"video_id"`
	ActionType int    `json:"action_type"` //1: 点赞，2:取消点赞
}
type FavoriteActionResponse struct {
	utils.Response
}

// FavoriteAction
// @Summary 点赞操作，点赞或取消点赞
// @Tags 点赞
// @version 1.0
// @Accept application/x-json-stream
// @Param favoriteAction body PostFavoriteActionForm true "文章"
// @Success 200 object CommentActionResponse 成功后返回值
// @Failure 409 object CommentActionResponse 失败后返回值
// @Router /douyin/favorite/action [post]
func FavoriteAction(c *gin.Context) {
	var favoriteActionForm *PostFavoriteActionForm
	if err := c.ShouldBind(&favoriteActionForm); err != nil {
		response := FavoriteActionResponse{
			Response: utils.Response{
				StatusCode: 1,
				StatusMsg:  "data数据有误",
			},
		}
		c.JSON(409, response)
		return
	}
	//具体业务
	//TODO
	// 业务代码
	//成功
	response := CommentActionResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Comment: nil,
	}
	c.JSON(http.StatusOK, response)
}

type Video struct {
	User          *User
	ID            uint   `json:"id"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
}
type FavoriteListResponse struct {
	utils.Response
	VideoList []*Video `json:"video_list"`
}

// FavoriteList
// @Summary 获取所有点赞过的视频
// @Tags 点赞
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/favorite/list [get]
func FavoriteList(c *gin.Context) {
	//TODO
	// 业务代码
	//成功返回
	resp := FavoriteListResponse{
		Response:  utils.Response{},
		VideoList: nil,
	}

	c.JSON(http.StatusOK, resp)
}
