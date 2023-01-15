package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
type Comment struct {
	ID         uint   `json:"id"`
	User       *User  `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
type CommentActionResponse struct {
	utils.Response
	Comment *Comment `json:"comment"`
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
// @Router /douyin/comment/action [post]
func CommentAction(c *gin.Context) {
	auth := Auth{}.GetAuth(c)
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
	var video models.Video
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(video).Set("gorm:query_option", "FOR UPDATE").Where("id = ?", videoID).Find(&video).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误"
		c.JSON(409, failureResponse)
		return
	}
	// 3、具体业务
	var returnComment Comment
	if action == 1 {
		comment := models.Comment{
			Model:      gorm.Model{},
			Content:    commentText,
			CreateTime: time.Now().UTC(),
			VideoID:    uint(videoID),
			UserID:     auth.UserID,
		}
		if err := tx.Model(models.Comment{}).Create(&comment).Error; err != nil {
			failureResponse.StatusMsg = "创建失败" + err.Error()
			c.JSON(409, failureResponse)
			tx.Rollback()
			return
		}
		returnComment = Comment{
			ID: comment.ID,
			User: &User{
				ID:            auth.UserID,
				Name:          auth.UserName,
				FollowCount:   auth.FollowCount,
				FollowerCount: auth.FollowerCount,
				IsFollow:      false,
			},
			Content:    comment.Content,
			CreateDate: utils.GetMonthAndDay(comment.CreateTime),
		}
		video.CommentCount += 1
	} else {
		var comment models.Comment
		if err := tx.Model(comment).Where("id = ?", commentID).Delete(&comment).Error; err != nil {
			failureResponse.StatusMsg = "删除失败" + err.Error()
			c.JSON(409, failureResponse)
			tx.Rollback()
			return
		}

		returnComment = Comment{
			ID: uint(commentID),
			User: &User{
				ID:            auth.UserID,
				Name:          auth.UserName,
				FollowCount:   auth.FollowCount,
				FollowerCount: auth.FollowerCount,
				IsFollow:      false,
			},
			Content:    commentText,
			CreateDate: "",
		}
		video.CommentCount -= 1
	}
	// 更新video信息
	if err := tx.Model(video).Save(&video).Error; err != nil {
		failureResponse.StatusMsg = "存储video信息错误"
		c.JSON(409, failureResponse)
		return
	}
	if err := tx.Commit().Error; err != nil {
		failureResponse.StatusMsg = "事务提交失败" + err.Error()
		c.JSON(409, failureResponse)
		tx.Rollback()
		return
	}

	// 4、装配返回值
	successResponse.Comment = &returnComment
	c.JSON(http.StatusOK, successResponse)
}

type CommentListResponse struct {
	utils.Response
	CommentList []Comment `json:"comment_list"`
}

// CommentList
// @Summary 获取所有评论
// @Tags 评论
// @version 1.0
// @Param video_id query int true "视频id"
// @Param token query string true "token"
// @Success 200 object CommentListResponse 成功后返回值
// @Failure 409 object CommentListResponse 失败后返回值
// @Router /douyin/comment/list [get]
func CommentList(c *gin.Context) {
	//获取参数
	//videoID := c.Query("video_id")

	auth := Auth{}.GetAuth(c)

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
	if err != nil {
		failureResponse.StatusMsg = "解析video_id错误"
		c.JSON(409, failureResponse)
		return
	}
	// 3、具体业务
	var commentList []models.Comment
	if err := driver.Db.Debug().Model(models.Comment{}).Where("video_id=?", videoID).Order("create_time DESC").Find(&commentList).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	//TODO
	// 联表查询优化

	// 4、装配返回值
	returnComment := make([]Comment, len(commentList))
	for i := 0; i < len(commentList); i++ {
		var relation models.Relation
		var user models.User
		returnComment[i].ID = commentList[i].ID
		returnComment[i].Content = commentList[i].Content
		returnComment[i].CreateDate = utils.GetMonthAndDay(commentList[i].CreateTime)
		if err := driver.Db.Debug().Model(user).Where("id=?", commentList[i].UserID).Find(&user).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
		returnComment[i].User = &User{
			ID:            user.ID,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      false,
		}
		if err := driver.Db.Debug().Model(relation).Where("user_id=?", auth.UserID).Where("target_id=?", user.ID).Where("exist=1").Where("type=1 or type=2").Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
		if relation.ID != 0 && relation.Exist {
			returnComment[i].User.IsFollow = true
		}
	}
	successResponse.CommentList = returnComment
	c.JSON(http.StatusOK, successResponse)
}
