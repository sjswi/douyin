package controllers

import "github.com/gin-gonic/gin"

func FavoriteAction(c *gin.Context) {

}

// FavoriteList
// @Summary 获取所有点赞过的视频
// @Tags 评论
// @version 1.0
// @Accept application/x-json-stream
// @Param video query true
// @Success 200 object PostCommentListForm 成功后返回值
// @Failure 409 object CommentListResponse 失败后返回值
// @Router /douyin/favorite/list [get]
func FavoriteList(c *gin.Context) {

}
