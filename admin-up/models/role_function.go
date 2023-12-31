package models

import (
	"admin-up/helper"
	"fmt"

	"gorm.io/gorm"
)

type RoleFunction struct {
	gorm.Model
	RoleId     uint `gorm:"column:role_id;type:int(11);" json:"role_id"`         // 角色ID
	FunctionId uint `gorm:"column:function_id;type:int(11);" json:"function_id"` // 功能ID
}

func (table *RoleFunction) TableName() string {
	return "role_function"
}

// GetRoleFunctions 获取指定角色授权的功能列表
// roleIdentity 角色唯一标识
// isAdmin 是否是超管
func GetRoleFunctions(roleIdentity string, isAdmin bool) *gorm.DB {
	tx := DB.Model(new(FunctionBasic)).Select("function_basic.identity, mb.identity menu_identity, " +
		"function_basic.name, function_basic.uri, function_basic.sort").
		Joins("LEFT JOIN menu_basic mb ON mb.id = function_basic.menu_id")
	if !isAdmin {
		// 1. 获取角色ID
		var roleId int
		err := DB.Model(new(RoleBasic)).Select("id").Where("identity = ?", roleIdentity).
			Scan(&roleId).Error
		if err != nil {
			helper.Error("[DB ERROR] get role function err: %v", err)
		}
		// 2. 获取当前角色授权的功能
		tx.Joins("LEFT JOIN role_function rf ON rf.function_id = function_basic.id").
			Where("rf.role_id = ?", roleId)
	}
	tx.Order("function_basic.sort ASC")
	return tx
}

type authFuncReply struct {
	Uri string `json:"uri"`
}

// GetAuthFuncUri 获取给特定角色授权功能的URI
func GetAuthFuncUri(roleIdentity string) (map[string]interface{}, error) {
	roleBasic := new(RoleBasic)
	err := DB.Model(new(RoleBasic)).Select("id").Where("identity = ?", roleIdentity).Find(roleBasic).Error
	if err != nil {
		return nil, err
	}
	afr := make([]*authFuncReply, 0)
	err = DB.Model(new(RoleFunction)).Select("fb.uri").
		Joins("LEFT JOIN function_basic fb ON fb.id = role_function.function_id").
		Where("role_function.role_id = ?", roleBasic.ID).Find(&afr).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("afr: %v\n", afr)
	data := make(map[string]interface{})
	for _, v := range afr {
		data[v.Uri] = "1"
	}
	return data, nil
}
