package service

import (
	"admin-up/define"
	"admin-up/helper"
	"admin-up/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 新增菜单
func DevMenuAdd(c *gin.Context) {
	in := new(DevMenuAddRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[CINDJSON ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	//获取父级id
	var parentId uint
	if in.ParentIdentity != "" {
		err = models.DB.Model(new(models.MenuBasic)).Select("id").Where("identity= ?", in.ParentIdentity).Scan(&parentId).Error
		if err != nil {
			helper.Error("[DB ERROR]:%v", err)
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "数据库异常",
			})
			return
		}
		//保存数据
		err = models.DB.Create(&models.MenuBasic{
			Identity: helper.UUID(),
			ParentId: parentId,
			Name:     in.Name,
			WebIcon:  in.WebIcon,
			Sort:     in.Sort,
			Path:     in.Path,
			Level:    in.Level,
		}).Error
		if err != nil {
			helper.Error("[DB ERROR]:%v", err)
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "数据库异常",
			})
			return
		}
		//清空菜单缓存
		err = models.RDB.Del(context.Background(), define.RedisMenuPrefix).Err()
		if err != nil {
			helper.Error("[RDB ERROR]:%v", err)
			c.JSON(200, gin.H{
				"code": -1,
				"msg":  "缓存异常",
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "新增成功",
		})
	}

}

// 修改菜单
func DevMenuUpdate(c *gin.Context) {
	in := new(DevMenuUpdateRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[BindJSON ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	if in.Identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
		})
		return
	}
	//获取父级id
	var parentId uint
	if in.ParentIdentity != "" {
		err = models.DB.Model(new(models.MenuBasic)).Select("id").
			Where("identity = ?", in.ParentIdentity).Scan(&parentId).Error
		if err != nil {
			helper.Error("[DB ERROR] : %v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "数据库异常",
			})
			return
		}
	}
	//更新数据
	err = models.DB.Model(new(models.MenuBasic)).Where("identity= ?", in.Identity).Updates(map[string]interface{}{
		"parent_id": parentId,
		"name":      in.Name,
		"web_icon":  in.WebIcon,
		"sort":      in.Sort,
		"path":      in.Path,
		"level":     in.Level,
	}).Error
	if err != nil {
		helper.Error("[ED ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库错误",
		})
		return
	}
	//清空菜单缓存
	err = models.RDB.Del(context.Background(), define.RedisMenuPrefix).Err()
	if err != nil {
		helper.Error("[RDB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "缓存异常",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})

}

// 删除菜单
func DevMenuDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数为空",
		})
		return
	}
	//删除数据库中的数据
	err := models.DB.Where("identity= ?", identity).Delete(new(models.MenuBasic)).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "缓存异常",
		})
		return
	}
	//清空缓存数据
	err = models.RDB.Del(context.Background(), define.RedisMenuPrefix).Err()
	if err != nil {
		helper.Error("[RDB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "缓存错误",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// 新增功能
func DevFuncAdd(c *gin.Context) {
	in := new(DevFuncAddRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[BIND ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	if in.MenuIdentity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
		})
		return
	}
	//获取菜单的id
	var menuId uint
	err = models.DB.Model(new(models.MenuBasic)).Where("identity= ?", in.MenuIdentity).Select("id").Scan(&menuId).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库错误",
		})
		return
	}
	//保存数据
	err = models.DB.Create(&models.FunctionBasic{
		Identity: helper.UUID(),
		MenuId:   menuId,
		Name:     in.Name,
		Uri:      in.Uri,
		Sort:     in.Sort,
	}).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "新增成功",
	})
}

// 功能修改
func DevFuncUpdate(c *gin.Context) {
	in := new(DevFuncUpdateRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[BIND ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	if in.Identity == "" || in.MenuIdentity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
		})
		return
	}
	//获取菜单id
	var menuId uint
	err = models.DB.Model(new(models.MenuBasic)).Select("id").Where("idnetity= ?", in.MenuIdentity).Scan(&menuId).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	//更新数据
	err = models.DB.Model(new(models.FunctionBasic)).Where("identity= ?", in.Identity).Updates(map[string]interface{}{
		"menu_id": menuId,
		"name":    in.Name,
		"uri":     in.Uri,
		"sort":    in.Sort,
	}).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// 删除功能
func DevFuncDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数不能为空",
		})
		return
	}
	err := models.DB.Where("identity= ?", identity).Delete(new(models.FunctionBasic)).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
