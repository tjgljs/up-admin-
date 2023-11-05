package service

import (
	"admin-up/define"
	"admin-up/helper"
	"admin-up/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)
	type UserInfoReply struct {
		Username string `json:"username"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		RoleName string `json:"role_name"`
	}
	data := new(UserInfoReply)
	err := models.GetUserInfo(userClaim.Identity).Find(data).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "网络异常",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"data": data,
	})
}

// 修改密码
func UserPasswordChange(c *gin.Context) {
	in := new(UserPasswordChangeRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[BindJSON ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	userClaim := c.MustGet("UserClaim").(*define.UserClaim)

	//判断旧密码是否正确
	var cnt int64
	err = models.DB.Model(new(models.UserBasic)).Where("identity = ? AND password= ?", userClaim.Identity, in.OldPassword).Count(&cnt).Error
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "网络异常",
		})
		return
	}
	if cnt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "旧密码不正确",
		})
		return
	}
	//修改密码
	err = models.DB.Model(new(models.UserBasic)).Where("identity= ?", userClaim.Identity).Update("password", in.NewPassword).Error
	if err != nil {
		helper.Error("[DB ERROR] : %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "网络异常",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

func SetUserList(c *gin.Context) {
	in := &SetUserListRequest{NewQueryRequest()}
	err := c.ShouldBindQuery(in)
	if err != nil {
		helper.Error("[INPUT ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}

	var (
		cnt  int64
		list = make([]*SetUserListReply, 0)
	)
	err = models.GetUserList(in.Keyword).Count(&cnt).Offset((in.Page - 1) * in.Size).Limit(in.Size).Find(&list).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库错误",
		})
		return
	}
	for _, v := range list {
		v.CreatedAt = helper.RFC3339ToNormalTime(v.CreatedAt)
		v.UpdatedAt = helper.RFC3339ToNormalTime(v.UpdatedAt)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "加载成功",
		"data": gin.H{
			"list":  list,
			"count": cnt,
		},
	})
}

// 管理员的创建
func SetUserAdd(c *gin.Context) {
	in := new(SetUserAddRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[INPUT ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	//判断用户是否已经存在
	var cnt int64
	err = models.DB.Model(new(models.UserBasic)).Where("username= ?", in.Username).Count(&cnt).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	if cnt > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户已经存在",
		})
		return
	}
	//创建用户
	err = models.DB.Create(&models.UserBasic{
		Identity:     helper.UUID(),
		RoleIdentity: in.RoleIdentity,
		Username:     in.Username,
		Password:     in.Password,
		Phone:        in.Phone,
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
		"msg":  "创建成功",
	})
}

// 管理员修改
func SetUserUpdate(c *gin.Context) {
	in := new(SetUserUpdateRequest)
	err := c.ShouldBindJSON(in)
	if err != nil {
		helper.Error("[INPUT ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "参数异常",
		})
		return
	}
	//判断用户是否存在
	var cnt int64
	err = models.DB.Model(new(models.UserBasic)).Where("identity != ? AND username = ?", in.Identity, in.Username).Count(&cnt).Error
	if err != nil {
		helper.Error("[DB ERROR]:%v", err)
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "数据库异常",
		})
		return
	}
	if cnt > 0 {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "用户名已经存在",
		})
		return
	}
	//修改数据
	err = models.DB.Model(new(models.UserBasic)).Where("identity= ?", in.Identity).Updates(map[string]interface{}{
		"role_identity": in.Identity,
		"username":      in.Username,
		"phone":         in.Phone,
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
		"code": -1,
		"msg":  "修改成功",
	})
}

// 删除管理员
func SetUserDelete(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(200, gin.H{
			"code": -1,
			"msg":  "必填参数为空",
		})
		return
	}
	//删除管理员
	err := models.DB.Where("identity= ?", identity).Delete(new(models.UserBasic)).Error
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
