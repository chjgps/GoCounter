package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/order"
	"github.com/robfig/cron"
)

type AutoConfController struct {
	BaseController
}

var (
	cr = cron.New()
	// 秒 分 时 日 月 周
	autoConfSpec = "0 %d %d * * *"
)

// socket消息
func init() {
	startAutoConf()
}

func startAutoConf() {
	// 查询所有auto配置
	list, err := order.AutoConfList()
	if err != nil {
		log.Error("%v", errors.As(err))
		return
	}

	for _, v := range list {
		if err := cr.AddFunc(fmt.Sprintf(autoConfSpec, v.Minute, v.Hour), runAutoConf); err != nil {
			log.Error("%v", errors.As(err))
			return
		}
	}

	cr.Start()
}

func reloadAutoConf() {
	cr.Stop()

	// 查询所有auto配置
	list, err := order.AutoConfList()
	if err != nil {
		log.Error("%v", errors.As(err))
		return
	}

	for _, v := range list {
		if err := cr.AddFunc(fmt.Sprintf(autoConfSpec, v.Minute, v.Hour), runAutoConf); err != nil {
			log.Error("%v", errors.As(err))
			return
		}
	}

	cr.Start()
}

func runAutoConf() {
	// socket
	Server.BroadcastTo("login", "auto", "success")
}

// 添加
func (c *AutoConfController) AddAutoConf() {
	obj := &order.AutoConf{}

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
	acc, err := order.AutoConfByName(name)
	if err != nil {
		if !order.ErrAutoConfNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	if acc != nil {
		c.WriteHttpResponse(400, nil, errors.As(order.ErrAutoConfNameAlreadyExist))
		return
	}

	obj.Created = timex.String()
	if err := order.InsertAutoConf(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	reloadAutoConf()

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 修改
func (c *AutoConfController) EditAutoConf() {
	obj := &order.AutoConf{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	autoId := obj.Id

	// 查询是否存在
	acc, err := order.AutoConfById(autoId)
	if err != nil {
		if !order.ErrAutoConfNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(order.ErrAutoConfNotFound))
		return
	}

	obj.Created = acc.Created
	obj.Updated = timex.String()
	if err := order.UpdateAutoConf(obj); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	reloadAutoConf()

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除
func (c *AutoConfController) DelAutoConf() {
	autoIdStr := c.Ctx.Input.Param(":id")
	log.Debug(autoIdStr)
	if len(autoIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("auto id is empty"))
		return
	}

	autoId, err := strconv.Atoi(autoIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if err := order.DelAutoConf(autoId); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	reloadAutoConf()

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据ID查询
func (c *AutoConfController) AutoConfById() {
	autoIdStr := c.Ctx.Input.Param(":id")
	log.Debug(autoIdStr)
	if len(autoIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("auto id is empty"))
		return
	}

	autoId, err := strconv.Atoi(autoIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	acc, err := order.AutoConfById(autoId)
	if err != nil {
		if !order.ErrAutoConfNotFound.Equal(err) {
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
func (c *AutoConfController) AutoConfList() {
	list, err := order.AutoConfList()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
