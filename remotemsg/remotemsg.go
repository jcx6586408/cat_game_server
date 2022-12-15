package remotemsg

const (
	// 排行榜
	RANKPULL   = 101 // 拉取排行榜
	RANKUPDATE = 102 // 更新排行榜
	RANKSELF   = 103 // 排行榜自己

	// 微信openid
	WXOPENID = 901

	// Excel表格配置拉取
	EXCELCONFIG = 801

	// 存储Storage
	STORAGE       = 201 // 拉起存储
	STORAGEUPDATE = 202 // 存储更新

	// 匹配消息
	ROOMCREATE            = 301 // 创建房间
	ROOMLEAVE             = 302 // 离开房间
	ROOMADD               = 303 // 加入房间
	ROOMMATCH             = 304 // 单人匹配
	ROOMSTARTPLAY         = 305 // 比赛开始
	ROOMENDPLAY           = 306 // 比赛结束
	ROOMANSWEREND         = 307 // 单次答题结束
	ROOMANSWERSTART       = 318 // 单次答题开始
	ROOMOVER              = 308 // 房间结束解散
	ROOMPREPARE           = 309 // 准备
	ROOMPREPARECANCEL     = 310 // 取消准备
	ROOMCHANGEMASTER      = 311 // 转移房主
	ROOMTIME              = 312 // 计时
	ROOMMATCHROOMCANCEL   = 313 // 取消匹配房间的准备
	ROOMMATCHMEMBERCANCEL = 314 // 取消个人匹配的准备
	ROOMMATCHROOM         = 315 // 房间匹配
	ROOMANSWER            = 316 // 回答题目
	ROOMGET               = 317 // 主动请求房间信息
)
