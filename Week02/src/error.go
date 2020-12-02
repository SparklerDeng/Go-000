package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
)

var RepositoryError = errors.New("数据资源错误！")
var DataNotFound = errors.New("无此id用户数据！")

var InvalidParam = errors.New("非法参数！")

type User struct {
	Id   int64
	Name string
}

func main() {
	user, err := biz(3)
	if err != nil {
		log.Fatalf("查询用户失败，原因: %+v", err)
		//err = errors.Unwrap(err)
		//log.Fatalf("查询用户失败，原因: %+v", err)
	}
	fmt.Println(user)
}

func biz(id int64) (user *User, err error) {
	if id <= 0 {
		return nil, errors.Wrap(InvalidParam, fmt.Sprintf("id不能小于1，当前值：%d", id))
	}
	return FindById(id)
}

func FindById(id int64) (user *User, err error) {
	db, err := dbConn()
	if err != nil {
		err = errors.Wrap(RepositoryError, "连接数据资源失败！")
		return
	}
	stmt, err := db.Prepare("select id, name from go_user where id=?")
	if err != nil {
		return user, errors.Wrap(RepositoryError, "查询数据资源失败！")
	}
	user = &User{}
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Wrap(DataNotFound, fmt.Sprintf("id=%d", id))
		}
		return
	}
	defer db.Close()
	return user, nil
}

func dbConn() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "gogogo"
	dbPass := "gogogo123"
	addr := "127.0.0.1:13306"
	dbName := "gogogo_db"
	return sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+addr+")/"+dbName)
}
