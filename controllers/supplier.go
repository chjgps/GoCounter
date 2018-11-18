package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/material"
)

type SupplierController struct {
	BaseController
}

// 添加
func (c *SupplierController) AddSupplier() {
	obj := &material.Supplier{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	shortName := obj.ShortName

	// 查询用户名是否重复
	acc, err := material.SupplierByName(shortName)
	if err != nil {
		if !material.ErrSupplierNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(material.ErrSupplierAlreadyExist))
		return
	}

	obj.Status = 1
	obj.Created = timex.String()
	if err := material.InsertSupplier(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, struct {
		SupplierId int `json:"supplierIdId"`
	}{
		SupplierId: obj.Id,
	}, nil)
	return
}

// 修改
func (c *SupplierController) EditSupplier() {
	obj := &material.Supplier{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	supplierId := obj.Id

	// 查询是否存在
	acc, err := material.SupplierById(supplierId)
	if err != nil {
		if !material.ErrSupplierNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(material.ErrSupplierNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := material.UpdateSupplier(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *SupplierController) DelSupplier() {
	supplierIdStr := c.Ctx.Input.Param(":id")
	log.Debug(supplierIdStr)
	if len(supplierIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("supplier id is empty"))
		return
	}

	supplierId, err := strconv.Atoi(supplierIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := material.DelSupplier(supplierId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *SupplierController) SupplierById() {
	supplierIdStr := c.Ctx.Input.Param(":id")
	log.Debug(supplierIdStr)
	if len(supplierIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("supplier id is empty"))
		return
	}

	supplierId, err := strconv.Atoi(supplierIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := material.SupplierById(supplierId)
	if err != nil {
		if !material.ErrSupplierNotFound.Equal(err) {
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
func (c *SupplierController) SupplierList() {
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

	total, list, err := material.SupplierList(map[string]interface{}{
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

// 查询所有物料
func (c *SupplierController) MaterialList() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	name := c.GetString("name")

	categoryIdStr := c.Input().Get("categoryId")

	var categoryId int
	var err error
	if len(categoryIdStr) > 0 {
		categoryId, err = strconv.Atoi(categoryIdStr)
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

	total, list, err := material.SupplierMaterialList(map[string]interface{}{
		"startDate":  startDate,
		"endDate":    endDate,
		"name":       name,
		"categoryId": categoryId,
	}, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	var data interface{}
	if list == nil {
		data = make([]interface{}, 0)
	} else {
		for _, v := range list {
			v.Img = fmt.Sprintf("%s/%s%s", ImgHost, ImgMaterialDir, v.Img)
		}

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
