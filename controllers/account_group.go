package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/account"
)

type AccountGroupController struct {
	BaseController
}

// 添加用户组绑定关系
func (c *AccountGroupController) AddAccountGroup() {
	obj := &account.AccountGroup{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	groupId := obj.GroupId
	accountId := obj.AccountId

	// 查询是否存在
	acc, err := account.AccountGroupById(groupId, accountId)
	if err != nil {
		if !account.ErrAccountGroupNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	// 不存在添加
	if acc == nil {
		obj.Status = 1
		obj.Created = timex.String()
		if err := account.InsertAccountGroup(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	} else {
		// 存在修改
		acc.Status = 1
		acc.Updated = timex.String()
		if err := account.UpdateAccountGroup(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除用户组
func (c *AccountGroupController) DelAccountGroup() {
	groupIdStr := c.Ctx.Input.Param(":groupId")
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

	accountIdStr := c.Ctx.Input.Param(":accountId")
	log.Debug(accountIdStr)
	if len(accountIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("account id is empty"))
		return
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// 查询是否存在
	acc, err := account.AccountGroupById(groupId, accountId)
	if err != nil {
		if account.ErrAccountGroupNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	// 存在
	if err := account.DelAccountGroup(acc.Id); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据组ID查询所有用户
func (c *AccountGroupController) AccountByGroupId() {
	groupId, err := c.GetInt("groupId")
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

	list, err := account.AccountByGroupId(groupId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
