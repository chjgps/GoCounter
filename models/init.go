package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/ms304w-client/models/account"
	"github.com/beego/ms304w-client/models/box"
	"github.com/beego/ms304w-client/models/material"
	"github.com/beego/ms304w-client/models/order"
	"github.com/beego/ms304w-client/models/permission"
	"github.com/beego/ms304w-client/models/sensor"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func init() {
	// db
	driver := beego.AppConfig.String("db_driver")
	dsn := beego.AppConfig.String("db_dsn")
	debug, err := beego.AppConfig.Bool("db_debug")
	if err != nil {
		panic(err)
	}

	// 设置为 UTC 时间
	// sqlite3 默认使用UTC时间
	// 强制用UTC时间
	orm.DefaultTimeLoc = time.UTC

	if err := orm.RegisterDriver(driver, orm.DRSqlite); err != nil {
		panic(err)
	}

	if err := orm.RegisterDataBase("default", driver, dsn); err != nil {
		panic(err)
	}

	// models
	orm.RegisterModel(
		// account
		new(account.Account),
		new(account.Group),
		new(account.AccountGroup),
		// material
		new(material.Material),
		new(material.GroupMaterial),
		new(material.Category),
		new(material.Sensor),
		new(material.MaterialRfid),
		new(material.Supplier),
		// box
		new(box.Box),
		new(box.Grid),
		new(box.Channel),
		new(box.Account),
		new(box.Correct),
		// sensor
		new(sensor.Sensor),
		// order
		new(order.Order),
		new(order.OrderRfid),
		new(order.Stock),
		new(order.Detail),
		new(order.Auto),
		new(order.AutoConf),
		// permission
		new(permission.User),
		new(permission.Role),
		new(permission.Permission),
		new(permission.UserRole),
		new(permission.RolePermission),
	)

	// sync
	if err := orm.RunSyncdb("default", false, true); err != nil {
		panic(err)
	}

	orm.SetMaxIdleConns("default", 100)
	// orm.SetMaxOpenConns("default", 50)
	orm.SetMaxOpenConns("default", 1)

	orm.Debug = debug
}
