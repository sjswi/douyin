package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostFavoriteActionForm struct {
	Token      string `json:"token"`
	VideoID    uint   `json:"video_id"`
	ActionType int    `json:"action_type"` //1: 点赞，2:取消点赞
}

// FavoriteAction
// @Summary 点赞操作，点赞或取消点赞
// @Tags 点赞
// @version 1.0
// @Accept application/x-json-stream
// @Param token query string true "token"
// @Param video_id query int true "视频id"
// @Param action_type query int true "事件类型，1点赞，2取消点赞"
// @Success 200 object CommentActionResponse 成功后返回值
// @Failure 409 object CommentActionResponse 失败后返回值
// @Router /douyin/favorite/action [post]
func FavoriteAction(c *gin.Context) {
	auth := Auth{}.GetAuth(c)
	successResponse := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	failureResponse := utils.Response{
		StatusCode: 1,
		StatusMsg:  "",
	}
	// 1、解析参数
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	// 2、验证参数
	action, err := strconv.Atoi(actionType)
	if err != nil || (action != 1 && action != 2) {
		failureResponse.StatusMsg = "action必须为1或2"
		c.JSON(409, failureResponse)
		return
	}
	videoID, err := strconv.Atoi(videoId)
	if err != nil || videoID <= 0 {
		failureResponse.StatusMsg = "视频id必须大于0"
		c.JSON(409, failureResponse)
		return
	}
	// 3、具体业务
	// 3.1 业务分析：
	// 		1）查询视频是否存在，不存在直接返回
	// 		2）查询点赞表查看是否点赞过，点赞过直接返回
	//
	var video models.Video
	// 事务开始
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 3.2 查询数据库获取视频信息
	if err := tx.Model(video).Set("gorm:query_option", "FOR UPDATE").Where("id = ?", videoID).Find(&video).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误"
		c.JSON(409, failureResponse)
		return
	}
	var favorite models.Favorite
	// 3.3 查询favorite表获取信息，
	if action == 1 {
		// 3.3.1 查看点赞是否存在，如果存在返回
		if err := tx.Model(favorite).Where("user_id = ? and video_id = ?", auth.UserID, videoID).Find(&favorite).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		if favorite.ID == 0 {
			// 3.3.2 不存在，需要创建
			favorite.UserID = auth.UserID
			favorite.Exist = true
			favorite.VideoID = video.ID

			if err := tx.Model(favorite).Create(&favorite).Error; err != nil {
				failureResponse.StatusMsg = "查询数据库错误"
				c.JSON(409, failureResponse)
				return
			}
		} else if !favorite.Exist {
			// 3.3.3 之前点赞过后面取消了
			favorite.Exist = true
			if err := tx.Model(favorite).Updates(&favorite).Error; err != nil {
				failureResponse.StatusMsg = "更新favorite信息错误"
				c.JSON(409, failureResponse)
				return
			}
		} else {
			// 3.3.4 赞现在存在
			failureResponse.StatusMsg = "点赞存在，无需再次点赞"
			c.JSON(409, failureResponse)
			return
		}
		video.FavoriteCount += 1
	} else {
		// 3.4.1 查看点赞是否存在，如果存在返回
		if err := tx.Model(favorite).Where("user_id = ? and video_id = ?", auth.UserID, videoID).Find(&favorite).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		if favorite.ID == 0 || !favorite.Exist {
			// 3.3.2 不存在，无法取消
			failureResponse.StatusMsg = "点赞存在，无需再次点赞"
			c.JSON(409, failureResponse)
			return

		} else {
			// 3.3.4 赞现在存在
			favorite.Exist = false
			if err := tx.Model(favorite).Save(&favorite).Error; err != nil {
				failureResponse.StatusMsg = "存储favorite信息错误"
				c.JSON(409, failureResponse)
				return
			}
		}
		video.FavoriteCount -= 1
	}
	// 3.4 更新video信息
	if err := tx.Model(video).Save(&video).Error; err != nil {
		failureResponse.StatusMsg = "存储video信息错误"
		c.JSON(409, failureResponse)
		return
	}
	if err := tx.Commit().Error; err != nil {
		failureResponse.StatusMsg = "事务提交错误"
		c.JSON(409, failureResponse)
		return
	}
	// 4、装配返回值

	c.JSON(http.StatusOK, successResponse)
}

type Video struct {
	Author        *User
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
	VideoList []Video `json:"video_list"`
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
	auth := Auth{}.GetAuth(c)
	successResponse := FavoriteListResponse{
		Response:  utils.Response{},
		VideoList: nil,
	}
	failureResponse := FavoriteListResponse{
		Response:  utils.Response{},
		VideoList: nil,
	}

	// 1、解析参数

	userId := c.Query("user_id")

	// 2、验证参数
	userID, err := strconv.Atoi(userId)
	if err != nil || userID <= 0 {
		failureResponse.StatusMsg = "user_id必须大于等于0"
		c.JSON(409, failureResponse)
		return
	}
	//TODO
	// mysql联表查询优化，redis缓存优化

	// 3、具体业务
	// 3.1、首先查询用户所有的点赞信息
	var favoriteList []models.Favorite
	if err := driver.Db.Debug().Model(models.Favorite{}).Where("exist=1").Where("user_id=?", userID).Find(&favoriteList).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误"
		c.JSON(409, failureResponse)
		return
	}
	returnFavorite := make([]Video, len(favoriteList))
	for i := 0; i < len(favoriteList); i++ {
		var video models.Video
		var author models.User
		var relation models.Relation
		var favorite models.Favorite
		// 3.2、查询视频
		if err := driver.Db.Debug().Model(video).Where("id=?", favoriteList[i].VideoID).Find(&video).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		// 3.3 查询视频author信息
		if err := driver.Db.Debug().Model(author).Where("id=?", video.AuthorID).Find(&author).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		// 3.4 查询视频作者与用户auth的关系
		if err := driver.Db.Debug().Model(relation).Where("exist=1").Where("type=1 or type=2").Where("user_id=? and target_id=?", auth.UserID, author.ID).Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		returnFavorite[i] = Video{
			Author: &User{
				UserID:        author.ID,
				UserName:      author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      false,
			},
			ID:            video.ID,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false,
			Title:         video.Title,
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
		}
		if relation.ID != 0 && relation.Exist {
			returnFavorite[i].Author.IsFollow = true
		}
		// 3.5 查询视频自己是否点过赞
		if err := driver.Db.Debug().Model(favorite).Where("exist=1").Where("user_id=? and video_id=?", auth.UserID, video.ID).Find(&favorite).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误"
			c.JSON(409, failureResponse)
			return
		}
		if favorite.ID != 0 && favorite.Exist {
			returnFavorite[i].IsFavorite = true
		}
	}
	// 4、装配返回值
	successResponse.VideoList = returnFavorite
	c.JSON(http.StatusOK, successResponse)
}
