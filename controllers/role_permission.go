package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/permission"
)

type RolePermissionController struct {
	BaseController
}

// 添加绑定关系
func (c *RolePermissionController) AddRolePermission() {
	obj := &permission.RolePermission{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	roleId := obj.RoleId
	permissionId := obj.PermissionId

	// 查询是否存在
	acc, err := permission.RolePermissionById(roleId, permissionId)
	if err != nil {
		if !permission.ErrRolePermissionNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	// 不存在添加
	if acc == nil {
		obj.Status = 1
		obj.Created = timex.String()
		if err := permission.InsertRolePermission(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	} else {
		// 存在修改
		acc.Status = 1
		acc.Updated = timex.String()
		if err := permission.UpdateRolePermission(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *RolePermissionController) DelRolePermission() {
	// roleId
	roleIdStr := c.Ctx.Input.Param(":roleId")
	log.Debug(roleIdStr)
	if len(roleIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("role id is empty"))
		return
	}

	roleId, err := strconv.Atoi(roleIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// permissionId
	permissionIdStr := c.Ctx.Input.Param(":permissionId")
	log.Debug(permissionIdStr)
	if len(permissionIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("permission id is empty"))
		return
	}

	permissionId, err := strconv.Atoi(permissionIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// 查询是否存在
	acc, err := permission.RolePermissionById(roleId, permissionId)
	if err != nil {
		if !permission.ErrRolePermissionNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	// 存在
	if err := permission.DelRolePermission(acc.Id); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据角色ID查询所有权限
func (c *RolePermissionController) PermissionByRoleId() {
	roleId, err := c.GetInt("roleId")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	page, err := c.GetInt("page")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	pageSize, err := c.GetInt("pageSize")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	list, err := permission.PermissionByRoleId(roleId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}

// 根据用户ID查询所有权限
func (c *RolePermissionController) PermissionByUserId() {
	userId, err := c.GetInt("userId")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	page, err := c.GetInt("page")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	pageSize, err := c.GetInt("pageSize")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// 根据用户查询角色
	roles, err := permission.RoleByUserId(userId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(roles) == 0 {
		c.WriteHttpResponse(404, nil, errors.New("role not exist"))
		return
	}

	roleId := roles[0].Id

	list, err := permission.PermissionByRoleId(roleId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
