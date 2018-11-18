package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/permission"
)

type UserController struct {
	BaseController
}

// 添加用户
func (c *UserController) AddUser() {
	obj := &permission.User{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	username := obj.Username

	// 查询用户名是否重复
	acc, err := permission.UserByName(username)
	if err != nil {
		if !permission.ErrUserNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(permission.ErrUsernameAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := permission.InsertUser(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改用户
func (c *UserController) EditUser() {
	obj := &permission.User{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	userId := obj.Id

	// 查询用户是否存在
	acc, err := permission.UserById(userId)
	if err != nil {
		if !permission.ErrUserNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(permission.ErrUserNotFound))
		return
	}

	// 密码为空，不修改
	if len(obj.Password) == 0 {
		obj.Password = acc.Password
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := permission.UpdateUser(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改密码
func (c *UserController) EditPassword() {
	obj := &permission.User{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	userId := obj.Id

	// 查询用户是否存在
	acc, err := permission.UserById(userId)
	if err != nil {
		if !permission.ErrUserNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(permission.ErrUserNotFound))
		return
	}

	// 密码为空
	password := obj.Password
	if len(password) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("password is empty"))
		return
	}

	acc.Password = password
	if err := permission.UpdateUser(acc); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除用户
func (c *UserController) DelUser() {
	userIdStr := c.Ctx.Input.Param(":id")
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

	if err := permission.DelUser(userId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询用户
func (c *UserController) UserById() {
	userIdStr := c.Ctx.Input.Param(":id")
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

	acc, err := permission.UserById(userId)
	if err != nil {
		if !permission.ErrUserNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 根据用户名登录
func (c *UserController) Login() {
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

	acc, err := permission.LoginByUsername(username, password)
	if err != nil {
		if !permission.ErrUserNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
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
func (c *UserController) UserList() {
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

	total, list, err := permission.UserList(map[string]interface{}{
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
