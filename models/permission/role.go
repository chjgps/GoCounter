package permission

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrNameAlreadyExist = errors.New("name already exist")
)

type Role struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 姓名
	Name string `orm:"column(name);unique" json:"name"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *Role) TableName() string {
	return "role"
}

// 添加
func InsertRole(obj *Role) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelRole(id int) error {
	o := orm.NewOrm()

	obj := &Role{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateRole(obj *Role) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func RoleById(id int) (*Role, error) {
	o := orm.NewOrm()

	obj := &Role{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrRoleNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func RoleByName(name string) (*Role, error) {
	o := orm.NewOrm()

	obj := &Role{
		Name: name,
	}

	if err := o.Read(obj, "Name"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrRoleNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func RoleList(where map[string]interface{}, page, pageSize int) (int64, []*Role, error) {
	o := orm.NewOrm()

	list := []*Role{}

	sql := " 1 "
	if len(where) > 0 {
		startDate := where["startDate"]
		if startDate != "" {
			sql += " AND t1.created >= '" + fmt.Sprintf("%s", startDate) + "' "
		}

		endDate := where["endDate"]
		if endDate != "" {
			sql += " AND t1.created <= '" + fmt.Sprintf("%s", endDate) + "' "
		}

		name := where["name"]
		if name != "" {
			sql += " AND (t1.name LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(roleListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(roleListSql+sql+" ORDER BY id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const roleListCountSql = `
SELECT
    COUNT(*)
FROM
    role AS t1
WHERE
`

const roleListSql = `
SELECT
    t1.id, 
    t1.created, 
    t1.created_by,
    t1.name, 
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    role AS t1
WHERE
`
