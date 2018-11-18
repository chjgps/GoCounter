package permission

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/basis/errors"
)

var (
	ErrUserRoleNotFound = errors.New("user role not found")
)

type UserRole struct {
	Id        int    `orm:"column(id);auto;pk" json:"id"`
	Created   string `orm:"column(created)" json:"created"`
	CreatedBy string `orm:"column(created_by)" json:"createdBy"`
	// 客户ID
	UserId int `orm:"column(user_id)" json:"userId"`
	// 角色ID
	RoleId int `orm:"column(role_id)" json:"roleId"`
	// 状态0停用1启用
	Status int `orm:"column(status);default(1)" json:"status"`
	// 更新时间
	Updated   string `orm:"column(updated)" json:"updated"`
	UpdatedBy string `orm:"column(updated_by)" json:"updatedBy"`
}

func (t *UserRole) TableName() string {
	return "rel_user_role"
}

// 添加
func InsertUserRole(obj *UserRole) error {
	o := orm.NewOrm()

	if _, err := o.Insert(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 修改
func UpdateUserRole(obj *UserRole) error {
	o := orm.NewOrm()

	if _, err := o.Update(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 删除
func DelUserRole(id int) error {
	o := orm.NewOrm()

	obj := &UserRole{
		Id: id,
	}

	if _, err := o.Delete(obj); err != nil {
		return errors.As(err)
	}

	return nil
}

// 根据ID查询
func UserRoleById(userId, roleId int) (*UserRole, error) {
	o := orm.NewOrm()

	obj := &UserRole{
		UserId: userId,
		RoleId: roleId,
	}

	if err := o.Read(obj, "UserId", "RoleId"); err != nil {
		if err == orm.ErrNoRows {
			return nil, errors.As(ErrUserRoleNotFound)
		}

		return nil, errors.As(err)
	}

	return obj, nil
}

// 查询所有
func UserByRoleId(roleId, page, pageSize int) ([]*User, error) {
	o := orm.NewOrm()

	list := []*User{}
	if _, err := o.Raw(userByRoleIdSql, roleId, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const userByRoleIdSql = `
SELECT
    t1.id,
    t1.created, 
    t1.created_by,
    t1.username, 
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    user t1
INNER JOIN
    rel_user_role t2
ON
    t1.id = t2.user_id
WHERE
    t2.role_id = ?
ORDER BY
    t2.created
DESC
LIMIT ? OFFSET ?
`

// 查询所有角色
func RoleByUserId(userId, page, pageSize int) ([]*Role, error) {
	o := orm.NewOrm()

	list := []*Role{}
	if _, err := o.Raw(roleByUserIdSql, userId, pageSize, (page-1)*pageSize).QueryRows(&list); err != nil {
		return nil, errors.As(err)
	}

	return list, nil
}

const roleByUserIdSql = `
SELECT
    t1.id,
    t1.created, 
    t1.created_by,
    t1.name, 
    t1.status, 
    t1.updated,
    t1.updated_by
FROM
    role t1
INNER JOIN
    rel_user_role t2
ON
    t1.id = t2.role_id
WHERE
    t2.user_id = ?
ORDER BY
    t2.created
DESC
LIMIT ? OFFSET ?
`
