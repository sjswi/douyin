package controllers

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
// @Router /douyin/favorite/action/ [post]
func FavoriteAction(c *gin.Context) {
	auth := auth2.Auth{}.GetAuth(c)
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
	err = service.FavoriteActionPost(uint(videoID), action, auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 4、装配返回值

	c.JSON(http.StatusOK, successResponse)
}

type FavoriteListResponse struct {
	utils.Response
	VideoList []vo.Video `json:"video_list"`
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
// @Router /douyin/favorite/list/ [get]
func FavoriteList(c *gin.Context) {
	auth := auth2.Auth{}.GetAuth(c)
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
	if err != nil || userID < 0 {
		failureResponse.StatusMsg = "user_id必须大于等于0"
		c.JSON(409, failureResponse)
		return
	}
	if userID == 0 {
		userID = int(auth.UserID)
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
	returnFavorite := make([]vo.Video, len(favoriteList))
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
		returnFavorite[i] = vo.Video{
			Author: &vo.User{
				ID:            author.ID,
				Name:          author.Name,
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
