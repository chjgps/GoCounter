package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/material"
	"github.com/beego/ms304w-client/models/order"
)

type CallbackController struct {
	BaseController
}

type SerialRequest struct {
	Api  string  `json:"api"`
	Code int     `json:"code"`
	Err  string  `json:"err"`
	Data *CbData `json:"data"`
}

type CbData struct {
	ApiData
}

func (r *CbData) String() string {
	bytes, err := json.MarshalIndent(r, "  ", "    ")
	if err != nil {
		return fmt.Sprintf("%#v", err)
	}

	return string(bytes)
}

// 关门传感器结果返回
type ResData struct {
	MaterialId int `json:"materialId"`
	Type       int `json:"type"`
	Qty        int `json:"qty"`
}

func (r *ResData) String() string {
	bytes, err := json.MarshalIndent(r, "  ", "    ")
	if err != nil {
		return fmt.Sprintf("%#v", err)
	}

	return string(bytes)
}

// -------------------------
// 刷卡回调

func (c *CallbackController) Card() {
	log.Info("Card: %s", string(c.Ctx.Input.RequestBody))

	obj := &SerialRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	Server.BroadcastTo("login", "loginByCard", obj.Data)

	c.WriteHttpResponse(200, nil, nil)
	return
}

// -------------------------
// 关门结果回调

