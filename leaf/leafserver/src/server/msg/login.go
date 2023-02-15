package msg

type Login struct {
	Uuid string
}

type LoginReply struct {
	State bool
	Value string
}

type LoginRequest struct {
	Uuid     string
	Uid      string
	Icon     string
	Nickname string
}

type Ping struct {
	Uuid string
}

type Pong struct {
}
