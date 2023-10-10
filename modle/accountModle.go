package modle

import (
	"github.com/gin-contrib/sessions"
	"yxd/conf"
)

type a struct {
	num int
}
//
// IsOk
//  @ 判断用户信息是否实时正确
//  @param s
//  @return bool
//
func IsOk(s sessions.Session) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT COUNT(*) from user where id = ? and username = ? and token = ? LIMIT 1")
	if err != nil {
		return false
	}
	var b a
	stmt.QueryRow(s.Get("id"),s.Get("username"),s.Get("token")).Scan(&b.num)
	if b.num == 1 {
		return true
	}
	return false
}
