package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/order"
)

type OrderController struct {
	BaseController
}

// 添加
func (c *OrderController) AddOrder() {
	obj := &order.Order{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
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

	obj.Status = 1
	obj.Created = timex.String()
	if err := order.InsertOrder(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *OrderController) EditOrder() {
	obj := &order.Order{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	orderId := obj.Id

	// 查询用户是否存在
	acc, err := order.OrderById(orderId)
	if err != nil {
		if !order.ErrOrderNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(order.ErrOrderNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := order.UpdateOrder(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *OrderController) DelOrder() {
	orderIdStr := c.Ctx.Input.Param(":id")
	log.Debug(orderIdStr)
	if len(orderIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("order id is empty"))
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := order.DelOrder(orderId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *OrderController) OrderById() {
	orderIdStr := c.Ctx.Input.Param(":id")
	log.Debug(orderIdStr)
	if len(orderIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("order id is empty"))
		return
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := order.OrderById(orderId)
	if err != nil {
		if !order.ErrOrderNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, acc, nil)
	return
}

// 查询待归还数据
func (c *OrderController) OrderByAccountId() {
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

	acc, err := order.OrderByAccountId(accountId)
	if err != nil {
		if !order.ErrOrderNotFound.Equal(err) {
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
func (c *OrderController) OrderList() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	name := c.GetString("name")

	accountIdStr := c.Input().Get("accountId")
	typeIdStr := c.Input().Get("type")

	var err error

	var accountId int
	if len(accountIdStr) > 0 {
		accountId, err = strconv.Atoi(accountIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	var typeId int
	if len(typeIdStr) > 0 {
		typeId, err = strconv.Atoi(typeIdStr)
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

	total, list, err := order.OrderList(map[string]interface{}{
		"startDate": startDate,
		"endDate":   endDate,
		"name":      name,
		"accountId": accountId,
		"type":      typeId,
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

// 查询所有待回收物料
func (c *OrderController) RecycleList() {
	accountIdStr := c.Input().Get("accountId")

	if len(accountIdStr) < 0 {
		c.WriteHttpResponse(400, nil, errors.New("accountId is empty"))
		return
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	list, err := order.RecycleList(accountId)
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

	c.WriteHttpResponse(200, data, nil)
	return
}
