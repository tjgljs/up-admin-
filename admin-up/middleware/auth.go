package middleware

import (
	"admin-up/define"
	"admin-up/helper"
	"admin-up/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoginAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaim, err := helper.AnalyzeToken(c.GetHeader("AccessToken"))
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"code": 60403,
				"msg":  "登陆过期，请重新登陆",
			})
		} else {
			if userClaim.RoleIdentity == "" {
				c.Abort()
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "非法请求",
				})
			}
			//判断是不是超级管理员
			isAdmin, err := models.RDB.Get(context.Background(), define.RedisRoleAdminPrefix+userClaim.RoleIdentity).Result()
			adminRoleKey := define.RedisRoleAdminPrefix + userClaim.RoleIdentity
			if err != nil {
				roleBasic := new(models.RoleBasic)
				err = models.DB.Model(new(models.RoleBasic)).Select("is_admin").Where("identity= ?", userClaim.RoleIdentity).Find(roleBasic).Error
				if err != nil {
					helper.Error("[DB ERROR]:%v", err)
					c.Abort()
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "网络异常，请重试",
					})
				}
				//加入redis缓存，默认保存一周
				models.RDB.Set(context.Background(), adminRoleKey, roleBasic.IsAdmin, time.Second*3600*24*7)
				if roleBasic.IsAdmin == 1 {
					isAdmin = "1"
				} else {
					isAdmin = "0"
				}
			} else {
				//查到数据，本身就在redis缓存中，再增加一周的缓存时间
				models.RDB.Expire(context.Background(), adminRoleKey, time.Second*3600*24*7)

			}
			if isAdmin == "1" {
				userClaim.IsAdmin = true
			} else {
				userClaim.IsAdmin = false
			}
			c.Set("UserClaim", userClaim)
			c.Next()
		}
	}
}

// FuncAuthCheck 功能权限验证
// 判断 用户是否有调用该路由的权力
// FuncAuthCheck 功能权限验证
func FuncAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 不是超管时，判断用户是否具有相关权限
		userClaim := c.MustGet("UserClaim").(*define.UserClaim)
		// 可操作函数的 key
		funcKey := define.RedisFuncPrefix + userClaim.RoleIdentity
		fmt.Printf("funcKey: %v\n", funcKey)
		fmt.Printf("userClaim.IsAdmin: %v\n", userClaim.IsAdmin)
		if !userClaim.IsAdmin {
			// 判断key是否在Redis中存在
			keyExist, _ := models.RDB.Exists(context.Background(), funcKey).Result()

			fmt.Printf("keyExist: %v\n", keyExist)
			if keyExist > 0 {
				// key存在，再续一周
				models.RDB.Expire(context.Background(), funcKey, time.Second*3600*24*7)

				fieldExist, _ := models.RDB.HExists(context.Background(), funcKey, c.Request.RequestURI).Result()
				fmt.Printf("c.Request.RequestURI: %v\n", c.Request.RequestURI)
				fmt.Printf("fieldExist: %v\n", fieldExist)
				if !fieldExist { // 权限不存在，非法访问
					c.Abort()
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "非法请求",
					})
				}
			} else {
				// key不存在，从DB中查询数据，并保存

				data, err := models.GetAuthFuncUri(userClaim.RoleIdentity)
				fmt.Printf("data: %v\n", data)
				if err != nil {
					helper.Error("[DB ERROR] : %v", err)
					c.Abort()
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "网络异常",
					})
				}
				if len(data) == 0 {
					data["up-admin-empty"] = "get"
				}
				fmt.Printf("data: %v\n", data)
				//roleidentity  url  1
				err = models.RDB.HSet(context.Background(), funcKey, data).Err()
				fmt.Printf("data: %v\n", data)
				if err != nil {
					helper.Error("[DB ERROR] : %v", err)
					c.Abort()
					c.JSON(http.StatusOK, gin.H{
						"code": -1,
						"msg":  "网络异常",
					})
				}

				models.RDB.Expire(context.Background(), funcKey, time.Second*3600*24*7)
			}
		}
		c.Next()
	}
}
