package permission

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrPermissionNotFound     = errors.New("permission not found")
	ErrPermissionAlreadyExist = errors.New("permission already exist")
)

type Permission struct {
	// 类别ID
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 父类ID
	ParentId int `orm:"column(parent_id)" json:"parentId"`
	// 标签
	Tag string `orm:"column(tag);unique" json:"tag"`
	// 标注
	Label string `orm:"column(label)" json:"label"`
	// 路径
	Path string `orm:"column(path)" json:"path"`
	// ICON
	Icon string `orm:"column(icon)" json:"icon"`
	// 等级
	Level int `orm:"column(level)" json:"level"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *Permission) TableName() string {
	return "permission"
}

// 添加
func InsertPermission(obj *Permission) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdatePermission(obj *Permission) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelPermission(id int) error {
	o := orm.NewOrm()

	obj := &Permission{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func PermissionById(id int) (*Permission, error) {
	o := orm.NewOrm()

	obj := &Permission{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrPermissionNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据tag查询
func PermissionByTag(tag string) (*Permission, error) {
	o := orm.NewOrm()

	obj := &Permission{
		Tag: tag,
	}

	if err := o.Read(obj, "Tag"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrPermissionNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询分类
func PermissionList() ([]*Permission, error) {
	o := orm.NewOrm()

	list := []*Permission{}

	if _, err := o.Raw(permissionListSql).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const permissionListSql = `
WITH RECURSIVE
MATERIAL_CTE(id, parent_id, tag, label, path, icon, level, floor) AS (
    SELECT
        t1.id,
        t1.parent_id,
        t1.tag,
        t1.label,
        t1.path,
        t1.icon,
        t1.level,
		0 AS floor
    FROM
        permission AS t1
    WHERE
        t1.parent_id = 0
    UNION ALL
        SELECT
            t2.id,
            t2.parent_id,
            t2.tag,
            t2.label,
            t2.path,
            t2.icon,
			t2.level,
            t3.floor + 1
        FROM
            permission AS t2
        INNER JOIN MATERIAL_CTE AS t3 ON t2.parent_id = t3.id
			ORDER BY level DESC
)
SELECT
    *
FROM
    MATERIAL_CTE;
`
