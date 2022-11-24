package client

type RemoteMsg struct {
	ID   int         `json:"id"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func NewRemoteMsg(id int, code int, data interface{}) *RemoteMsg {
	msg := &RemoteMsg{}
	msg.ID = id
	msg.Code = code
	msg.Data = data
	return msg
}
