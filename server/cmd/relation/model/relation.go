package model

import (
	"douyin_rpc/common/cache"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

// Relation TODO
// 建立联合索引
type Relation struct {
	gorm.Model
	ID       int64 `gorm:"primarykey; not null"`
	UserID   int64 `gorm:"index; not null"`
	TargetID int64 `gorm:"index; not null"`
	Type     int   `gorm:"index; not null"` // 1:UserID关注TargetID   2:互相关注
	Exist    bool  // 判断是否存在，避免一个用户频繁关注取消关注造成数据膨胀，并且这样可以使用唯一索引
}

func (b *Relation) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const RelationCachePrefix string = "relation:relation_"

func queryRelationByUserID(tx *gorm.DB, userID int64) ([]Relation, error) {
	var relations []Relation
	if err := tx.Table("relation").Where("exist=1").Where("type=1 or type=2").Where("user_id=?", userID).Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}

func QueryRelationByUserIDWithCache(tx *gorm.DB, userID int64) ([]Relation, error) {
	key := RelationCachePrefix + "UserID_" + strconv.Itoa(int(userID))
	// 查看key是否存在
	//不存在

	var result string
	var relations []Relation
	var err error
	if !cache.Exist(key) {
		relations, err = queryRelationByUserID(tx, userID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, relations)
		if err != nil {
			return nil, err
		}
		return relations, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			relations, err = queryRelationByUserID(tx, userID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, relations)
			if err != nil {
				return nil, err
			}
			return relations, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &relations)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func queryRelationByTargetID(tx *gorm.DB, targetID int64) ([]Relation, error) {
	var relations []Relation
	if err := tx.Table("relation").Where("exist=1").Where("type=1 or type=2").Where("target_id=?", targetID).Find(&relations).Error; err != nil {
		return nil, err
	}
	return relations, nil
}
func QueryRelationIsFriend(tx *gorm.DB, userId int64) ([]Relation, error) {
	relations, err := QueryRelationByUserIDWithCache(tx, userId)
	if err != nil {
		return nil, err
	}
	returnR := make([]Relation, len(relations))
	j := 0
	for i := 0; i < len(relations); i++ {
		if relations[i].Type == 2 {
			returnR[j] = relations[i]
			j++
		}
	}
	return returnR[:j], nil
}
func QueryRelationByTargetIDWithCache(tx *gorm.DB, targetID int64) ([]Relation, error) {
	key := RelationCachePrefix + "TargetID_" + strconv.Itoa(int(targetID))
	// 查看key是否存在
	//不存在
	var result string
	var relations []Relation
	var err error
	if !cache.Exist(key) {
		relations, err = queryRelationByTargetID(tx, targetID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, relations)
		if err != nil {
			return nil, err
		}
		return relations, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			relations, err = queryRelationByTargetID(tx, targetID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, relations)
			if err != nil {
				return nil, err
			}
			return relations, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &relations)
	if err != nil {
		return nil, err
	}
	return relations, nil
}

func queryRelationByUserIDAndTargetID(tx *gorm.DB, userID, targetID int64) (*Relation, error) {
	var relations Relation
	if err := tx.Table("relation").Where("exist=1").Where("type=1 or type=2").Where("user_id=? and target_id=?", userID, targetID).Find(&relations).Error; err != nil {
		return nil, err
	}
	return &relations, nil
}
func CountRelation(tx *gorm.DB, userID int64) (friendCount, followCount, followerCount int64, err error) {
	if err = tx.Table("relation").Where("exist=1").Where("type=1 or type=2").Where("user_id=?", userID).Count(&followCount).Error; err != nil {
		return
	}
	if err = tx.Table("relation").Where("exist=1").Where("type=1 or type=2").Where("target_id=?", userID).Count(&followerCount).Error; err != nil {
		return
	}
	if err = tx.Table("relation").Where("exist=1").Where("type=2").Where("user_id=?", userID).Count(&friendCount).Error; err != nil {
		return
	}
	return
}
func QueryRelationByUserIDAndTargetIDWithCache(tx *gorm.DB, userID, targetID int64) (*Relation, error) {
	key := RelationCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_TargetID_" + strconv.Itoa(int(targetID))
	var result string
	var relation *Relation
	var err error
	// 查看key是否存在
	//不存在
	if !cache.Exist(key) {
		relation, err = queryRelationByUserIDAndTargetID(tx, userID, targetID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, relation)
		if err != nil {
			return nil, err
		}
		return relation, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			relation, err = queryRelationByUserIDAndTargetID(tx, userID, targetID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, relation)
			if err != nil {
				return nil, err
			}
			return relation, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &relation)
	if err != nil {
		return nil, err
	}
	return relation, nil
}

func UpdateRelation(tx *gorm.DB, relation Relation) error {
	if err := tx.Save(&relation).Error; err != nil {
		return err
	}
	return nil
}

func CreateRelation(tx *gorm.DB, relation Relation) error {
	if err := tx.Model(relation).Create(&relation).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOrCreateRelation(tx *gorm.DB, relation Relation) error {
	if err := tx.Clauses(clause.OnConflict{
		Columns:      []clause.Column{{Name: "user_id"}, {Name: "target_id"}},
		Where:        clause.Where{},
		TargetWhere:  clause.Where{},
		OnConstraint: "",
		DoNothing:    false,
		DoUpdates:    clause.Assignments(map[string]interface{}{"exist": true}),
		UpdateAll:    false,
	}).Create(&relation).Error; err != nil {
		return err
	}
	return nil
}
