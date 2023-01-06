package msg

type Login struct {
	Uuid string
}

type LoginRequest struct {
	Uuid     string
	Icon     string
	Nickname string
}
