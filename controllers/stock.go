package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	// "github.com/beego/ms304w-client/basis/bytex"
	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/httpx/rest"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/box"
	"github.com/beego/ms304w-client/models/material"
	"github.com/beego/ms304w-client/models/order"
)

type StockController struct {
	BaseController
}

// 上料
func (c *StockController) StockIn() {
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

	materialId := obj.MaterialId
	if materialId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("materialId is illegal"))
		return
	}

	qty := obj.Qty
	if qty <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("qty is illegal"))
		return
	}

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

	// 查询物料绑定的格子
	gridList, err := box.GridByMaterialId(materialId)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	// 物料未绑定格子
	if len(gridList) == 0 {
		c.WriteHttpResponse(404, nil, errors.As(box.ErrGridNotFound))
		return
	}

	// 如果有指定格子，判断是否已超过最大容量
	var boxAddr int
	var gridId int
	var gridChannel int

	for _, v := range gridList {
		maxQty := v.Qty
		totalQty := v.TotalQty

		// 判断是否大于最大库存
		if totalQty >= maxQty && maxQty != 0 {
			// 格子已满
			continue
		} else {
			// 指定格子
			boxAddr = v.Addr
			gridId = v.Id
			gridChannel = v.Channel
			break
		}
	}

	// 没有分配格子
	if gridId <= 0 {
		c.WriteHttpResponse(404, nil, errors.As(box.ErrGridNotFound))
		return
	}

	log.Info("boxAddr %d, gridId %d, gridChannel %d", boxAddr, gridId, gridChannel)

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

	log.Info("Post %s", SERIAL_OPEN_URL)

	bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	log.Info("StockIn %s", string(bytes))

	c.WriteHttpResponse(200, struct {
		Channel int `json:"channel"`
	}{
		Channel: channel,
	}, nil)

	return
}

// 上料确认
func (c *StockController) StockInConfirm() {
	log.Info("StockInConfirm: %s", string(c.Ctx.Input.RequestBody))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 领料
func (c *StockController) StockOut() {
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

	materialId := obj.MaterialId
	if materialId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("materialId is illegal"))
		return
	}

	qty := obj.Qty
	if qty <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("qty is illegal"))
		return
	}

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

	// 查询物料绑定的格子
	gridList, err := box.GridByMaterialId(materialId)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	// 物料未绑定格子
	if len(gridList) == 0 {
		c.WriteHttpResponse(404, nil, errors.As(box.ErrGridNotFound))
		return
	}

	// 如果有指定格子，判断是否已超过最大容量
	boxAddr := gridList[0].Addr
	gridId := gridList[0].Id
	gridChannel := gridList[0].Channel

	log.Info("boxAddr %d, gridId %d, gridChannel %d", boxAddr, gridId, gridChannel)

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

		// 开门
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

	log.Info("StockOut %s", string(bytes))

	c.WriteHttpResponse(200, struct {
		Channel int `json:"channel"`
	}{
		Channel: channel,
	}, nil)
	return
}

// 领料确认
func (c *StockController) StockOutConfirm() {
	log.Info("StockOutConfirm: %s", string(c.Ctx.Input.RequestBody))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 回收
func (c *StockController) StockRecycle() {
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

	materialId := obj.MaterialId
	if materialId <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("materialId is illegal"))
		return
	}

	qty := obj.Qty
	if qty <= 0 {
		c.WriteHttpResponse(400, nil, errors.New("qty is illegal"))
		return
	}

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

	// 查询物料绑定的格子
	gridList, err := box.GridByMaterialId(materialId)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	// 物料未绑定格子
	if len(gridList) == 0 {
		c.WriteHttpResponse(404, nil, errors.As(box.ErrGridNotFound))
		return
	}

	// 如果有指定格子，判断是否已超过最大容量
	var boxAddr int
	var gridId int
	var gridChannel int

	for _, v := range gridList {
		maxQty := v.Qty
		totalQty := v.TotalQty

		// 判断是否大于最大库存
		if totalQty >= maxQty && maxQty != 0 {
			// 格子已满
			continue
		} else {
			// 指定格子
			boxAddr = v.Addr
			gridId = v.Id
			gridChannel = v.Channel
			break
		}
	}

	// 没有分配格子
	if gridId <= 0 {
		c.WriteHttpResponse(404, nil, errors.As(box.ErrGridNotFound))
		return
	}

	log.Info("boxAddr %d, gridId %d, gridChannel %d", boxAddr, gridId, gridChannel)

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
		Channel int `json:"channel"`
	}{
		Channel: channel,
	}, nil)
	return
}

