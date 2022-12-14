package internal

type Roombaseer interface {
	GetID() int // 房间ID
	OnClose()   // 房间关闭监听
	OnInit()    // 房间启用初始化

	Full() bool          // 房间是否满员
	GetMemberCount() int // 获取当前成员数量
}
