package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/permission"
)

type PermissionController struct {
	BaseController
}

// 添加
func (c *PermissionController) AddPermission() {
	obj := &permission.Permission{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	tag := obj.Tag

	// 查询名称是否重复
	acc, err := permission.PermissionByTag(tag)
	if err != nil {
		if !permission.ErrPermissionNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(permission.ErrPermissionAlreadyExist))
		return
	}

	obj.Created = timex.String()
	if err := permission.InsertPermission(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *PermissionController) EditPermission() {
	obj := &permission.Permission{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	permissionId := obj.Id

	// 查询用户是否存在
	acc, err := permission.PermissionById(permissionId)
	if err != nil {
		if !permission.ErrPermissionNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(permission.ErrPermissionNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := permission.UpdatePermission(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *PermissionController) DelPermission() {
	permissionIdStr := c.Ctx.Input.Param(":id")
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

	if err := permission.DelPermission(permissionId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *PermissionController) PermissionById() {
	permissionIdStr := c.Ctx.Input.Param(":id")
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

	acc, err := permission.PermissionById(permissionId)
	if err != nil {
		if !permission.ErrPermissionNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有
func (c *PermissionController) PermissionList() {
	list, err := permission.PermissionList()
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
