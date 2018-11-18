package material

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrGroupMaterialNotFound = errors.New("account material not found")
)

type GroupMaterial struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 组ID
	GroupId int `orm:"column(group_id)" json:"groupId"`
	// 物料ID
	MaterialId int `orm:"column(material_id)" json:"materialId"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *GroupMaterial) TableName() string {
	return "rel_group_material"
}

// 添加
func InsertGroupMaterial(obj *GroupMaterial) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateGroupMaterial(obj *GroupMaterial) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelGroupMaterial(id int) error {
	o := orm.NewOrm()

	obj := &GroupMaterial{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func GroupMaterialById(groupId, materialId int) (*GroupMaterial, error) {
	o := orm.NewOrm()

	obj := &GroupMaterial{
		GroupId:    groupId,
		MaterialId: materialId,
	}

	if err := o.Read(obj, "GroupId", "MaterialId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrGroupMaterialNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func MaterialByGroupId(groupId, page, pageSize int) ([]*Material, error) {
	o := orm.NewOrm()

	list := []*Material{}

	if _, err := o.Raw(materialByGroupIdSql, groupId, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const materialByGroupIdSql = `
SELECT
    t1.id, 
    t1.created, 
    t1.created_by,
    t1.name, 
    t1.category_id, 
    t1.material_code,
    t1.material_spec,
    t1.supplier,
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    material t1
INNER JOIN
    rel_group_material t2
ON
    t1.id = t2.material_id
WHERE
    t2.group_id = ?
LIMIT ? OFFSET ?
`
