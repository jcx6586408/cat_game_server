package storage

import (
	"config"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-mysql-org/go-mysql/client"
	_ "github.com/go-sql-driver/mysql"
	"github.com/name5566/leaf/log"
)

type UserStorageDB struct {
	Uid      string // 用户uid
	Nickname string // 昵称
	Icon     string // 头像
	Online   int    // 是否在线
	Forever  string // 永久存储信息
}

type UserStorage struct {
	Uid          string            `bson:"uid"`
	Nickname     string            `bson:"nickname"`
	Icon         string            `bson:"icon"`
	Online       int               `bson:"online"`
	RegisterTime time.Time         `bson:"registerTime"`
	Forever      map[string]string `bson:"forever"`
}

var (
	Conf *config.Config
	Pool *client.Pool
)

func Connect() {
	Conf = config.Read()
	Pool = client.NewPool(log.Debug, Conf.DB.MinAlive, Conf.DB.MaxAlive, Conf.DB.MaxIdle, Conf.DB.Host, Conf.DB.User, Conf.DB.Password, Conf.DB.DB)
	conn, _ := Pool.GetConn(context.Background())
	defer Pool.PutConn(conn)
	sqlBytes, err := ioutil.ReadFile("./user.sql")
	if err != nil {
		return
	}
	sqlTable := string(sqlBytes)
	conn.Execute(sqlTable) // 初始建表
}

func OfflineHandle(uu *UserStorage) {
	result := Queryuser(uu.Uid)
	uu.Online = 0 // 变更在线状态为离线
	if result != nil {
		u := toDbUser(uu)
		Updateuser(uu, "Forever", u.Forever) // 更新用户
	} else {
		Insectuser(uu) // 插入用户
	}
}

func toDbUser(uu *UserStorage) *UserStorageDB {
	u := &UserStorageDB{}
	u.Uid = uu.Uid
	u.Nickname = uu.Nickname
	u.Icon = uu.Icon
	data, _ := json.Marshal(uu.Forever)
	u.Forever = string(data)
	return u
}

func Insectuser(uu *UserStorage) {
	conn, _ := Pool.GetConn(context.Background())
	defer Pool.PutConn(conn)
	u := toDbUser(uu)
	str := fmt.Sprintf(`insert into user (uid, nickname, icon, online, Forever) values ('%s', '%s', '%s', '%d','%s')`, uu.Uid, uu.Nickname, uu.Icon, uu.Online, u.Forever)
	_, err := conn.Execute(str)
	if err != nil {
		log.Error("插入失败%v", err)
	}
}

func UpdateOnline(uu *UserStorage, online int) {
	conn, _ := Pool.GetConn(context.Background())
	defer Pool.PutConn(conn)
	_, err := conn.Execute(fmt.Sprintf(`UPDATE user SET online='%d', where uid = '%s'`, uu.Online, uu.Uid))
	if err != nil {
		log.Debug("更新出错%v", err)
		return
	}
}

func Updateuser(uu *UserStorage, fieldName string, content string) {
	conn, _ := Pool.GetConn(context.Background())
	defer Pool.PutConn(conn)
	_, err := conn.Execute(fmt.Sprintf(`UPDATE user SET nickname='%s', icon='%s', online='%d',Forever='%s' where uid = '%s'`, uu.Nickname, uu.Icon, uu.Online, content, uu.Uid))
	if err != nil {
		log.Debug("更新出错%v", err)
		return
	}
}

func Queryuser(uid string) *UserStorage {
	//查詢資料
	conn, _ := Pool.GetConn(context.Background())
	defer Pool.PutConn(conn)
	// Select
	r, err := conn.Execute(fmt.Sprintf(`select uid, nickname, icon, online, Forever from user where uid = '%s'`, uid))
	if err != nil {
		log.Debug("查找出错%v", err)
	}
	if len(r.Values) > 0 {
		uid, _ := r.GetStringByName(0, "uid")
		nickName, _ := r.GetStringByName(0, "nickname")
		icon, _ := r.GetStringByName(0, "icon")
		online, _ := r.GetIntByName(0, "online")
		ff, _ := r.GetStringByName(0, "Forever")
		forever := make(map[string]string)
		log.Debug("结果数量:%v, 查找结果: %v-%v-%v-%v", len(r.Values), uid, nickName, icon, ff)
		json.Unmarshal([]byte(ff), &forever)
		return &UserStorage{
			Uid:      uid,
			Nickname: nickName,
			Icon:     icon,
			Online:   int(online),
			Forever:  forever,
		}
	}
	return nil
}
