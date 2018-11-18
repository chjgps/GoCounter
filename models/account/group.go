package account

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrGroupNotFound    = errors.New("group not found")
	ErrNameAlreadyExist = errors.New("group name already exist")
)

type Group struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 名称
	Name string `orm:"column(name)" json:"name"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *Group) TableName() string {
	return "group"
}

// 添加
func InsertGroup(obj *Group) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelGroup(id int) error {
	o := orm.NewOrm()

	obj := &Group{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateGroup(obj *Group) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func GroupById(id int) (*Group, error) {
	o := orm.NewOrm()

	obj := &Group{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据名称查询
func GroupByName(name string) (*Group, error) {
	o := orm.NewOrm()

	obj := &Group{
		Name: name,
	}

	if err := o.Read(obj, "Name"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrGroupNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func GroupList(page, pageSize int) ([]*Group, error) {
	o := orm.NewOrm()

	list := []*Group{}
	if _, err := o.Raw(groupListSql, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const groupListSql = `
SELECT
    id, 
    created, 
    created_by,
    name, 
    status, 
    updated,
    updated_by
FROM
    'group'
LIMIT ? OFFSET ?
`
