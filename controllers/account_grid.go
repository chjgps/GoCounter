package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/box"
)

type AccountGridController struct {
	BaseController
}

// 添加
func (c *AccountGridController) AddAccount() {
	obj := &box.Account{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	if obj.AccountId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("account is empty"))
		return
	}

	if obj.BoxId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("box is empty"))
		return
	}

	if obj.GridId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("grid is empty"))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := box.InsertAccount(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *AccountGridController) EditAccount() {
	obj := &box.Account{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	id := obj.Id

	// 查询用户是否存在
	acc, err := box.AccountById(id)
	if err != nil {
		if box.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(box.ErrAccountNotFound))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := box.UpdateAccount(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *AccountGridController) DelAccount() {
	idStr := c.Ctx.Input.Param(":id")
	log.Debug(idStr)
	if len(idStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("id is empty"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := box.DelAccount(id); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *AccountGridController) AccountById() {
	idStr := c.Ctx.Input.Param(":id")
	log.Debug(idStr)
	if len(idStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("id is empty"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := box.AccountById(id)
	if err != nil {
		if box.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有
func (c *AccountGridController) AccountList() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	name := c.GetString("name")

	accountIdStr := c.Input().Get("accountId")
	boxIdStr := c.Input().Get("boxId")

	var accountId int
	var boxId int
	var err error

	if len(accountIdStr) > 0 {
		accountId, err = strconv.Atoi(accountIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	if len(boxIdStr) > 0 {
		boxId, err = strconv.Atoi(boxIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
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

	total, list, err := box.AccountList(map[string]interface{}{
		"startDate": startDate,
		"endDate":   endDate,
		"name":      name,
		"accountId": accountId,
		"boxId":     boxId,
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
