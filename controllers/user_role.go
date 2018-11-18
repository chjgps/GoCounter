package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/permission"
)

type UserRoleController struct {
	BaseController
}

// 添加用户角色绑定关系
func (c *UserRoleController) AddUserRole() {
	obj := &permission.UserRole{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	userId := obj.UserId
	roleId := obj.RoleId

	// 查询是否存在
	acc, err := permission.UserRoleById(userId, roleId)
	if err != nil {
		if !permission.ErrUserRoleNotFound.Equal(err) {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// 不存在添加
	if acc == nil {
		obj.Status = 1
		obj.Created = timex.String()
		if err := permission.InsertUserRole(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	} else {
		// 存在修改
		acc.Status = 1
		acc.Updated = timex.String()
		if err := permission.UpdateUserRole(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除用户角色
func (c *UserRoleController) DelUserRole() {
	// userId
	userIdStr := c.Ctx.Input.Param(":userId")
	log.Debug(userIdStr)
	if len(userIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("user id is empty"))
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

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

	// 查询是否存在
	acc, err := permission.UserRoleById(userId, roleId)
	if err != nil {
		if !permission.ErrUserRoleNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	// 存在
	if err := permission.DelUserRole(acc.Id); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据角色ID查询所有用户
func (c *UserRoleController) UserByRoleId() {
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

	list, err := permission.UserByRoleId(roleId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}

// 根据用户查询所有角色
func (c *UserRoleController) RoleByUserId() {
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

	list, err := permission.RoleByUserId(userId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
