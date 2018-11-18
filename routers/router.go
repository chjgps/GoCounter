package routers

import (
	// "net/http/httputil"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/beego/ms304w-client/basis/errors"
	l "github.com/beego/ms304w-client/basis/syslog"
	"github.com/beego/ms304w-client/controllers"
)

var log = l.New()

func init() {
	ns := beego.NewNamespace("/v1",
		// 域名验证
		beego.NSCond(func(ctx *context.Context) bool {
			/*
				r := ctx.Request
				bytes, err := httputil.DumpRequest(r, true)
				if err != nil {
					log.Warn("%v", err)
				} else {
					log.Info("request: %s", string(bytes))
				}
			*/

			return true
		}),

		// 权限验证
		beego.NSBefore(func(ctx *context.Context) {
			oauth, err := beego.AppConfig.Bool("oauth")
			if err != nil {
				panic(err)
			}

			if oauth {
				token := ctx.Request.Header.Get("Token")
				uri := ctx.Request.RequestURI
				log.Info("uri: %s, token: %s", uri, token)

				if uri == "/v1/account/login/card" ||
					uri == "/v1/account/login/password" ||
					uri == "/v1/permission/user/login" ||
					strings.Contains(uri, "/v1/account/user/card") {
					controllers.OAuth.Del(token)
					log.Info("uri %s", uri)
				} else {
					// ctx.WriteString("Forbidden")
					// 查询token是否存在
					// 是否过期
					log.Info("token %s", token)
					if !controllers.OAuth.Get(token) {
						log.Info("token %s", token)
						ctx.ResponseWriter.WriteHeader(401)
						data := controllers.NewHttpResponse(401, nil, errors.New("Forbidden"))
						ctx.WriteString(data.String())
						return
					}

					log.Info("token %s", token)
					// 设置时间延长
					controllers.OAuth.Set(token)
					log.Info("token %s", token)
				}
			}
		}),

		// --------------------------
		// Error
		beego.NSNamespace("/error",
			beego.NSRouter("/404", &controllers.ErrorController{}, "POST:Error404"),
			beego.NSRouter("/500", &controllers.ErrorController{}, "*:Error500"),
		),

		// --------------------------
		// Account
		beego.NSNamespace("/account",
			// user
			beego.NSRouter("/user", &controllers.AccountController{}, "POST:AddAccount"),
			beego.NSRouter("/user", &controllers.AccountController{}, "PUT:EditAccount"),
			beego.NSRouter("/user/:id:int", &controllers.AccountController{}, "DELETE:DelAccount"),
			beego.NSRouter("/user/:id:int", &controllers.AccountController{}, "GET:AccountById"),
			beego.NSRouter("/user/card/:cardId:string", &controllers.AccountController{}, "GET:AccountByCard"),
			beego.NSRouter("/login/card", &controllers.AccountController{}, "POST:LoginByCard"),
			beego.NSRouter("/login/username", &controllers.AccountController{}, "POST:LoginByUsername"),
			beego.NSRouter("/login/finger/:id:int", &controllers.AccountController{}, "GET:AccountByFinger"),
			beego.NSRouter("/user", &controllers.AccountController{}, "GET:AccountList"),
			beego.NSRouter("/user/password", &controllers.AccountController{}, "PUT:EditPassword"),
			// group
			beego.NSRouter("/group", &controllers.GroupController{}, "POST:AddGroup"),
			beego.NSRouter("/group", &controllers.GroupController{}, "PUT:EditGroup"),
			beego.NSRouter("/group/:id:int", &controllers.GroupController{}, "DELETE:DelGroup"),
			beego.NSRouter("/group/:id:int", &controllers.GroupController{}, "GET:GroupById"),
			beego.NSRouter("/group", &controllers.GroupController{}, "GET:GroupList"),
			// rel_group_user
			beego.NSRouter("/group/user", &controllers.AccountGroupController{}, "POST:AddAccountGroup"),
			// beego.NSRouter("/group/user", &controllers.AccountGroupController{}, "PUT:EditAccountGroup"),
			beego.NSRouter("/group/user/:groupId:int/:accountId:int", &controllers.AccountGroupController{}, "DELETE:DelAccountGroup"),
			beego.NSRouter("/group/user", &controllers.AccountGroupController{}, "GET:AccountByGroupId"),
		),

		// --------------------------
		// Supplier
		beego.NSNamespace("/supplier",
			// user
			beego.NSRouter("/", &controllers.SupplierController{}, "POST:AddSupplier"),
			beego.NSRouter("/", &controllers.SupplierController{}, "PUT:EditSupplier"),
			beego.NSRouter("/:id:int", &controllers.SupplierController{}, "DELETE:DelSupplier"),
			beego.NSRouter("/:id:int", &controllers.SupplierController{}, "GET:SupplierById"),
			beego.NSRouter("/", &controllers.SupplierController{}, "GET:SupplierList"),
			beego.NSRouter("/material", &controllers.SupplierController{}, "GET:MaterialList"),
		),

		// --------------------------
		// Material
		beego.NSNamespace("/material",
			// category
			beego.NSRouter("/category", &controllers.CategoryController{}, "POST:AddCategory"),
			beego.NSRouter("/category", &controllers.CategoryController{}, "PUT:EditCategory"),
			beego.NSRouter("/category/:id:int", &controllers.CategoryController{}, "DELETE:DelCategory"),
			beego.NSRouter("/category", &controllers.CategoryController{}, "GET:CategoryList"),

			// material
			beego.NSRouter("/", &controllers.MaterialController{}, "POST:AddMaterial"),
			beego.NSRouter("/", &controllers.MaterialController{}, "PUT:EditMaterial"),
			beego.NSRouter("/:id:int", &controllers.MaterialController{}, "DELETE:DelMaterial"),
			beego.NSRouter("/:id:int", &controllers.MaterialController{}, "GET:MaterialById"),
			// 根据分类ID查询所有物料
			beego.NSRouter("/", &controllers.MaterialController{}, "GET:MaterialList"),

			// rel_material_sensor
			beego.NSRouter("/sensor", &controllers.MaterialSensorController{}, "POST:AddSensor"),
			beego.NSRouter("/sensor", &controllers.MaterialSensorController{}, "PUT:EditSensor"),
			beego.NSRouter("/sensor/:id:int", &controllers.MaterialSensorController{}, "DELETE:DelSensor"),
			beego.NSRouter("/sensor/:id:int", &controllers.MaterialSensorController{}, "GET:SensorById"),
			beego.NSRouter("/sensor", &controllers.MaterialSensorController{}, "GET:SensorList"),

			// rel_group_material
			beego.NSRouter("/group", &controllers.GroupMaterialController{}, "POST:AddGroupMaterial"),
			// beego.NSRouter("/group", &controllers.GroupMaterialController{}, "PUT:EditGroupMaterial"),
			beego.NSRouter("/group/:groupId:int/:materialId:int", &controllers.GroupMaterialController{}, "DELETE:DelGroupMaterial"),
			beego.NSRouter("/group", &controllers.GroupMaterialController{}, "GET:MaterialByGroupId"),

			/*
				// rel_material_rfid
				beego.NSRouter("/rfid", &controllers.MaterialRfidController{}, "POST:AddMaterialRfid"),
				beego.NSRouter("/rfid", &controllers.MaterialRfidController{}, "PUT:EditMaterialRfid"),
				beego.NSRouter("/rfid/:id:int", &controllers.MaterialRfidController{}, "DELETE:DelMaterialRfid"),
				beego.NSRouter("/rfid/:rfid:string", &controllers.MaterialRfidController{}, "GET:MaterialByRfid"),
				beego.NSRouter("/rfid", &controllers.MaterialRfidController{}, "GET:MaterialRfidList"),
			*/
		),

		// --------------------------
		// Sensor
		beego.NSNamespace("/sensor",
			beego.NSRouter("/", &controllers.SensorController{}, "POST:AddSensor"),
			beego.NSRouter("/", &controllers.SensorController{}, "PUT:EditSensor"),
			beego.NSRouter("/:id:int", &controllers.SensorController{}, "DELETE:DelSensor"),
			beego.NSRouter("/:id:int", &controllers.SensorController{}, "GET:SensorById"),
			beego.NSRouter("/", &controllers.SensorController{}, "GET:SensorList"),
			// 物料和传感器的关系
			// beego.NSRouter("/material", &controllers.SensorController{}, "POST:AddMaterialSensor"),
			// beego.NSRouter("/material", &controllers.SensorController{}, "DELETE:DelMaterialSensor"),
			// beego.NSRouter("/material/:materialId:int/:sensorId:int", &controllers.SensorController{}, "GET:SensorByMaterialId"),
		),


		// --------------------------
		// Stock
		beego.NSNamespace("/stock",
			// 入库单
			beego.NSRouter("/in", &controllers.StockController{}, "POST:StockIn"),
			// 入库单确认
			beego.NSRouter("/in/confirm", &controllers.StockController{}, "POST:StockInConfirm"),
			// 出库单
			beego.NSRouter("/out", &controllers.StockController{}, "POST:StockOut"),
			// 出库单确认
			beego.NSRouter("/out/confirm", &controllers.StockController{}, "POST:StockOutConfirm"),
			// 回收单
			beego.NSRouter("/recycle", &controllers.StockController{}, "POST:StockRecycle"),
			// 回收单确认
			beego.NSRouter("/recycle/confirm", &controllers.StockController{}, "POST:StockRecycleConfirm"),
			// 格子库存
			beego.NSRouter("/", &controllers.StockController{}, "GET:StockList"),
			// 物料库存
			beego.NSRouter("/material", &controllers.StockController{}, "GET:MaterialStockList"),
			// 组物料库存
			beego.NSRouter("/material/group", &controllers.StockController{}, "GET:GroupStockList"),
			// 根据物料ID查询格子
			// 优先已存在相同物料格子
			beego.NSRouter("/grid/:materialId:int", &controllers.StockController{}, "GET:GridByMaterialId"),
			// 自动盘点
			beego.NSRouter("/auto", &controllers.StockController{}, "POST:Auto"),
			// 查询
			beego.NSRouter("/auto", &controllers.StockController{}, "GET:AutoList"),
		),

		/*
			// --------------------------
			// RFID
			beego.NSNamespace("/rfid",
				// 入库单
				beego.NSRouter("/in", &controllers.StockRfidController{}, "POST:StockIn"),
				// 入库单确认
				beego.NSRouter("/in/confirm", &controllers.StockRfidController{}, "POST:StockInConfirm"),
				// 出库单
				beego.NSRouter("/out", &controllers.StockRfidController{}, "POST:StockOut"),
				// 出库单确认
				beego.NSRouter("/out/confirm", &controllers.StockRfidController{}, "POST:StockOutConfirm"),
				// 回收单
				beego.NSRouter("/recycle", &controllers.StockRfidController{}, "POST:StockRecycle"),
				// 回收单确认
				beego.NSRouter("/recycle/confirm", &controllers.StockRfidController{}, "POST:StockRecycleConfirm"),
				// 库存
				beego.NSRouter("/stock", &controllers.StockRfidController{}, "GET:StockList"),
				// 根据用户查询订单
				beego.NSRouter("/order", &controllers.StockRfidController{}, "GET:OrderList"),
				// 待回收列表
				// beego.NSRouter("/recycle", &controllers.StockRfidController{}, "GET:RecycleList"),
			),
		*/

		// --------------------------
		// Code
		beego.NSNamespace("/code",
			// 入库单
			beego.NSRouter("/in", &controllers.StockCodeController{}, "POST:StockIn"),
			// 出库单
			beego.NSRouter("/out", &controllers.StockCodeController{}, "POST:StockOut"),
			// 回收单
			beego.NSRouter("/recycle", &controllers.StockCodeController{}, "POST:StockRecycle"),
		),

		// --------------------------
		// Order
		beego.NSNamespace("/order",
			beego.NSRouter("/", &controllers.OrderController{}, "POST:AddOrder"),
			beego.NSRouter("/", &controllers.OrderController{}, "PUT:EditOrder"),
			beego.NSRouter("/:id:int", &controllers.OrderController{}, "DELETE:DelOrder"),
			beego.NSRouter("/:id:int", &controllers.OrderController{}, "GET:OrderById"),
			beego.NSRouter("/account/:accountId:int", &controllers.OrderController{}, "GET:OrderByAccountId"),
			beego.NSRouter("/", &controllers.OrderController{}, "GET:OrderList"),
			beego.NSRouter("/recycle", &controllers.OrderController{}, "GET:RecycleList"),
		),

		// --------------------------
		// Conf
		beego.NSNamespace("/conf",
			beego.NSRouter("/auto", &controllers.AutoConfController{}, "POST:AddAutoConf"),
			beego.NSRouter("/auto", &controllers.AutoConfController{}, "PUT:EditAutoConf"),
			beego.NSRouter("/auto/:id:int", &controllers.AutoConfController{}, "DELETE:DelAutoConf"),
			beego.NSRouter("/auto/:id:int", &controllers.AutoConfController{}, "GET:AutoConfById"),
			beego.NSRouter("/auto", &controllers.AutoConfController{}, "GET:AutoConfList"),
		),

		// --------------------------
		// File
		beego.NSNamespace("/file",
			beego.NSRouter("/material", &controllers.FileController{}, "POST:AddMaterial"),
		),

		// --------------------------
		// Serial
		beego.NSNamespace("/serial",
			// 开门
			beego.NSRouter("/open", &controllers.SerialController{}, "POST:Open"),
			// 称重
			beego.NSRouter("/weight", &controllers.SerialController{}, "POST:Weight"),
			// 矫正
			// 第1步:清0
			beego.NSRouter("/weight/zero", &controllers.SerialController{}, "POST:Zero"),
			// 第2步:矫正
			beego.NSRouter("/weight/measure", &controllers.SerialController{}, "POST:Measure"),
			// 盘点
			beego.NSRouter("/weight/check", &controllers.SerialController{}, "POST:Check"),
			// 获取所有门状态
			beego.NSRouter("/box/:id:int/status", &controllers.SerialController{}, "GET:BoxStatus"),
			// 开关灯
			beego.NSRouter("/light", &controllers.SerialController{}, "POST:Light"),
		),

		// --------------------------
		// Callback
		beego.NSNamespace("/callback",
			// 刷卡
			beego.NSRouter("/card", &controllers.CallbackController{}, "POST:Card"),
			// 称重回调
			beego.NSRouter("/weight", &controllers.CallbackController{}, "POST:Weight"),
			// 清0回调
			beego.NSRouter("/weight/zero", &controllers.CallbackController{}, "POST:Zero"),
			// 矫正回调
			beego.NSRouter("/weight/measure", &controllers.CallbackController{}, "POST:Measure"),
			// 盘点
			beego.NSRouter("/weight/check", &controllers.CallbackController{}, "POST:Check"),
			// 所有门状态回调
			beego.NSRouter("/box/status", &controllers.CallbackController{}, "POST:BoxStatus"),
			// 当前门状态回调
			beego.NSRouter("/door/status", &controllers.CallbackController{}, "POST:DoorStatus"),
			// 开关灯回调
			beego.NSRouter("/light/status", &controllers.CallbackController{}, "POST:Light"),
			// 扫条码
			beego.NSRouter("/barcode", &controllers.CallbackController{}, "POST:Barcode"),
			// 扫二维码
			beego.NSRouter("/qrcode", &controllers.CallbackController{}, "POST:Qrcode"),
			// 指纹
			beego.NSRouter("/finger", &controllers.CallbackController{}, "POST:Finger"),
		),

		// --------------------------
		// Permission
		beego.NSNamespace("/permission",
			// user
			beego.NSRouter("/user", &controllers.UserController{}, "POST:AddUser"),
			beego.NSRouter("/user", &controllers.UserController{}, "PUT:EditUser"),
			beego.NSRouter("/user/password", &controllers.UserController{}, "PUT:EditPassword"),
			beego.NSRouter("/user/:id:int", &controllers.UserController{}, "DELETE:DelUser"),
			beego.NSRouter("/user/:id:int", &controllers.UserController{}, "GET:UserById"),
			beego.NSRouter("/user", &controllers.UserController{}, "GET:UserList"),
			beego.NSRouter("/user/login", &controllers.UserController{}, "POST:Login"),
			// role
			beego.NSRouter("/role", &controllers.RoleController{}, "POST:AddRole"),
			beego.NSRouter("/role", &controllers.RoleController{}, "PUT:EditRole"),
			beego.NSRouter("/role/:id:int", &controllers.RoleController{}, "DELETE:DelRole"),
			beego.NSRouter("/role/:id:int", &controllers.RoleController{}, "GET:RoleById"),
			beego.NSRouter("/role", &controllers.RoleController{}, "GET:RoleList"),
			// permission
			beego.NSRouter("/permission", &controllers.PermissionController{}, "POST:AddPermission"),
			beego.NSRouter("/permission", &controllers.PermissionController{}, "PUT:EditPermission"),
			beego.NSRouter("/permission/:id:int", &controllers.PermissionController{}, "DELETE:DelPermission"),
			beego.NSRouter("/permission/:id:int", &controllers.PermissionController{}, "GET:PermissionById"),
			beego.NSRouter("/permission", &controllers.PermissionController{}, "GET:PermissionList"),
			// userRole
			beego.NSRouter("/user/role", &controllers.UserRoleController{}, "POST:AddUserRole"),
			beego.NSRouter("/user/role/:userId:int/:roleId:int", &controllers.UserRoleController{}, "DELETE:DelUserRole"),
			// 角色查用户
			beego.NSRouter("/user/role", &controllers.UserRoleController{}, "GET:UserByRoleId"),
			// 用户查角色
			beego.NSRouter("/role/user", &controllers.UserRoleController{}, "GET:RoleByUserId"),
			// rolePermission
			beego.NSRouter("/role/permission", &controllers.RolePermissionController{}, "POST:AddRolePermission"),
			beego.NSRouter("/role/permission/:roleId:int/:permissionId:int", &controllers.RolePermissionController{}, "DELETE:DelRolePermission"),
			// 角色权限
			beego.NSRouter("/role/permission", &controllers.RolePermissionController{}, "GET:PermissionByRoleId"),
			// 用户权限
			beego.NSRouter("/user/permission", &controllers.RolePermissionController{}, "GET:PermissionByUserId"),
		),
	)
	beego.AddNamespace(ns)
}
