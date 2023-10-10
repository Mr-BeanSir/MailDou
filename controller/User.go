package controller

import (
	"fmt"
	"yxd/conf"
)

type UserInfo struct {
	Id       string `json:"id" form:"e"`
	Username string `json:"username" form:"b"`
	Token    string `json:"token" form:"c"`
	Time     string `json:"time" form:"a"`
	Regtoken string `json:"regtoken" form:"d"`
}

// Login
//  @ 登录失败返回type=1，error为错误提示
//	@ 登录成功返回type=0，user为返回的用户信息
//  @receiver u
//  @return map[string]interface{}
//
func (u *UserInfo) Login() map[string]interface{} {
	re := make(map[string]interface{})
	db, err := conf.SqlConn()
	if err != nil {
		re["type"] = false
		re["error"] = "错误：数据库连接失败"
		re["tag"] = "000"
		return re
	}
	defer db.Close()
	is, user := SelectPrepare(u.Token, u.Id, u.Username)
	if !is {
		re["type"] = false
		re["error"] = "登录失败：错误的token"
		re["tag"] = "005"
		return re
	}
	re["type"] = true
	re["error"] = "登录成功：稍等跳转"
	re["user"] = user
	return re
}

//
// Reg
//  @ type为false则注册失败，error为错误信息
//  @ type为true注册成功，user返回users结构体
//  @receiver userinfo结构体
//  @return map[string]interface{}
//
func (u *UserInfo) Reg() map[string]interface{} {
	re := make(map[string]interface{})
	fmt.Println(u.Id)
	db, err := conf.SqlConn()
	if err != nil {
		re["type"] = false
		re["error"] = "错误：数据库连接失败"
		re["tag"] = "000"
		return re
	}
	defer db.Close()
	sqlStr := "insert into user value(?,?,?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		re["type"] = false
		re["error"] = "注册失败：操作数据库错误，可能存在相同的用户名或ID"
		re["tag"] = "001"
		return re
	}
	defer stmt.Close()
	entry := conf.Md5V(conf.GetRegKey() + u.Token + u.Time)
	if entry == u.Regtoken {
		_, err := stmt.Exec(u.Id, u.Username, u.Token)
		if err != nil {
			re["type"] = false
			re["error"] = "注册失败：操作数据库错误"
			re["tag"] = "002"
			return re
		}
		var user UserInfo
		user.Id = u.Id
		user.Token = u.Token
		user.Username = u.Username
		re["type"] = true
		re["error"] = "注册成功"
		re["user"] = user
		return re
	} else {
		re["type"] = false
		re["error"] = "注册失败：错误的token"
		re["tag"] = "003"
		return re
	}
}

// SelectPrepare
//  @param db 数据库连接
//  @param token
//  @param id
//  @param username
//  @return bool false为查询失败
//  @return Users 结构体
//
func SelectPrepare(token string, id string, username string) (bool, UserInfo) {
	db, err := conf.SqlConn()
	if err != nil {
		return false, UserInfo{}
	}
	defer db.Close()
	sqlStr := "select * from user where token = ? and id = ? and username = ?"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return false, UserInfo{}
	}
	defer stmt.Close()
	var user UserInfo
	err = stmt.QueryRow(token, id, username).Scan(&user.Id, &user.Username, &user.Token)
	if err != nil {
		return false, UserInfo{}
	}
	return true, user

}
