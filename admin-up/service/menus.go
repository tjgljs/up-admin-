package service

import (
	"admin-up/define"
	"admin-up/helper"
	"admin-up/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Menus 获取菜单列表
func Menus(c *gin.Context) {
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
	menuField := userClaim.RoleIdentity
	if userClaim.IsAdmin {
		menuField = "ADMIN"
	}
	data := make([]*MenuReply, 0)
	roleMenus := make([]*RoleMenu, 0)
	res, err := models.RDB.HGet(context.Background(), define.RedisMenuPrefix, menuField).Result()
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		tx, err := models.GetRoleMenus(userClaim.RoleIdentity, userClaim.IsAdmin)
		if err != nil {
			helper.Error("[DB ERROR]:%v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "数据异常",
			})
			return
		}
		err = tx.Find(&roleMenus).Error
		if err != nil {
			helper.Error("[DB ERROR]:%v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "数据异常",
			})
			return
		}
		data = roleMenuToMenuReply(roleMenus)
		b, _ := json.Marshal(data)
		models.RDB.HSet(context.Background(), define.RedisMenuPrefix, map[string]string{
			menuField: string(b),
		})

	} else {
		err = json.Unmarshal([]byte(res), &data)
		if err != nil {
			helper.Error("[DB ERROR]:%v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "数据异常",
			})
			return
		}
		//redis查到数据后，再续两周
		models.RDB.Expire(context.Background(), define.RedisMenuPrefix, time.Second*3600*24*14)

	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"data": data,
	})

}

func roleMenuToMenuReply(roleMenus []*RoleMenu) []*MenuReply {
	reply := make([]*MenuReply, 0)
	//一级循环，得到顶级菜单
	for _, v := range roleMenus {
		if v.ParentId == 0 {
			reply = append(reply, &MenuReply{
				Identity: v.Identity,
				Name:     v.Name,
				WebIcon:  v.WebIcon,
				Sort:     v.Sort,
				Path:     v.Path,
				Level:    v.Level,
				SubMenus: getChildrenMenu(v.Id, v.Identity, roleMenus),
			})
		}

	}
	return reply
}

func getChildrenMenu(parentId int, parentIdentity string, roleMenus []*RoleMenu) []*MenuReply {
	data := make([]*MenuReply, 0)
	for _, v := range roleMenus {
		if v.ParentId == parentId {
			data = append(data, &MenuReply{
				Identity:       v.Identity,
				Name:           v.Name,
				WebIcon:        v.WebIcon,
				Sort:           v.Sort,
				Path:           v.Path,
				Level:          v.Level,
				ParentIdentity: parentIdentity,
				SubMenus:       getChildrenMenu(v.Id, v.Identity, roleMenus),
			})
		}
	}
	return data
}

// SetMenuList 设置-获取菜单列表
func SetMenuList(c *gin.Context) {
	Menus(c)
}
