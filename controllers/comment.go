package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PostCommentActionForm struct {
	Token       string `json:"token"`
	VideoID     uint   `json:"video_id"`
	ActionType  int    `json:"action_type"`
	CommentText string `json:"comment_text"`
	CommentID   string `json:"comment_id"`
}
type User struct {
	UserID        uint   `json:"user_id"`
	UserName      string `json:"user_name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
type Comment struct {
	ID         uint `json:"id"`
	User       *User
	Content    string    `json:"content"`
	CreateDate time.Time `json:"create_date"`
}
type CommentActionResponse struct {
	utils.Response
	Comment *Comment
}

// CommentAction
// @Summary 评论操作，删除或增加
// @Tags 评论
// @version 1.0
// @Accept application/x-json-stream
// @Param commentAction body PostCommentActionForm true "文章"
// @Success 200 object CommentActionResponse 成功后返回值
// @Failure 409 object CommentActionResponse 失败后返回值
// @Router /douyin/comment/action [post]
func CommentAction(c *gin.Context) {
	var commentActionForm *PostCommentActionForm
	if err := c.ShouldBind(&commentActionForm); err != nil {
		response := CommentActionResponse{
			Response: utils.Response{
				StatusCode: 1,
				StatusMsg:  "data数据有误",
			},
			Comment: nil,
		}
		c.JSON(409, response)
		return
	}
	//具体业务

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

type CommentListResponse struct {
	utils.Response
	CommentList []*Comment `json:"comment_list"`
}

// CommentList
// @Summary 获取所有评论
// @Tags 评论
// @version 1.0
// @Accept application/x-json-stream
// @Param video_id query true
// @Success 200 object CommentListResponse 成功后返回值
// @Failure 409 object CommentListResponse 失败后返回值
// @Router /douyin/comment/list [get]
func CommentList(c *gin.Context) {
	//获取参数
	//videoID := c.Query("video_id")

	//
	response := CommentListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CommentList: nil,
	}
	c.JSON(http.StatusOK, response)
}
