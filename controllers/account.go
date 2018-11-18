package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/account"
)

type AccountController struct {
	BaseController
}

// 添加用户
func (c *AccountController) AddAccount() {
	obj := &account.Account{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	username := obj.Username
	card := obj.Card

	// 查询用户名是否重复
	acc, err := account.AccountByName(username)
	if err != nil {
		if !account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(account.ErrUsernameAlreadyExist))
		return
	}

	// 查询卡号是否重复
	acc, err = account.AccountByCard(card)
	if err != nil {
		if !account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(account.ErrCardAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := account.InsertAccount(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改用户
func (c *AccountController) EditAccount() {
	obj := &account.Account{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	accountId := obj.Id

	// 查询用户是否存在
	acc, err := account.AccountById(accountId)
	if err != nil {
		if !account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc == nil {
		c.WriteHttpResponse(404, nil, errors.As(account.ErrAccountNotFound))
		return
	}

	// 密码为空，不修改
	if len(obj.Password) == 0 {
		obj.Password = acc.Password
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := account.UpdateAccount(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改密码
func (c *AccountController) EditPassword() {
	obj := &account.Account{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	accountId := obj.Id

	// 查询用户是否存在
	acc, err := account.AccountById(accountId)
	if err != nil {
		if !account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc == nil {
		c.WriteHttpResponse(404, nil, errors.As(account.ErrAccountNotFound))
		return
	}

	// 密码为空
	password := obj.Password
	if len(password) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("password is empty"))
		return
	}

	acc.Password = password
	if err := account.UpdateAccount(acc); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除用户
func (c *AccountController) DelAccount() {
	accountIdStr := c.Ctx.Input.Param(":id")
	log.Debug(accountIdStr)
	if len(accountIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("account id is empty"))
		return
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if err := account.DelAccount(accountId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询用户
func (c *AccountController) AccountById() {
	accountIdStr := c.Ctx.Input.Param(":id")
	log.Debug(accountIdStr)
	if len(accountIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("account id is empty"))
		return
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	acc, err := account.AccountById(accountId)
	if err != nil {
		if account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 根据卡查询用户
func (c *AccountController) AccountByCard() {
	cardId := c.Ctx.Input.Param(":cardId")
	log.Debug(cardId)
	if len(cardId) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("card id is empty"))
		return
	}

	acc, err := account.AccountByCard(cardId)
	if err != nil {
		if account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 根据卡登录
func (c *AccountController) LoginByCard() {
	card := c.GetString("card")
	log.Debug(card)
	if len(card) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("card is empty"))
		return
	}

	password := c.GetString("password")
	log.Debug(password)
	if len(password) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("password is empty"))
		return
	}

	acc, err := account.LoginByCard(card, password)
	if err != nil {
		if account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	acc.Token = OAuth.Add()

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 根据用户名登录
func (c *AccountController) LoginByUsername() {
	username := c.GetString("username")
	log.Debug(username)
	if len(username) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("username is empty"))
		return
	}

	password := c.GetString("password")
	log.Debug(password)
	if len(password) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("password is empty"))
		return
	}

	acc, err := account.LoginByUsername(username, password)
	if err != nil {
		if account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	acc.Token = OAuth.Add()

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询所有用户
func (c *AccountController) AccountList() {
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

	total, list, err := account.AccountList(map[string]interface{}{
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

// 根据指纹ID查询用户
func (c *AccountController) AccountByFinger() {
	fingerIdStr := c.Ctx.Input.Param(":id")
	log.Debug(fingerIdStr)
	if len(fingerIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("finger id is empty"))
		return
	}

	fingerId, err := strconv.Atoi(fingerIdStr)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	acc, err := account.AccountByFinger(fingerId)
	if err != nil {
		if account.ErrAccountNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	acc.Token = OAuth.Add()

	c.WriteHttpResponse(200, acc, nil)
	return
}
