package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/permission"
)

type RoleController struct {
	BaseController
}

// 添加
func (c *RoleController) AddRole() {
	obj := &permission.Role{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	name := obj.Name

	// 查询用户名是否重复
	acc, err := permission.RoleByName(name)
	if err != nil {
		if !permission.ErrRoleNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(permission.ErrNameAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := permission.InsertRole(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, obj, nil)
	return
}

// 修改
func (c *RoleController) EditRole() {
	obj := &permission.Role{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	roleId := obj.Id

	// 查询是否存在
	acc, err := permission.RoleById(roleId)
	if err != nil {
		if !permission.ErrRoleNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(permission.ErrRoleNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := permission.UpdateRole(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *RoleController) DelRole() {
	roleIdStr := c.Ctx.Input.Param(":id")
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

	if err := permission.DelRole(roleId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *RoleController) RoleById() {
	roleIdStr := c.Ctx.Input.Param(":id")
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

	acc, err := permission.RoleById(roleId)
	if err != nil {
		if !permission.ErrRoleNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有
func (c *RoleController) RoleList() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	name := c.GetString("name")

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

	total, list, err := permission.RoleList(map[string]interface{}{
		"startDate": startDate,
		"endDate":   endDate,
		"name":      name,
	}, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	var data interface{}
	if list == nil {
		data = make([]interface{}, 0)
	} else {
		data = list
	}

	c.WriteHttpResponse(200, struct {
		Total int64       `json:"total"`
		Data  interface{} `json:"data"`
	}{
		Total: total,
		Data:  data,
	}, nil)

	return
}
