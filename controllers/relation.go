package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//TODO
// 关注，粉丝，好友列表机器相似，可以整合为一个函数

type RelationList struct {
	utils.Response
	UserList []User `json:"user_list"`
}

// RelationAction
// @Summary 关注和取消关注操作
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Param action_type query int true "操作类型"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/relation/action [post]
func RelationAction(c *gin.Context) {
	successResponse := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	failureResponse := utils.Response{
		StatusCode: 1,
		StatusMsg:  "",
	}
	auth := Auth{}.GetAuth(c)
	//TODO
	// 业务代码
	// 1、解析参数
	toUserID := c.Query("to_user_id")
	actionType := c.Query("action_type")
	// 2、验证参数
	// 2.1 toUserID必须是整数
	id, err := strconv.Atoi(toUserID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "to_user_id必须为非0正整数"
		c.JSON(409, failureResponse)
		return
	}
	// 2.2 actionType 必须为1或2
	action, err := strconv.Atoi(actionType)
	if err != nil || (action != 1 && action != 2) {
		failureResponse.StatusMsg = "action必须为1或2"
		c.JSON(409, failureResponse)
		return
	}
	if auth.UserID == uint(id) {
		failureResponse.StatusMsg = "没办法关注自己或取消关注自己"
		c.JSON(409, failureResponse)
		return
	}
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 3、查询数据库获取两个用户信息，使用for update加锁（用户一般都存在）
	var user, targetUser models.User
	if err := driver.Db.Debug().Model(user).Set("gorm:query_option", "FOR UPDATE").Where("id = ?", auth.UserID).Find(&user).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	//TODO
	// 用户不存在的情况，布隆过滤器？
	if err := driver.Db.Debug().Model(targetUser).Set("gorm:query_option", "FOR UPDATE").Where("id = ?", id).Find(&targetUser).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 用户不存在
	if targetUser.ID == 0 {
		failureResponse.StatusMsg = "用户不存在"
		c.JSON(409, failureResponse)
		return
	}
	// 4、查询关系表，判断是否存在关系
	// 4.1 关系表的关系分析：
	// 		1）user关注targetUser，则表中存在一条 userID=user.ID，targetID=targetUser.ID，Type=1的数据
	//      2）targetUser关注user，则表中存在一条 userID=targetUser.ID，targetID=user.ID，Type=1的数据
	//      3）两个人互相关注，则存在两条数据，Type都等于2
	// 4.2 验证：
	//		1）如果actionType=1，且表中存在一条userID=user.ID，targetID=targetUser.ID，Type=1的数据，则报错，已关注
	//		2）如果actionType=2，且表中存在一条userID=targetUser.ID，targetID=user.ID，Type=1的数据，则报错，没有关注无需取消关注
	// 4.3查询两个关系数据
	var relation1, relation2 models.Relation
	if err := driver.Db.
		Debug().
		Model(relation1).
		Where("user_id = ?", user.ID).
		Where("target_id = ?", targetUser.ID).
		Find(&relation1).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库出错" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	if err := driver.Db.
		Debug().
		Model(relation2).
		Where("user_id = ?", targetUser.ID).
		Where("target_id = ?", user.ID).
		Find(&relation2).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库出错" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 由于不能确定这两个关系同时存在，因此不要使用for update加锁（使用for update时确保索引存在。不存在会锁住表）
	// for update在数据存在时加的是行级锁，不存在加的是间隙锁。之后进行insert时容易形成死锁
	if action == 1 {
		if relation1.ID == 0 {
			relation1.UserID = user.ID
			relation1.TargetID = targetUser.ID
			relation1.Exist = true

			//数据不存在，第一次关注，创建数据
			if relation2.ID != 0 && relation2.Exist && relation2.Type == 1 {
				// target关注了user，修改两个关系数据为Type=2
				relation2.Type = 2
				relation1.Type = 2
				if err := tx.Model(relation2).Updates(&relation2).Error; err != nil {
					failureResponse.StatusMsg = "更新数据库失败" + err.Error()
					c.JSON(409, failureResponse)
					tx.Rollback()
					return
				}

			} else {
				// 只需要增加一条Type为1的数据
				relation1.Type = 1

			}

			if err := tx.Model(relation1).Create(&relation1).Error; err != nil {
				failureResponse.StatusMsg = "创建关系失败" + err.Error()
				c.JSON(409, failureResponse)
				tx.Rollback()
				return
			}
		} else if relation1.Exist {
			// 关注操作且数据库显示已关注，错误
			failureResponse.StatusMsg = "关注信息已存在，您已关注无需再次关注"
			c.JSON(409, failureResponse)
			tx.Rollback()
			return

		} else {
			relation1.Exist = true
			// user关注过target，但是取消了，因此存在一条exist=false的数据，修改exist为true
			if relation2.ID != 0 && relation2.Exist && relation2.Type == 1 {
				// target关注了user，修改两个关系数据为Type=2
				relation2.Type = 2
				relation1.Type = 2
				if err := tx.Model(relation2).Updates(&relation2).Error; err != nil {
					failureResponse.StatusMsg = "更新数据库失败" + err.Error()
					c.JSON(409, failureResponse)
					tx.Rollback()
					return
				}
			} else {
				// 修改数据为Type=1
				relation1.Type = 1
			}

			if err := tx.Model(relation1).Updates(&relation1).Error; err != nil {
				failureResponse.StatusMsg = "更新数据库失败" + err.Error()
				c.JSON(409, failureResponse)
				tx.Rollback()
				return
			}
		}
		user.FollowCount += 1
		targetUser.FollowerCount += 1
	} else if action == 2 {
		// 取消关注，数据不存在直接报错
		if relation1.ID == 0 || !relation1.Exist {
			failureResponse.StatusMsg = "关注数据不存在，无需取消关注"
			c.JSON(409, failureResponse)
			tx.Rollback()
			return
		} else {
			//数据存在
			if relation1.Type == 1 {
				// Type为1，只需要将Exist改为false
				relation1.Exist = false
				if err := tx.Model(relation1).Save(&relation1).Error; err != nil {
					failureResponse.StatusMsg = "更新数据库失败" + err.Error()
					c.JSON(409, failureResponse)
					tx.Rollback()
					return
				}
			} else {
				// Type为2，修改relation2的Type为1
				relation1.Exist = false
				relation1.Type = 1
				if err := tx.Model(relation1).Save(&relation1).Error; err != nil {
					failureResponse.StatusMsg = "更新数据库失败" + err.Error()
					c.JSON(409, failureResponse)
					tx.Rollback()
					return
				}
				relation2.Type = 1
				if err := tx.Model(relation2).Updates(&relation2).Error; err != nil {
					failureResponse.StatusMsg = "更新数据库失败" + err.Error()
					c.JSON(409, failureResponse)
					tx.Rollback()
					return
				}
			}

		}
		user.FollowCount -= 1
		targetUser.FollowerCount -= 1
	}
	//注意使用gorm有可能修改到零值的需要使用Save而不能使用updates
	// 5、修改用户的关注数和粉丝数
	if err := tx.Model(user).Save(&user).Error; err != nil {
		failureResponse.StatusMsg = "更新数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		tx.Rollback()
		return
	}
	if err := tx.Model(targetUser).Save(&targetUser).Error; err != nil {
		failureResponse.StatusMsg = "更新数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		tx.Rollback()
		return
	}
	if err := tx.Commit().Error; err != nil {
		failureResponse.StatusMsg = "更新数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		tx.Rollback()
		return
	}
	c.JSON(http.StatusOK, successResponse)
}

