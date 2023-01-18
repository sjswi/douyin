package controllers

import (
	auth2 "douyin/auth"
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentActionResponse struct {
	utils.Response
	Comment *vo.Comment `json:"comment"`
}

// CommentAction
// @Summary 评论操作，删除或增加
// @Tags 评论
// @version 1.0
// @Accept application/x-json-stream
// @Param video_id query int true "视频id"
// @Param token query string true "token"
// @Param action_type query int true "操作类型"
// @Param comment_text query string false "评论内容"
// @Param comment_id query int false "评论id"
// @Success 200 object CommentActionResponse 成功后返回值
// @Failure 409 object CommentActionResponse 失败后返回值
// @Router /douyin/comment/action/ [post]
func CommentAction(c *gin.Context) {
	auth := auth2.Auth{}.GetAuth(c)
	successResponse := CommentActionResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Comment: nil,
	}
	failureResponse := CommentActionResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		Comment: nil,
	}
	// 1、解析参数
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentId := c.Query("comment_id")
	// 2、验证参数
	// videoID必须为正整数，commentID一样，actionType必须为1或2
	videoID, err := strconv.Atoi(videoId)
	if err != nil {
		failureResponse.StatusMsg = "video_id必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	action, err := strconv.Atoi(actionType)
	if err != nil || (action != 1 && action != 2) {
		failureResponse.StatusMsg = "action_type必须为1或2"
		c.JSON(409, failureResponse)
		return
	}
	commentID := -1
	// 操作为1，则不能存在commentID
	if action == 1 {
		if commentText == "" {
			failureResponse.StatusMsg = "评论不能为空"
			c.JSON(409, failureResponse)
			return
		}
		if commentId != "" {
			failureResponse.StatusMsg = "commentId不允许存在"
			c.JSON(409, failureResponse)
			return
		}
	} else {
		if commentText != "" {
			failureResponse.StatusMsg = "删除操作不允许出现评论内容"
			c.JSON(409, failureResponse)
			return
		}

		id, err := strconv.Atoi(commentId)
		if err != nil || id <= 0 {
			failureResponse.StatusMsg = "评论id必须大于等于0"
			c.JSON(409, failureResponse)
			return
		}
		commentID = id

	}

	// 4、装配返回值
	successResponse.Comment, err = service.CommentActionPost(uint(videoID), uint(commentID), action, commentText, auth)
	if err != nil {
		failureResponse.StatusMsg = err.Error()
		c.JSON(409, failureResponse)
		return
	}
	c.JSON(http.StatusOK, successResponse)
}

type CommentListResponse struct {
	utils.Response
	CommentList []vo.Comment `json:"comment_list"`
}

// CommentList
// @Summary 获取所有评论
// @Tags 评论
// @version 1.0
// @Param video_id query int true "视频id"
// @Param token query string true "token"
// @Success 200 object CommentListResponse 成功后返回值
// @Failure 409 object CommentListResponse 失败后返回值
// @Router /douyin/comment/list/ [get]
func CommentList(c *gin.Context) {
	//获取参数
	//videoID := c.Query("video_id")

	auth := auth2.Auth{}.GetAuth(c)

	successResponse := CommentListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CommentList: nil,
	}
	failureResponse := CommentListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CommentList: nil,
	}

	// 1、解析参数
	videoId := c.Query("video_id")
	// 2、验证参数
	videoID, err := strconv.Atoi(videoId)
	if err != nil || videoID < 0 {
		failureResponse.StatusMsg = "解析video_id错误"
		c.JSON(409, failureResponse)
		return
	}
	// 3、具体业务
	returnComment, err := service.CommentListGet(videoID, auth)
	if err != nil {
		failureResponse.StatusMsg = "服务失败"
		c.JSON(409, failureResponse)
		return
	}
	successResponse.CommentList = returnComment
	c.JSON(http.StatusOK, successResponse)
}
