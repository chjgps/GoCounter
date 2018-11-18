package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/httpx/rest"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/box"
	"github.com/beego/ms304w-client/models/material"
	"github.com/beego/ms304w-client/models/order"
)

type StockCodeController struct {
	BaseController
}

// 上料
func (c *StockCodeController) StockIn() {
	obj := &order.Request{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	accountId := obj.AccountId
	if accountId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("accountId is illegal"))
		return
	}

	code := obj.Code
	if len(code) <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("code is illegal"))
		return
	}

	qty := obj.Qty

	// 查询code对应的格子和物料
	g, err := box.GridByCode(code)
	if err != nil {
		if box.ErrGridNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	materialId := g.MaterialId
	boxAddr := g.Addr
	gridId := g.Id
	gridChannel := g.Channel

	log.Info("%#v", gridChannel)

	// 根据物料查询传感器
	_, sensor, err := material.SensorList(map[string]interface{}{
		"startDate":  "",
		"endDate":    "",
		"name":       "",
		"materialId": materialId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	// 每个物料只有一个传感器
	if len(sensor) != 1 {
		c.WriteHttpResponse(500, nil, errors.New("material and sensor not connect"))
		return
	}

	// 物料传感器
	sensorId := sensor[0].SensorId

	// 查询通道
	_, ch, err := box.ChannelList(map[string]interface{}{
		"startDate": "",
		"endDate":   "",
		"name":      "",
		"gridId":    gridId,
		"sensorId":  sensorId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(ch) == 0 {
		c.WriteHttpResponse(500, nil, errors.New("gridId and sensorId not connect"))
		return
	}

	channel := ch[0].Channel

	// 添加订单
	o := &order.Order{
		Created:    timex.String(),
		AccountId:  accountId,
		Type:       order.IN,
		GridId:     gridId,
		MaterialId: materialId,
		SensorId:   sensorId,
		Channel:    channel,
		Qty:        qty,
		Status:     1,
	}

	if err := order.InsertOrder(o); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	/*
		boxBytes, err := bytex.IntToBytes(int32(boxAddr))
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		// 打开柜门
		if err := SerialOpen(strconv.Itoa(o.Id), boxBytes, byte(gridChannel), LOCK_SUM); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	*/

	where := make(map[string]interface{})
	where["uuid"] = strconv.Itoa(o.Id)
	where["boxId"] = boxAddr
	where["gridId"] = gridChannel
	where["operation"] = 2

	bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("StockIn %s", string(bytes))

	c.WriteHttpResponse(200, struct {
		MaterialId   int    `json:"materialId"`
		MaterialCode string `json:"materialCode"`
		Qty          int    `json:"qty"`
	}{
		MaterialId:   materialId,
		MaterialCode: g.MaterialCode,
		Qty:          g.TotalQty,
	}, nil)

	return
}

// 领料
func (c *StockCodeController) StockOut() {
	obj := &order.Request{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	accountId := obj.AccountId
	if accountId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("accountId is illegal"))
		return
	}

	code := obj.Code
	if len(code) <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("code is illegal"))
		return
	}

	qty := obj.Qty

	// 查询code对应的格子和物料
	g, err := box.GridByCode(code)
	if err != nil {
		if box.ErrGridNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	materialId := g.MaterialId
	boxAddr := g.Addr
	gridId := g.Id
	gridChannel := g.Channel

	log.Info("%#v", gridChannel)

	// 根据物料查询传感器
	_, sensor, err := material.SensorList(map[string]interface{}{
		"startDate":  "",
		"endDate":    "",
		"name":       "",
		"materialId": materialId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(sensor) != 1 {
		c.WriteHttpResponse(500, nil, errors.New("material and sensor not connect"))
		return
	}

	sensorId := sensor[0].SensorId

	// 查询通道
	_, ch, err := box.ChannelList(map[string]interface{}{
		"startDate": "",
		"endDate":   "",
		"name":      "",
		"gridId":    gridId,
		"sensorId":  sensorId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(ch) == 0 {
		c.WriteHttpResponse(500, nil, errors.New("gridId and sensorId not connect"))
		return
	}

	channel := ch[0].Channel

	// 添加订单
	o := &order.Order{
		Created:    timex.String(),
		AccountId:  accountId,
		Type:       order.OUT,
		GridId:     gridId,
		MaterialId: materialId,
		SensorId:   sensorId,
		Channel:    channel,
		Qty:        qty,
		Status:     1,
	}

	if err := order.InsertOrder(o); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	/*
		boxBytes, err := bytex.IntToBytes(int32(boxAddr))
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		// 打开柜门
		if err := SerialOpen(strconv.Itoa(o.Id), boxBytes, byte(gridChannel), LOCK_SUM); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	*/

	where := make(map[string]interface{})
	where["uuid"] = strconv.Itoa(o.Id)
	where["boxId"] = boxAddr
	where["gridId"] = gridChannel
	where["operation"] = 2

	bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("StockIn %s", string(bytes))

	c.WriteHttpResponse(200, struct {
		MaterialId   int    `json:"materialId"`
		MaterialCode string `json:"materialCode"`
		Qty          int    `json:"qty"`
		MaterialName string `json:"materialName"`
	}{
		MaterialId:   materialId,
		MaterialCode: g.MaterialCode,
		Qty:          g.TotalQty,
		MaterialName: g.MaterialName,
	}, nil)

	return
}

// 回收
func (c *StockCodeController) StockRecycle() {
	obj := &order.Request{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	accountId := obj.AccountId
	if accountId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("accountId is illegal"))
		return
	}

	code := obj.Code
	if len(code) <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("code is illegal"))
		return
	}

	qty := obj.Qty

	// 查询code对应的格子和物料
	g, err := box.GridByCode(code)
	if err != nil {
		if box.ErrGridNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	materialId := g.MaterialId
	boxAddr := g.Addr
	gridId := g.Id
	gridChannel := g.Channel

	log.Info("%#v", gridChannel)

	// 根据物料查询传感器
	_, sensor, err := material.SensorList(map[string]interface{}{
		"startDate":  "",
		"endDate":    "",
		"name":       "",
		"materialId": materialId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(sensor) != 1 {
		c.WriteHttpResponse(500, nil, errors.New("material and sensor not connect"))
		return
	}

	sensorId := sensor[0].SensorId

	// 查询通道
	_, ch, err := box.ChannelList(map[string]interface{}{
		"startDate": "",
		"endDate":   "",
		"name":      "",
		"gridId":    gridId,
		"sensorId":  sensorId,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	if len(ch) == 0 {
		c.WriteHttpResponse(500, nil, errors.New("gridId and sensorId not connect"))
		return
	}

	channel := ch[0].Channel

	// 添加订单
	o := &order.Order{
		Created:    timex.String(),
		AccountId:  accountId,
		Type:       order.RECYCLE,
		GridId:     gridId,
		MaterialId: materialId,
		SensorId:   sensorId,
		Channel:    channel,
		Qty:        qty,
		Status:     1,
	}

	if err := order.InsertOrder(o); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	/*
		boxBytes, err := bytex.IntToBytes(int32(boxAddr))
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		// 打开柜门
		if err := SerialOpen(strconv.Itoa(o.Id), boxBytes, byte(gridChannel), LOCK_SUM); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	*/

	where := make(map[string]interface{})
	where["uuid"] = strconv.Itoa(o.Id)
	where["boxId"] = boxAddr
	where["gridId"] = gridChannel
	where["operation"] = 2

	bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("StockRecycle %s", string(bytes))

	c.WriteHttpResponse(200, struct {
		MaterialId   int    `json:"materialId"`
		MaterialCode string `json:"materialCode"`
		Qty          int    `json:"qty"`
	}{
		MaterialId:   materialId,
		MaterialCode: g.MaterialCode,
		Qty:          g.TotalQty,
	}, nil)

	return
}
