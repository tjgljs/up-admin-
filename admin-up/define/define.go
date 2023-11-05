package define

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	// DefaultSize 默认每页查询20条数据
	DefaultSize    = 20
	FrameName      = "UpAdmin"
	DateTimeLayout = "2006-01-02 15:04:05"
	JwtKey         = "up-admin"
	// TokenExpire token 有效期，7天
	TokenExpire = time.Now().Add(time.Second * 3600 * 24 * 7).Unix()
	// RefreshTokenExpire 刷新 token 有效期，14天
	RefreshTokenExpire = time.Now().Add(time.Second * 3600 * 24 * 14).Unix()

	// RedisRoleAdminPrefix 判断角色是否是超管的前缀
	RedisRoleAdminPrefix = "ADMIN-"
	// RedisMenuPrefix 菜单
	RedisMenuPrefix = "MENU"
	// RedisFuncPrefix 功能的前缀
	RedisFuncPrefix = "FUNC-"
)

type UserClaim struct {
	Id           uint
	Identity     string
	Name         string
	RoleIdentity string // 角色唯一标识
	IsAdmin      bool   // 是否是超管
	jwt.StandardClaims
}
