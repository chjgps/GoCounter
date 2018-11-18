package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/material"
)

type CategoryController struct {
	BaseController
}

// 添加
func (c *CategoryController) AddCategory() {
	obj := &material.Category{}

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
	acc, err := material.CategoryByName(name)
	if err != nil {
		if !material.ErrCategoryNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(material.ErrCategoryAlreadyExist))
		return
	}

	obj.Created = timex.String()
	if err := material.InsertCategory(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *CategoryController) EditCategory() {
	obj := &material.Category{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	categoryId := obj.Id

	// 查询用户是否存在
	acc, err := material.CategoryById(categoryId)
	if err != nil {
		if !material.ErrCategoryNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(material.ErrCategoryNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := material.UpdateCategory(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *CategoryController) DelCategory() {
	categoryIdStr := c.Ctx.Input.Param(":id")
	log.Debug(categoryIdStr)
	if len(categoryIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("category id is empty"))
		return
	}

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := material.DelCategory(categoryId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 查询所有
func (c *CategoryController) CategoryList() {
	list, err := material.CategoryList()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
