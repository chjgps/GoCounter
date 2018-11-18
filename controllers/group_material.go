package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/beego/ms304w-client/basis/errors"
	"github.com/beego/ms304w-client/basis/timex"
	"github.com/beego/ms304w-client/models/material"
)

type GroupMaterialController struct {
	BaseController
}

// 添加
func (c *GroupMaterialController) AddGroupMaterial() {
	obj := &material.GroupMaterial{}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &obj); err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	if obj == nil {
		c.WriteHttpResponse(400, nil, errors.New("params is empty"))
		return
	}

	groupId := obj.GroupId
	materialId := obj.MaterialId

	// 查询是否存在
	gm, err := material.GroupMaterialById(groupId, materialId)
	if err != nil {
		if !material.ErrGroupMaterialNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}
	}

	// 不存在添加
	if gm == nil {
		obj.Status = 1
		obj.Created = timex.String()
		if err := material.InsertGroupMaterial(obj); err != nil {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(200, nil, nil)
		return
	}

	// 存在修改
	gm.Status = 1
	gm.Updated = timex.String()
	if err := material.UpdateGroupMaterial(gm); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 删除组物料
func (c *GroupMaterialController) DelGroupMaterial() {
	groupIdStr := c.Ctx.Input.Param(":groupId")
	log.Debug(groupIdStr)
	if len(groupIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("group id is empty"))
		return
	}

	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	materialIdStr := c.Ctx.Input.Param(":materialId")
	log.Debug(materialIdStr)
	if len(materialIdStr) == 0 {
		c.WriteHttpResponse(400, nil, errors.New("material id is empty"))
		return
	}

	materialId, err := strconv.Atoi(materialIdStr)
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
	}

	// 查询是否存在
	gm, err := material.GroupMaterialById(groupId, materialId)
	if err != nil {
		if !material.ErrGroupMaterialNotFound.Equal(err) {
			c.WriteHttpResponse(500, nil, errors.As(err))
			return
		}

		c.WriteHttpResponse(404, nil, errors.As(err))
		return
	}

	// 存在
	if err := material.DelGroupMaterial(gm.Id); err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, nil, nil)
	return
}

// 根据组ID查询所有物料
func (c *GroupMaterialController) MaterialByGroupId() {
	groupId, err := c.GetInt("groupId")
	if err != nil {
		c.WriteHttpResponse(400, nil, errors.As(err))
		return
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

	list, err := material.MaterialByGroupId(groupId, page, pageSize)
	if err != nil {
		c.WriteHttpResponse(500, nil, errors.As(err))
		return
	}

	c.WriteHttpResponse(200, list, nil)
	return
}
