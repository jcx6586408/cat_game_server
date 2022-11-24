package storage

import (
	"catLog"
	"config"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	db, _ := sql.Open("mysql", config.Read().DB)
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		fmt.Println("open database fail", err)
		return
	}
	DB = db
	CreateTable()
	fmt.Println("mysql数据库 connnect success")
}

func CreateTable() {
	sqlBytes, err := ioutil.ReadFile("./user.sql")
	if err != nil {
		catLog.Log("读取文件失败_", err)
		return
	}
	sqlTable := string(sqlBytes)
	result, err := DB.Exec(sqlTable)
	checkErr(err)
	catLog.Log("建表结果_", result)
}

func SaveToDB(uu *UserStorage) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		cc := queryuser(uu.Uid)
		result := <-cc
		if result == 1 {
			catLog.Log("开始查询角色")
			u := toDbUser(uu)
			ccc := updateuser(uu.Uid, "Forever", u.Forever)
			<-ccc
		} else {
			ccc := insectuser(uu)
			<-ccc
		}
	}()
	return c
}

func toDbUser(uu *UserStorage) *UserStorageDB {
	u := &UserStorageDB{}
	u.Uid = uu.Uid
	data, _ := json.Marshal(uu.Forever)
	u.Forever = string(data)
	return u
}

func toUser(u *UserStorageDB) *UserStorage {
	uu := &UserStorage{}
	uu.Uid = u.Uid
	forever := make(map[string]string)
	json.Unmarshal([]byte(u.Forever), &forever)
	uu.Forever = forever
	return uu
}

func insectuser(uu *UserStorage) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		tx, err := DB.Begin()
		if err != nil {
			fmt.Println("tx fail")
		}
		u := toDbUser(uu)
		stmt, err := tx.Prepare("INSERT user SET uid=?,Forever=?")
		checkErr(err)
		res, err := stmt.Exec(u.Uid, u.Forever)
		checkErr(err)
		//将事务提交
		tx.Commit()
		//获得上一个插入自增的id
		fmt.Println(res.LastInsertId())
		catLog.Log("插入完成")
	}()
	return c
}

func updateuser(uid string, fieldName string, content string) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		tx, err := DB.Begin()
		if err != nil {
			fmt.Println("tx fail")
		}
		var s = "update user set " + fieldName + "=? where uid=?"
		//更新資料
		stmt, err := tx.Prepare(s)
		checkErr(err)

		res, err := stmt.Exec(content, uid)
		checkErr(err)

		_, err = res.RowsAffected()
		checkErr(err)
		tx.Commit()
		catLog.Log("更新完成")
	}()
	return c
}

func queryuser(uid string) <-chan interface{} {
	//查詢資料
	c := make(chan interface{})
	go func() {
		defer close(c)
		row, err := DB.Query("SELECT * FROM user where uid=" + uid)
		if err != nil {
			c <- 2
			return
		}
		for row.Next() {
			var id int
			var uid string
			var Forever string
			err = row.Scan(&id, &uid, &Forever)
			catLog.Log("查询信息_", id, uid, Forever)
			c <- 1
		}
		checkErr(err)
	}()
	return c
}

func userFromDb(uid string, u *UserStorage) <-chan interface{} {
	//查詢資料
	c := make(chan interface{})
	go func() {
		defer close(c)
		row, err := DB.Query("SELECT * FROM user where uid=" + uid)
		if err != nil {
			return
		}
		for row.Next() {
			var id int
			var uid string
			var Forever string
			err = row.Scan(&id, &uid, &Forever)
			catLog.Log("查询信息_", id, uid, Forever)
			// forever
			forever := make(map[string]string)
			json.Unmarshal([]byte(Forever), &forever)
			u.Forever = forever
		}
		checkErr(err)
	}()
	return c
}

func NewUser() <-chan interface{} {
	//查詢資料
	c := make(chan interface{})
	go func() {
		defer close(c)
		row, err := DB.Query("SELECT * FROM user")
		for row.Next() {
			var id int
			var uid string
			var Forever string
			err = row.Scan(&id, &uid, &Forever)
			catLog.Log("查询信息_", id, uid, Forever)

			u := &UserStorageDB{Uid: uid, Forever: Forever}
			users[u.Uid] = toUser(u)
			c <- 1
		}
		checkErr(err)
	}()
	return c
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
