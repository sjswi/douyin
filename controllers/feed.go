package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/config"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	utils.Response
	VideoList []Video `json:"video_list"`
	NextTime  string  `json:"next_time"`
}

// Feed
// @Summary 视频推流
// @Tags feed
// @version 1.0
// @Accept application/x-json-stream
// @Param latest_time query int false "用户id"
// @Param token query string false "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/feed [get]
func Feed(c *gin.Context) {
	// feed允许未登录的用户调用，因此需要在这里判断是否登录了。而不能使用中间件
	successResponse := FeedResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: nil,
		NextTime:  "",
	}
	failureResponse := FeedResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		VideoList: nil,
		NextTime:  "",
	}
	// 1、解析参数
	token := c.Query("token")
	latestTime := c.Query("latest_time")
	// 验证参数
	var latest time.Time
	if latestTime == "" {
		latest = time.Now()
	} else {
		lateUnix, err := strconv.Atoi(latestTime)
		if err != nil {
			failureResponse.StatusMsg = "时间戳必须为int，解析失败"
			c.JSON(409, failureResponse)
			return
		}
		latest = time.Unix(int64(lateUnix), 0)
		if err != nil {
			failureResponse.StatusMsg = "解析时间错误错误"
			c.JSON(409, failureResponse)
			return
		}
	}

	// 3、具体业务
	// 判断是否提供了token
	if token != "" {

		auth, err := ParseToken(token)
		//auth :=
		if err != nil {
			failureResponse.StatusMsg = "解析token错误，token有误"
			c.JSON(409, failureResponse)
			return
		}
		//TODO
		// 推荐算法

		// 目前直接从video中获取给定数量的最新的视频
		var videoList []models.Video
		if err := driver.Db.Debug().Model(models.Video{}).Where("created_at<=?", latest).Limit(config.Number).Find(&videoList).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
		returnVideo := make([]Video, len(videoList))
		for i := 0; i < len(videoList); i++ {
			var author models.User
			var relation models.Relation
			var favorite models.Favorite
			// 查询每个视频的作者
			if err := driver.Db.Debug().Model(author).Where("id=?", videoList[i].AuthorID).Find(&author).Error; err != nil {
				failureResponse.StatusMsg = "查询数据库错误" + err.Error()
				c.JSON(409, failureResponse)
				return
			}
			// 查询auth用户是否关注了作者
			if err := driver.Db.Debug().Model(relation).Where("exist=1 and (type=1 or type=2) and user_id=? and target_id=?", auth.UserID, author.ID).Find(&relation).Error; err != nil {
				failureResponse.StatusMsg = "查询数据库错误" + err.Error()
				c.JSON(409, failureResponse)
				return
			}
			// 查询用户是否点赞该视频
			if err := driver.Db.Debug().Model(favorite).Where("exist=1 and user_id=? and video_id=?", auth.UserID, videoList[i].ID).Find(&favorite).Error; err != nil {
				failureResponse.StatusMsg = "查询数据库错误" + err.Error()
				c.JSON(409, failureResponse)
				return
			}
			returnVideo[i] = Video{
				Author: &User{
					UserID:        author.ID,
					UserName:      author.Name,
					FollowCount:   author.FollowCount,
					FollowerCount: author.FollowerCount,
					IsFollow:      false,
				},
				ID:            videoList[i].ID,
				FavoriteCount: videoList[i].FavoriteCount,
				CommentCount:  videoList[i].CommentCount,
				IsFavorite:    false,
				Title:         videoList[i].Title,
				PlayURL:       videoList[i].PlayURL,
				CoverURL:      videoList[i].CoverURL,
			}
			if relation.ID != 0 && relation.Exist {
				returnVideo[i].Author.IsFollow = true
			}
			if favorite.ID != 0 && favorite.Exist {
				returnVideo[i].IsFavorite = true
			}
		}
		successResponse.NextTime = videoList[len(videoList)-1].CreatedAt.Format("2006-01-02 15:04:05")
		successResponse.VideoList = returnVideo
	} else {
		var videoList []models.Video
		if err := driver.Db.Debug().Model(models.Video{}).Where("created_at<=?", latest).Limit(config.Number).Find(&videoList).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库错误" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
		returnVideo := make([]Video, len(videoList))
		for i := 0; i < len(videoList); i++ {
			var author models.User
			// 查询每个视频的作者
			if err := driver.Db.Debug().Model(author).Where("id=?", videoList[i].AuthorID).Find(&author).Error; err != nil {
				failureResponse.StatusMsg = "查询数据库错误" + err.Error()
				c.JSON(409, failureResponse)
				return
			}
			returnVideo[i] = Video{
				Author: &User{
					UserID:        author.ID,
					UserName:      author.Name,
					FollowCount:   author.FollowCount,
					FollowerCount: author.FollowerCount,
					IsFollow:      false,
				},
				ID:            videoList[i].ID,
				FavoriteCount: videoList[i].FavoriteCount,
				CommentCount:  videoList[i].CommentCount,
				IsFavorite:    false,
				Title:         videoList[i].Title,
				PlayURL:       videoList[i].PlayURL,
				CoverURL:      videoList[i].CoverURL,
			}

		}
		successResponse.NextTime = videoList[len(videoList)-1].CreatedAt.Format("2006-01-02 15:04:05")
		successResponse.VideoList = returnVideo
	}
	// 4、装配返回值
	c.JSON(http.StatusOK, successResponse)
	return

}
