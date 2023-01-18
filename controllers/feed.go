package controllers

import (
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	utils.Response
	VideoList []vo.Video `json:"video_list"`
	NextTime  int64      `json:"next_time"`
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
// @Router /douyin/feed/ [get]
func Feed(c *gin.Context) {
	// feed允许未登录的用户调用，因此需要在这里判断是否登录了。而不能使用中间件
	successResponse := FeedResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: nil,
		NextTime:  0,
	}
	failureResponse := FeedResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		VideoList: nil,
		NextTime:  0,
	}
	// 1、解析参数
	token := c.Query("token")
	latestTime := c.Query("latest_time")
	// 验证参数
	var latest time.Time
	if latestTime == "" || latestTime == "0" {
		latest = time.Now()
	} else {
		//if len(latestTime) > 10 {
		//	latestTime = latestTime[:10]
		//}
		lateUnix, err := strconv.Atoi(latestTime)
		if err != nil {
			failureResponse.StatusMsg = "时间戳必须为int，解析失败"
			c.JSON(409, failureResponse)
			return
		}
		latest = time.UnixMilli(int64(lateUnix))
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
		videos, nextTime, err := service.FeedGet(latest, auth)
		if err != nil {
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		successResponse.VideoList = videos
		successResponse.NextTime = nextTime
	} else {
		videos, nextTime, err := service.FeedGet(latest, nil)
		if err != nil {
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		successResponse.VideoList = videos
		successResponse.NextTime = nextTime
	}
	// 4、装配返回值
	c.JSON(http.StatusOK, successResponse)
	return

}
