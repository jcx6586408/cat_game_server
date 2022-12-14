package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/name5566/leaf/log"
)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	HttpAddr    string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
}

func init() {
	data, err := ioutil.ReadFile("conf/leafserver.json")
	if err != nil {
		// log.Logger("%v", err)
		data, err = ioutil.ReadFile("./conf/leafserver.json")
		if err != nil {
			log.Fatal("%v", err)
		}
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