// RelationFollowList
// @Summary 获取关注列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/follow/list [get]
func RelationFollowList(c *gin.Context) {

	auth := Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	// 3、查询数据库获取该用的所有关注
	var relations []models.Relation
	if err := driver.Db.Debug().Model(models.Relation{}).Where("user_id = ?", id).Where("exist=1").Where("type=1 or type=2").Find(&relations).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	// 4、构建User信息
	//TODO
	// 目前循环查找数据库，之后使用联表查询优化，缓存优化等等
	// 关系表需要添加索引的地方：联合索引userID和targetID以及exist
	userList := make([]User, len(relations))

	for i := 0; i < len(relations); i++ {
		var user models.User
		var relation models.Relation
		// 此处为TargetID
		userList[i].ID = relations[i].TargetID
		if err := driver.Db.Debug().Model(user).Where("id = ?", userList[i].ID).Find(&user).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		userList[i].Name = user.Name
		userList[i].FollowerCount = user.FollowerCount
		userList[i].FollowCount = user.FollowCount
		userList[i].IsFollow = false
		if err := driver.Db.Debug().Model(relation).Where("user_id = ?", auth.UserID).Where("target_id = ?", relations[i].TargetID).Where("exist=1").Where("type=1 or type=2").Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		// 再次判断是否存在
		if relation.ID != 0 && relation.Exist {
			userList[i].IsFollow = true
		}

	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}

// RelationFollowerList
// @Summary 获取关注者列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/follower/list [get]
func RelationFollowerList(c *gin.Context) {

	auth := Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	// 3、查询数据库获取该用的所有关注
	var relations []models.Relation
	if err := driver.Db.Debug().Model(models.Relation{}).Where("target_id = ?", id).Where("exist=1").Where("type=1 or type=2").Find(&relations).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	// 4、构建User信息
	//TODO
	// 目前循环查找数据库，之后使用联表查询优化，缓存优化等等
	// 关系表需要添加索引的地方：联合索引userID和targetID以及exist
	userList := make([]User, len(relations))

	for i := 0; i < len(relations); i++ {
		var user models.User
		var relation models.Relation
		// 此处为UserID
		userList[i].ID = relations[i].UserID
		if err := driver.Db.Debug().Model(user).Where("id = ?", userList[i].ID).Find(&user).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		userList[i].Name = user.Name
		userList[i].FollowerCount = user.FollowerCount
		userList[i].FollowCount = user.FollowCount
		userList[i].IsFollow = false
		if err := driver.Db.Debug().Model(relation).Where("user_id = ?", auth.UserID).Where("target_id = ?", relations[i].UserID).Where("exist=1").Where("type=1 or type=2").Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		// 再次判断是否存在
		if relation.ID != 0 && relation.Exist {
			userList[i].IsFollow = true
		}

	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}

// RelationFriendList
// @Summary 获取好友列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/friend/list [get]
func RelationFriendList(c *gin.Context) {

	auth := Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	// 3、查询数据库获取该用的所有关注
	var relations []models.Relation
	if err := driver.Db.Debug().Model(models.Relation{}).Where("user_id = ?", id).Where("exist=1").Where("type=2").Find(&relations).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	// 4、构建User信息
	//TODO
	// 目前循环查找数据库，之后使用联表查询优化，缓存优化等等
	// 关系表需要添加索引的地方：联合索引userID和targetID以及exist
	userList := make([]User, len(relations))
	for i := 0; i < len(relations); i++ {
		var user models.User
		var relation models.Relation
		userList[i].ID = relations[i].TargetID
		if err := driver.Db.Debug().Model(user).Where("id = ?", userList[i].ID).Find(&user).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		userList[i].Name = user.Name
		userList[i].FollowerCount = user.FollowerCount
		userList[i].FollowCount = user.FollowCount
		userList[i].IsFollow = false
		if err := driver.Db.Debug().Model(relation).Where("user_id = ?", auth.UserID).Where("target_id = ?", relations[i].TargetID).Where("exist=1").Where("type=1 or type=2").Find(&relation).Error; err != nil {
			failureResponse.StatusMsg = "查询数据库失败" + err.Error()
			c.JSON(http.StatusOK, failureResponse)
			return
		}
		// 再次判断是否存在
		if relation.ID != 0 && relation.Exist {
			userList[i].IsFollow = true
		}

	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}
