package remotemsg

const (
	LOGIN     = 1 // 登录消息ID
	HEARTBEAT = 2 // 心跳消息ID
	//-------------------------以上消息皆是特定不能改的消息ID，为框架运行基础，请勿修改----------------------------------

	//====================业务消息ID=======================

	// 排行榜
	RANKPULL   = 101 // 拉取排行榜
	RANKUPDATE = 102 // 更新排行榜

	// 微信openid
	WXOPENID = 901

	// Excel表格配置拉取
	EXCELCONFIG = 801

	// 存储Storage
	STORAGE       = 201 // 拉起存储
	STORAGEUPDATE = 202 // 存储更新

	// 匹配消息
	ROOMCREATE        = 301 // 创建房间
	ROOMLEAVE         = 302 // 离开房间
	ROOMADD           = 303 // 加入房间
	ROOMMATCH         = 304 // 单人匹配
	ROOMSTARTPLAY     = 305 // 比赛开始
	ROOMENDPLAY       = 306 // 比赛结束
	ROOMANSWEREND     = 307 // 单次答题结束
	ROOMOVER          = 308 // 房间结束解散
	ROOMPREPARE       = 309 // 准备
	ROOMPREPARECANCEL = 310 // 取消准备
	ROOMCHANGEMASTER  = 311 // 转移房主
	ROOMTIME          = 312 // 计时
)
