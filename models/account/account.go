package account

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrUsernameAlreadyExist = errors.New("username already exist")
	ErrCardAlreadyExist     = errors.New("card already exist")
)

const (
	// 操作员
	NORMAL_USER = iota + 1
	// 维护员
	ADMIN_USER
	// 系统用户
	SYS_USER
)

type Account struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 姓名
	Username string `orm:"column(username);unique" json:"username"`
	// 角色
	Role string `orm:"column(role)" json:"role"`
	// 卡
	Card string `orm:"column(card);unique" json:"card"`
	// 密码
	Password string `orm:"column(password)" json:"password"`
	// 指纹
	Finger int `orm:"column(finger);unique" json:"finger"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// other
	Token   string `json:"token"`
	GroupId int    `json:"groupId"`
}

func (t *Account) TableName() string {
	return "account"
}

// 添加
func InsertAccount(obj *Account) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelAccount(id int) error {
	o := orm.NewOrm()

	obj := &Account{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateAccount(obj *Account) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func AccountById(id int) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据用户名查询
func AccountByName(username string) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Username: username,
	}

	if err := o.Read(obj, "Username"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据卡查询
func AccountByCard(card string) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Card: card,
	}

	if err := o.Read(obj, "Card"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据卡登录
func LoginByCard(card, password string) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Card:     card,
		Password: password,
	}

	if err := o.Read(obj, "Card", "Password"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据用户名登录
func LoginByUsername(username, password string) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Username: username,
		Password: password,
	}

	if err := o.Read(obj, "Username", "Password"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func AccountList(where map[string]interface{}, page, pageSize int) (int64, []*Account, error) {
	o := orm.NewOrm()

	list := []*Account{}

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
			sql += " AND (t1.username LIKE '%" + fmt.Sprintf("%s", name) + "%'"
			sql += " OR t1.card LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(accountListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(accountListSql+sql+" ORDER BY t1.id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const accountListCountSql = `
SELECT
    COUNT(*)
FROM
    account AS t1
WHERE
`

const accountListSql = `
SELECT
    t1.id, 
    t1.created, 
    t1.created_by,
    t1.username, 
    t1.role, 
    t1.card,
    t1.status, 
    t1.updated,
    t1.updated_by,
    t2.group_id
FROM
    account AS t1
LEFT JOIN
    rel_account_group AS t2
ON
    t1.id = t2.account_id
WHERE
`

// 根据指纹查询
func AccountByFinger(finger int) (*Account, error) {
	o := orm.NewOrm()

	obj := &Account{
		Finger: finger,
	}

	if err := o.Read(obj, "Finger"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrAccountNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}
