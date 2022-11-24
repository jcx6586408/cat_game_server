package remotemsg

var (
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
)
