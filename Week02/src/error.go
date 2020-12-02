package main

import (
	"database/sql"
	gerrors "errors"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

var DataNotFound = gerrors.New("record not found")

type User struct {
	Id   int64
	Name string
}

func main() {
	user, err := biz(1)
	if err != nil {
		log.Fatalf("find user failed: %+v", err)
	}
	fmt.Println(user)
}

func biz(id int64) (user *User, err error) {
	return FindById(id)
}

func FindById(id int64) (user *User, err error) {
	db, err := dbConn()
	stmt, err := db.Prepare("select * from go_user where id=?")
	if err != nil {
		return user, errors.Wrap(err, "查询用户失败！")
	}
	err = stmt.QueryRow(id).Scan(&user)

	if err != nil && gerrors.Is(err, sql.ErrNoRows) {
		err = errors.Wrap(DataNotFound, "query failed")
		return
	}
	defer db.Close()
	return user, nil
}

func dbConn() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "goblog"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return db, errors.Wrap(err, "连续数据库失败！")
	}
	return db, nil
}