func (c *CallbackController) Weight() {
	log.Info("Weight: %s", string(c.Ctx.Input.RequestBody))

	obj := &SerialRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	cbData := obj.Data
	if cbData == nil {
		c.WriteHttpResponse(400, nil, errors.New("callback data is empty"))
		return
	}

	Server.BroadcastTo("login", "weight", cbData.Weight)

	if cbData.Operation != LOCK_WEIGHT {
		i := WeightCallback(cbData.UUID, cbData.Weight)
		log.Info("WeightCallback %d", i)
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

func WeightCallback(orderIdStr string, gridWeight int) int {
	if len(orderIdStr) == 0 {
		log.Error("%v", errors.New("uuid is empty"))
		return 0
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// 查询订单
	o, err := order.OrderById(orderId)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// 结果返回
	resData := &ResData{
		MaterialId: o.MaterialId,
		Type:       o.Type,
	}

	var qty int
	// 查询物料传感器配置参数
	materialSensor, err := material.SensorByMaterialId(o.MaterialId, o.SensorId)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	params, err := materialSensor.ParamsObj()
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	if gridWeight > 0 {
		minQty := params.Weight - params.Lower
		maxQty := params.Weight + params.ComeUp
		log.Info("WEIGHT value %d, params weight %d, min %d, max %d", gridWeight, params.Weight, minQty, maxQty)

		// 取余数量
		num := gridWeight / params.Weight

		// 剩余数量
		subQty := gridWeight - num*params.Weight

		log.Info("weight %d", num)

		if minQty <= subQty && subQty <= maxQty {
			qty = num + 1
		} else {
			qty = num
		}

		log.Info("weight %d", qty)
	}

	log.Info("----------QTY---------- %d", qty)

	// 查询库存,如果新上料没有库存,领料和回收有库存
	stockObj, err := order.StockByMaterialId(o.MaterialId, o.GridId)
	if err != nil {
		if !order.ErrStockNotFound.Equal(err) {
			log.Error("%v", errors.As(err))
			return 0
		}

		// 更新订单
		o.Updated = timex.String()
		o.BeforeQty = 0
		o.Qty = qty
		o.AfterQty = qty
		if err := order.UpdateOrder(o); err != nil {
			log.Error("%v", errors.As(err))
			return 0
		}

		// insert
		if err := order.InsertStock(&order.Stock{
			Created:    timex.String(),
			GridId:     o.GridId,
			SensorId:   o.SensorId,
			MaterialId: o.MaterialId,
			Qty:        qty,
		}); err != nil {
			log.Error("%v", errors.As(err))
			return 0
		}

		// res
		resData.Qty = qty

		log.Warn("result %s", resData.String())
		Server.BroadcastTo("login", "inventory", resData)
		return 1
	}

	var updateQty int
	switch o.Type {
	case order.IN:
		// 上料
		updateQty = qty - stockObj.Qty

	case order.OUT:
		// 领料
		updateQty = stockObj.Qty - qty

		// 删除空格子
		// TODO:格子已为空
		if qty == 0 {
			if err := order.DelStockByGridId(o.GridId); err != nil {
				log.Error("%v", errors.As(err))
				return 0
			}
		}

	case order.RECYCLE:
		// 回收
		updateQty = qty - stockObj.Qty
	}

	// update
	// 更新订单
	o.Updated = timex.String()
	o.BeforeQty = stockObj.Qty
	o.Qty = updateQty
	o.AfterQty = qty
	if err := order.UpdateOrder(o); err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// 更新库存
	stockObj.Qty = qty
	stockObj.Updated = timex.String()
	if err := order.UpdateStock(stockObj); err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// res
	resData.Qty = updateQty

	log.Warn("result %s", resData.String())
	Server.BroadcastTo("login", "inventory", resData)
	return 1
}

func (c *CallbackController) Zero() {
	resData := string(c.Ctx.Input.RequestBody)
	log.Info("Zero: %s", resData)

	Server.BroadcastTo("login", "zero", resData)

	c.WriteHttpResponse(200, nil, nil)
	return
}

func (c *CallbackController) Measure() {
	resData := string(c.Ctx.Input.RequestBody)
	log.Info("Measure: %s", resData)

	Server.BroadcastTo("login", "measure", resData)

	c.WriteHttpResponse(200, nil, nil)
	return
}

// -----------------------------
// 自动盘点结果回调

func (c *CallbackController) Check() {
	log.Info("Check: %s", string(c.Ctx.Input.RequestBody))

	obj := &CbData{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	i := AutoInventory(obj.UUID, obj.Weight)

	log.Info("AutoInventory %d", i)

	c.WriteHttpResponse(200, nil, nil)
	return
}

func AutoInventory(orderIdStr string, gridWeight int) int {
	if len(orderIdStr) == 0 {
		log.Error("%v", errors.New("uuid is empty"))
		return 0
	}

	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// 查询auto
	o, err := order.AutoById(orderId)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	var qty int
	// 1
	// 查询物料传感器配置参数
	materialSensor, err := material.SensorByMaterialId(o.MaterialId, o.SensorId)
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	params, err := materialSensor.ParamsObj()
	if err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	if gridWeight > 0 {
		minQty := params.Weight - params.Lower
		maxQty := params.Weight + params.ComeUp
		log.Info("WEIGHT value %d, params weight %d, min %d, max %d", gridWeight, params.Weight, minQty, maxQty)

		// 取余数量
		num := gridWeight / params.Weight

		// 剩余数量
		subQty := gridWeight - (num * params.Weight)

		log.Info("weight %d", num)

		if minQty <= subQty && subQty <= maxQty {
			qty = num + 1
		} else {
			qty = num
		}

		log.Info("weight %d", qty)
	}

	log.Info("----------QTY---------- %d", qty)

	// 查询库存
	stockObj, err := order.StockByMaterialId(o.MaterialId, o.GridId)
	if err != nil {
		log.Error("StockByMaterialId %v", errors.As(err))
		return 1
	}

	// update
	// 更新auto
	o.Updated = timex.String()
	o.Qty = qty
	if err := order.UpdateAuto(o); err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	// 更新库存
	stockObj.Qty = qty
	stockObj.Updated = timex.String()
	if err := order.UpdateStock(stockObj); err != nil {
		log.Error("%v", errors.As(err))
		return 0
	}

	return 1
}

func (c *CallbackController) BoxStatus() {
	resData := string(c.Ctx.Input.RequestBody)
	log.Info("BoxStatus: %s", resData)

	Server.BroadcastTo("login", "boxStatus", resData)

	c.WriteHttpResponse(200, nil, nil)
	return
}

func (c *CallbackController) DoorStatus() {
	resData := string(c.Ctx.Input.RequestBody)
	log.Info("DoorStatus: %s", resData)

	Server.BroadcastTo("login", "doorStatus", resData)

	c.WriteHttpResponse(200, nil, nil)
	return
}

func (c *CallbackController) Light() {
	resData := string(c.Ctx.Input.RequestBody)
	log.Info("Light: %s", resData)

	Server.BroadcastTo("login", "light", resData)

	c.WriteHttpResponse(200, nil, nil)
	return
}

// -------------------------
// 扫条码回调

func (c *CallbackController) Barcode() {
	log.Info("Barcode: %s", string(c.Ctx.Input.RequestBody))

	obj := &SerialRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	Server.BroadcastTo("login", "scanner", obj.Data)

	c.WriteHttpResponse(200, nil, nil)
	return
}

// -------------------------
// 扫二维码回调

func (c *CallbackController) Qrcode() {
	log.Info("Qrcode: %s", string(c.Ctx.Input.RequestBody))

	obj := &SerialRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	Server.BroadcastTo("login", "scanner", obj.Data)

	c.WriteHttpResponse(200, nil, nil)
	return
}

// -------------------------
// 指纹

func (c *CallbackController) Finger() {
	log.Info("Finger: %s", string(c.Ctx.Input.RequestBody))

	obj := &SerialRequest{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(500, nil, err)
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	Server.BroadcastTo("login", "finger", obj)

	c.WriteHttpResponse(200, nil, nil)
	return
}
