package account

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrAccountGroupNotFound = errors.New("account group not found")
)

type AccountGroup struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 组ID
	GroupId int `orm:"column(group_id)" json:"groupId"`
	// 客户ID
	AccountId int `orm:"column(account_id)" json:"accountId"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *AccountGroup) TableName() string {
	return "rel_account_group"
}

// 添加
func InsertAccountGroup(obj *AccountGroup) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateAccountGroup(obj *AccountGroup) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelAccountGroup(id int) error {
	o := orm.NewOrm()

	obj := &AccountGroup{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func AccountGroupById(groupId, accountId int) (*AccountGroup, error) {
	o := orm.NewOrm()

	obj := &AccountGroup{
		GroupId:   groupId,
		AccountId: accountId,
	}

	if err := o.Read(obj, "GroupId", "AccountId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountGroupNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func AccountByGroupId(groupId, page, pageSize int) ([]*Account, error) {
	o := orm.NewOrm()

	list := []*Account{}
	if _, err := o.Raw(accountByGroupIdSql, groupId, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const accountByGroupIdSql = `
SELECT
    t1.id,
    t1.created, 
    t1.created_by,
    t1.username, 
    t1.role, 
    t1.card,
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    account t1
INNER JOIN
    rel_account_group t2
ON
    t1.id = t2.account_id
WHERE
    t2.group_id = ?
ORDER BY
    t2.created
DESC
LIMIT ? OFFSET ?
`
