package service

import (
	"admin-up/define"
	"admin-up/helper"
	"admin-up/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 角色列表
func SetRoleList(c *gin.Context) {
	in := &SetRoleListRequest{NewQueryRequest()}
	err := c.ShouldBindQuery(in)
	if err != nil {
		helper.Info("[INPUT ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	var (
		cnt  int64
		list = make([]*SetRoleListReply, 0)
	)
	err = models.GetRoleList(in.Keyword).Count(&cnt).Offset((in.Page - 1) * in.Size).Limit(in.Size).Find(&list).Error
	if err != nil {
		helper.Info("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	for _, v := range list {
		v.CreatedAt = helper.RFC3339ToNormalTime(v.CreatedAt)
		v.UpdatedAt = helper.RFC3339ToNormalTime(v.UpdatedAt)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"data": gin.H{
			"list":  list,
			"count": cnt,
		},
	})
}

// 获取角色详情
func SetRoleDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "必填参数不能为空",
		})
		return
	}
	data := new(SetRoleDetailReply)
	//获取角色的基本信息
	rb, err := models.GetRoleBasic(identity)
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据异常",
		})
		return
	}
	data.Name = rb.Name
	data.IsAdmin = rb.IsAdmin
	data.Sort = rb.Sort

	//获取授权的菜单
	menuIdentities, err := models.GetRoleMenuIdentity(rb.ID, rb.IsAdmin == 1)
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	data.MenuIdentities = menuIdentities
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": data,
	})

}

// 删除角色相关的缓存数据
func redisRoleDelete(roleIdentity string) error {
	err := models.RDB.Del(context.Background(), define.RedisRoleAdminPrefix).Err()
	if err != nil {
		return err
	}
	err = models.RDB.HDel(context.Background(), define.RedisMenuPrefix, roleIdentity).Err()
	if err != nil {
		return err
	}
	err = models.RDB.Del(context.Background(), define.RedisFuncPrefix+roleIdentity).Err()
	if err != nil {
		return err
	}
	return nil
}

// 角色修改
// SetRoleUpdate 角色修改
func SetRoleUpdate(c *gin.Context) {
	in := new(SetRoleUpdateRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Info("[INPUT ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	var (
		menuIds = make([]uint, 0)
		rb      = new(models.RoleBasic)
	)
	// 获取菜单ID
	err = models.DB.Model(new(models.MenuBasic)).Select("id").
		Where("identity IN ?", in.MenuIdentities).Scan(&menuIds).Error
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	// 角色ID
	err = models.DB.Model(new(models.RoleBasic)).Select("id").
		Where("identity = ?", in.Identity).Find(rb).Error
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	// 组装数据
	// 授权的菜单
	rms := make([]*models.RoleMenu, len(menuIds))
	for i, _ := range rms {
		rms[i] = &models.RoleMenu{
			RoleId: rb.ID,
			MenuId: menuIds[i],
		}
	}
	// Redis Key 删除
	err = redisRoleDelete(in.Identity)
	if err != nil {
		helper.Info("[RDB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "缓存异常",
		})
		return
	}
	// 修改数据
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 更新角色
		err = tx.Model(new(models.RoleBasic)).Where("id = ?", rb.ID).Updates(map[string]interface{}{
			"name":     in.Name,
			"is_admin": in.IsAdmin,
			"sort":     in.Sort,
		}).Error
		if err != nil {
			return err
		}
		// 删除老数据
		// 授权的菜单
		err = tx.Where("role_id = ?", rb.ID).Delete(new(models.RoleMenu)).Error
		if err != nil {
			return err
		}
		// 增加新数据
		// 授权的菜单
		if len(rms) > 0 {
			err = tx.Create(rms).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	// 延迟双删
	go func() {
		time.Sleep(time.Millisecond * 200)
		err = redisRoleDelete(in.Identity)
		if err != nil {
			helper.Info("[RDB ERROR] : %v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "缓存异常",
			})
			return
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