// 回收确认
func (c *StockController) StockRecycleConfirm() {
	log.Info("StockRecycleConfirm: %s", string(c.Ctx.Input.RequestBody))

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 格子库存
func (c *StockController) StockList() {
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

	// gridId
	gridIdStr := c.Input().Get("gridId")

	var gridId int
	if len(gridIdStr) > 0 {
		gridId, err = strconv.Atoi(gridIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// categoryId
	categoryIdStr := c.Input().Get("categoryId")

	var categoryId int
	if len(categoryIdStr) > 0 {
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// materialId
	materialIdStr := c.Input().Get("materialId")

	var materialId int
	if len(materialIdStr) > 0 {
		materialId, err = strconv.Atoi(materialIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	total, list, err := order.StockList(map[string]interface{}{
		"startDate":  startDate,
		"endDate":    endDate,
		"name":       name,
		"gridId":     gridId,
		"categoryId": categoryId,
		"materialId": materialId,
	}, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
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

// 物料库存
func (c *StockController) MaterialStockList() {
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

	// categoryId
	categoryIdStr := c.Input().Get("categoryId")
	var categoryId int

	if len(categoryIdStr) > 0 {
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// materialId
	materialIdStr := c.Input().Get("materialId")
	var materialId int

	if len(materialIdStr) > 0 {
		materialId, err = strconv.Atoi(materialIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	total, list, err := order.MaterialStockList(map[string]interface{}{
		"startDate":  startDate,
		"endDate":    endDate,
		"name":       name,
		"categoryId": categoryId,
		"materialId": materialId,
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

// 组物料库存
func (c *StockController) GroupStockList() {
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	name := c.GetString("name")

	// accountId
	accountId, err := c.GetInt("accountId")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// page
	page, err := c.GetInt("page")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// pageSize
	pageSize, err := c.GetInt("pageSize")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// categoryId
	categoryIdStr := c.Input().Get("categoryId")
	var categoryId int

	if len(categoryIdStr) > 0 {
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// materialId
	materialIdStr := c.Input().Get("materialId")
	var materialId int

	if len(materialIdStr) > 0 {
		materialId, err = strconv.Atoi(materialIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	total, list, err := order.GroupStock(map[string]interface{}{
		"startDate":  startDate,
		"endDate":    endDate,
		"name":       name,
		"categoryId": categoryId,
		"materialId": materialId,
	}, accountId, page, pageSize)
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

// 根据物料ID查询格子
func (c *StockController) GridByMaterialId() {
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

	acc, err := order.GridByMaterialId(orderId)
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

// 自动盘点
func (c *StockController) Auto() {
	// accountId
	accountIdStr := c.GetString("accountId")
	if len(accountIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("account id is empty"))
		return
	}

	var accountId int
	var err error
	if len(accountIdStr) > 0 {
		accountId, err = strconv.Atoi(accountIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	log.Info("%#v", accountId)

	// boxId
	boxIdStr := c.GetString("boxId")
	if len(boxIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("box id is empty"))
		return
	}

	var boxId int
	if len(boxIdStr) > 0 {
		boxId, err = strconv.Atoi(boxIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// 查询柜子是否存在
	boxObj, err := box.BoxById(boxId)
	if err != nil {
		if box.ErrBoxNotFound.Equal(err) {
			c.WriteHttpResponse(404, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	boxAddr := boxObj.Addr
	/*
		boxBytes, err := bytex.IntToBytes(int32(boxAddr))
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	*/

	// 查询所有格子
	_, list, err := box.GridList(map[string]interface{}{
		"startDate": "",
		"endDate":   "",
		"name":      "",
		"boxId":     boxId,
		"sensorId":  0,
	}, 1, 1000)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	for _, v := range list {
		gridId := v.Id

		// 查询库存
		stock, err := order.StockByGridId(gridId)
		if err != nil {
			if order.ErrStockNotFound.Equal(err) {
				continue
			}

			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		auto := &order.Auto{
			Created:    timex.String(),
			AccountId:  accountId,
			GridId:     gridId,
			SensorId:   stock.SensorId,
			MaterialId: stock.MaterialId,
			BeforeQty:  stock.Qty,
		}

		// 添加盘点数据
		if err := order.InsertAuto(auto); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		// 查询格子配置
		_, list, err := box.ChannelList(map[string]interface{}{
			"startDate": "",
			"endDate":   "",
			"gridId":    gridId,
			"sensorId":  0,
			"name":      "",
		}, 1, 1000)
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		if len(list) != 1 {
			c.WriteHttpResponse(500, nil, errors.New("material and sensor not connect"))
			return
		}

		// 物料传感器
		channel := list[0].Channel

		/*
			// 盘点
			if err := SerialWeight(strconv.Itoa(auto.Id), boxBytes, byte(channel), LOCK_CHECK); err != nil {
				c.WriteHttpResponse(500, nil, errors.As(err))
				return
			}
		*/

		where := make(map[string]interface{})
		where["uuid"] = strconv.Itoa(auto.Id)
		where["boxId"] = boxAddr
		where["gridId"] = channel

		bytes, err := rest.Post(SERIAL_OPEN_URL).PostQuerys(where).End()
		if err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		log.Info("StockIn %s", string(bytes))

	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 自动盘点查询
func (c *StockController) AutoList() {
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

	// boxId
	boxIdStr := c.Input().Get("boxId")
	var boxId int

	if len(boxIdStr) > 0 {
		boxId, err = strconv.Atoi(boxIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	// accountId
	accountIdStr := c.Input().Get("accountId")
	var accountId int

	if len(accountIdStr) > 0 {
		accountId, err = strconv.Atoi(accountIdStr)
		if err != nil {
			c.WriteHttpResponse(400, nil, errors.As(err))
			return
		}
	}

	total, list, err := order.AutoList(map[string]interface{}{
		"startDate": startDate,
		"endDate":   endDate,
		"name":      name,
		"boxId":     boxId,
		"accountId": accountId,
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
