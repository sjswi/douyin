package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/storage"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// PublishAction
// @Summary 发布视频
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param token formData string true "token"
// @Param data formData file true "视频"
// @Param title formData string true "标题"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/publish/action [post]
func PublishAction(c *gin.Context) {

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
	form, err := c.MultipartForm()
	if err != nil {
		failureResponse.StatusMsg = "解析form表单错误"
		c.JSON(409, failureResponse)
		return
	}
	title := c.PostForm("title")
	// 2、验证参数
	data := form.File["data"]
	if len(data) != 1 {
		failureResponse.StatusMsg = "视频数量只能为1"
		c.JSON(409, failureResponse)
		return
	}
	video, err := data[0].Open()
	if err != nil {
		failureResponse.StatusMsg = "读取视频时错误"
		c.JSON(409, failureResponse)
		return
	}
	//path := config.TempPath + data[0].Filename
	//err = c.SaveUploadedFile(data[0], path)
	//if err != nil {
	//	failureResponse.StatusMsg = "存储视频错误"
	//	c.JSON(409, failureResponse)
	//	return
	//}
	// 3、具体业务
	//snapshot, err := utils.GetSnapshot(path)
	//if err != nil {
	//	return
	//}
	uid := uuid.New().String()
	videoURL := storage.OSS.Put(uid+data[0].Filename, video)
	//coverURL := storage.OSS.Put(uid+".jpeg", snapshot)
	coverURL := videoURL + "?x-oss-process=video/snapshot,t_7000,f_jpg,w_800,h_600,m_fast"
	videoModel := models.Video{
		Model:         gorm.Model{},
		AuthorID:      auth.UserID,
		Title:         title,
		CommentCount:  0,
		FavoriteCount: 0,
		PlayURL:       videoURL,
		CoverURL:      coverURL,
	}
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(videoModel).Create(&videoModel).Error; err != nil {
		failureResponse.StatusMsg = "创建视频错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		failureResponse.StatusMsg = "创建视频错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 4、装配返回值
	c.JSON(http.StatusOK, successResponse)
	return
}

type PublishListResponse struct {
	utils.Response
	VideoList []Video `json:"video_list"`
}

// PublishList
// @Summary 获取发布视频列表
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/publish/list [get]
func PublishList(c *gin.Context) {

	auth := Auth{}.GetAuth(c)

	successResponse := PublishListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: nil,
	}
	failureResponse := PublishListResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		VideoList: nil,
	}
	// 1、解析参数
	userId := c.Query("user_id")

	// 2、验证参数
	// user_id必须为正整数
	userID, err := strconv.Atoi(userId)
	if err != nil || userID < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	if userID == 0 {
		userID = int(auth.UserID)
	}
	// 3、具体业务
	// 3.1、查询用户的所有视频
	var videoList []models.Video
	if err := driver.Db.Debug().Model(models.Video{}).Where("author_id = ?", userID).Find(&videoList).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库出错" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 3.2、查询视频的作者，填充返回的视频信息
	returnVideoList := make([]Video, len(videoList))
	for i := 0; i < len(videoList); i++ {
		var author models.User
		var relation models.Relation
		returnVideoList[i].ID = videoList[i].ID
		returnVideoList[i].Title = videoList[i].Title
		returnVideoList[i].CommentCount = videoList[i].CommentCount
		returnVideoList[i].CoverURL = videoList[i].CoverURL
		returnVideoList[i].PlayURL = videoList[i].PlayURL
		returnVideoList[i].FavoriteCount = videoList[i].FavoriteCount
		returnVideoList[i].IsFavorite = false
		if err := driver.Db.Debug().Model(models.User{}).Where("id = ?", videoList[i].AuthorID).Find(&author).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库出错" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
		returnVideoList[i].Author = &User{
			ID:            author.ID,
			Name:          author.Name,
			FollowCount:   author.FollowCount,
			FollowerCount: author.FollowerCount,
			IsFollow:      false,
		}
		if err := driver.Db.Debug().Model(models.Relation{}).Where("user_id = ? and target_id=? and exist=1 and (type=1 or type=2)", auth.UserID, author.ID).Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库出错" + err.Error()
			c.JSON(409, failureResponse)
			return
		}
	}
	// 4、装配返回值
	successResponse.VideoList = returnVideoList
	c.JSON(http.StatusOK, successResponse)
	return
}
