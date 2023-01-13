package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	utils.Response
	VideoList []*Video  `json:"video_list"`
	NextTime  time.Time `json:"next_time"`
}

// Feed
// @Summary 获取所有点赞过的视频
// @Tags feed
// @version 1.0
// @Accept application/x-json-stream
// @Param latest_time query int true "用户id"
// @Param token query string false "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/feed [get]
func Feed(c *gin.Context) {
	// feed允许未登录的用户调用，因此需要在这里判断是否登录了。而不能使用中间件
	token := c.Query("token")

	if token != "" {
		_, err := ParseToken(token)
		if err != nil {
			feedResponse := &FeedResponse{
				Response: utils.Response{
					StatusCode: 1,
					StatusMsg:  "获取到token，解析失败",
				},
				VideoList: nil,
				NextTime:  time.Time{},
			}
			c.JSON(http.StatusOK, feedResponse)
			return
		}
		// 登录的业务处理
		//TODO
		// 业务代码
		feedResponse := &FeedResponse{
			Response: utils.Response{
				StatusCode: 0,
				StatusMsg:  "获取到token，解析失败",
			},
			VideoList: nil, //替换掉这个
			NextTime:  time.Time{},
		}
		c.JSON(http.StatusOK, feedResponse)
		return
	} else {
		//未登录时的业务处理
		//TODO
		// 业务代码
		feedResponse := &FeedResponse{
			Response: utils.Response{
				StatusCode: 0,
				StatusMsg:  "获取到token，解析失败",
			},
			VideoList: nil, //替换掉这个
			NextTime:  time.Time{},
		}
		c.JSON(http.StatusOK, feedResponse)
		return
	}

}
