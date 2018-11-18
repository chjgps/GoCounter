package permission

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrRolePermissionNotFound = errors.New("role permission not found")
)

type RolePermission struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 角色ID
	RoleId int `orm:"column(role_id)" json:"roleId"`
	// 权限ID
	PermissionId int `orm:"column(permission_id)" json:"permissionId"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *RolePermission) TableName() string {
	return "rel_role_permission"
}

// 添加
func InsertRolePermission(obj *RolePermission) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateRolePermission(obj *RolePermission) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelRolePermission(id int) error {
	o := orm.NewOrm()

	obj := &RolePermission{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func RolePermissionById(roleId, permissionId int) (*RolePermission, error) {
	o := orm.NewOrm()

	obj := &RolePermission{
		RoleId:       roleId,
		PermissionId: permissionId,
	}

	if err := o.Read(obj, "RoleId", "PermissionId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrRolePermissionNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func PermissionByRoleId(roleId, page, pageSize int) ([]*Permission, error) {
	o := orm.NewOrm()

	list := []*Permission{}

	if _, err := o.Raw(permissionByRoleIdSql, roleId, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const permissionByRoleIdSql = `
SELECT
    t1.id,
    t1.created, 
    t1.created_by,
    t1.parent_id, 
    t1.tag,
    t1.label,
    t1.path,
    t1.icon,
    t1.level,
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    permission t1
INNER JOIN
    rel_role_permission t2
ON
    t1.id = t2.permission_id
WHERE
    t2.role_id = ?
ORDER BY
    t1.level
DESC
LIMIT ? OFFSET ?
`
