package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/account"
	"github.com/beego/ms304w-client/models/material"
)

type GroupController struct {
	BaseController
}

// 添加
func (c *GroupController) AddGroup() {
	obj := &account.Group{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	name := obj.Name

	// 查询名称是否重复
	acc, err := account.GroupByName(name)
	if err != nil {
		if !account.ErrGroupNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(account.ErrNameAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := account.InsertGroup(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *GroupController) EditGroup() {
	obj := &account.Group{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	groupId := obj.Id

	// 查询是否存在
	acc, err := account.GroupById(groupId)
	if err != nil {
		if !account.ErrGroupNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(account.ErrGroupNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := account.UpdateGroup(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *GroupController) DelGroup() {
	groupIdStr := c.Ctx.Input.Param(":id")
	log.Debug(groupIdStr)
	if len(groupIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("group id is empty"))
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// 查询组是否有关联用户
	accounts, err := account.AccountByGroupId(groupId, 1, 10)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(accounts) > 0 {
		c.WriteHttpResponse(405, nil, errors.New("group have more user"))
		return
	}

	// 查询组是否有关联物料
	materials, err := material.MaterialByGroupId(groupId, 1, 10)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(materials) > 0 {
		c.WriteHttpResponse(405, nil, errors.New("group have more material"))
		return
	}

	if err := account.DelGroup(groupId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *GroupController) GroupById() {
	groupIdStr := c.Ctx.Input.Param(":id")
	log.Debug(groupIdStr)
	if len(groupIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("group id is empty"))
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := account.GroupById(groupId)
	if err != nil {
		if !account.ErrGroupNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有
func (c *GroupController) GroupList() {
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

	list, err := account.GroupList(page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
