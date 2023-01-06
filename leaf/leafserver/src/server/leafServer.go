package main

import (
	"config"
	"excel"
	"leafserver/src/server/conf"
	"leafserver/src/server/game"
	"leafserver/src/server/gate"
	"leafserver/src/server/login"
	"os"

	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
)

func main() {
	println("启动leaf服务器-------------------------------------------------")
	conf.ConfPath = os.Args[1]
	config.ConfPath = os.Args[2]
	config.RoomConfPath = os.Args[3]
	excel.TablePath = os.Args[4]
	config.IPLocationPath = os.Args[5]
	if len(os.Args) >= 7 {
		conf.Server.CertFile = os.Args[6]
	}
	if len(os.Args) >= 8 {
		conf.Server.KeyFile = os.Args[7]
	}
	println("配置加载完成-------------------------------------------------")
	conf.Read()
	log.Debug("服务器启动配置路径:\r\n %s\r\n %s\r\n %s\r\n %s", conf.ConfPath, config.ConfPath, config.RoomConfPath, excel.TablePath)
	lconf.LeafServerPath = conf.ConfPath
	lconf.ServerPath = config.ConfPath
	lconf.RoomPath = config.RoomConfPath
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
