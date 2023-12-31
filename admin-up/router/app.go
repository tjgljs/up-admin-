package router

import (
	"admin-up/middleware"
	"admin-up/service"

	"github.com/gin-gonic/gin"
)

func App() *gin.Engine {
	r := gin.Default()
	//登陆
	//用户名，密码登陆
	r.POST("/login/password", service.LoginPassword)

	//登陆信息的认证
	loginAuth := r.Group("/").Use(middleware.LoginAuthCheck())
	//获取菜单列表
	loginAuth.GET("/menus", service.Menus)

	// 获取功能列表
	loginAuth.GET("/functions", service.Functions)

	//获取用户信息
	loginAuth.GET("/user/info", service.UserInfo)

	//修改密码
	loginAuth.PUT("/user/password/change", service.UserPasswordChange)

	//url鉴权的接口
	auth := loginAuth.Use(middleware.FuncAuthCheck())

	//角色管理
	//角色列表
	auth.GET("/set/role/list", service.SetRoleList)
	//角色详情
	auth.GET("/set/role/detail", service.SetRoleDetail)
	//角色修改
	auth.PUT("/set/role/update", service.SetRoleUpdate)

	// 菜单功能管理
	// 菜单列表
	auth.GET("/set/menu/list", service.SetMenuList)
	// 功能列表
	auth.GET("/set/func/list", service.SetFuncList)
	// 管理员管理
	// 管理员列表
	auth.GET("/set/user/list", service.SetUserList)
	// 新增管理员
	auth.POST("/set/user/add", service.SetUserAdd)
	// 修改管理员
	auth.PUT("/set/user/update", service.SetUserUpdate)
	// 删除管理员
	auth.DELETE("/set/user/delete", service.SetUserDelete)

	// ---------------- END - 设置 ----------------

	// ---------------- BEGIN - dev ----------------
	// 新增菜单
	auth.POST("/dev/menu/add", service.DevMenuAdd)
	// 修改菜单
	auth.PUT("/dev/menu/update", service.DevMenuUpdate)
	// 删除菜单
	auth.DELETE("/dev/menu/delete", service.DevMenuDelete)
	// 新增功能
	auth.POST("/dev/func/add", service.DevFuncAdd)
	// 修改功能
	auth.PUT("/dev/func/update", service.DevFuncUpdate)
	// 删除功能
	auth.DELETE("/dev/func/delete", service.DevFuncDelete)
	// ---------------- END - dev ----------------
	return r
}
