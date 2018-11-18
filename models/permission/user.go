package permission

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameAlreadyExist = errors.New("username already exist")
)

type User struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 姓名
	Username string `orm:"column(username);unique" json:"username"`
	// 密码
	Password string `orm:"column(password)" json:"password"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`

	// other
	Token string `json:"token"`
}

func (t *User) TableName() string {
	return "user"
}

// 添加
func InsertUser(obj *User) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelUser(id int) error {
	o := orm.NewOrm()

	obj := &User{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateUser(obj *User) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func UserById(id int) (*User, error) {
	o := orm.NewOrm()

	obj := &User{
		Id: id,
	}

	if err := o.Read(obj, "Id"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrUserNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据用户名查询
func UserByName(username string) (*User, error) {
	o := orm.NewOrm()

	obj := &User{
		Username: username,
	}

	if err := o.Read(obj, "Username"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrUserNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 根据用户名登录
func LoginByUsername(username, password string) (*User, error) {
	o := orm.NewOrm()

	obj := &User{
		Username: username,
		Password: password,
	}

	if err := o.Read(obj, "Username", "Password"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrUserNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func UserList(where map[string]interface{}, page, pageSize int) (int64, []*User, error) {
	o := orm.NewOrm()

	list := []*User{}

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
			sql += " AND (t1.username LIKE '%" + fmt.Sprintf("%s", name) + "%')"
		}
	}

	sql += " AND 1 "

	// 查询总数
	var total int64
	if err := o.Raw(userListCountSql + sql).QueryRow(&total); err != nil {
		return -1, nil, errors.As(err)
	}

	// 查询所有
	if _, err := o.Raw(userListSql+sql+" ORDER BY id LIMIT ? OFFSET ?", pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return -1, nil, errors.As(err)
	}

	return total, list, nil
}

const userListCountSql = `
SELECT
    COUNT(*)
FROM
    user AS t1
WHERE
`

const userListSql = `
SELECT
    t1.id, 
    t1.created, 
    t1.created_by,
    t1.username, 
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    user AS t1
WHERE
`
