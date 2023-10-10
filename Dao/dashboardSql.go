package Dao

import (
	"encoding/json"
	"errors"
	regexp2 "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"yxd/conf"
)

type email struct {
	Id    string
	Mail  string
	Token string
}

type dEmail struct {
	status bool
	msg    string
}

func SelectAllEmail(s sessions.Session) (emails []email, err error) {
	db, err := conf.SqlConn()
	if err != nil {
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("select id,mail,token from mail where userid = ? and username = ?")
	if err != nil {
		return nil, errors.New("dashboard 查询邮箱错误 011")
	}
	rows, err := stmt.Query(s.Get("id"), s.Get("username"))
	//rows, err := stmt.Query("6","wcwc")
	if err != nil {
		return nil, errors.New("dashboard 查询邮箱错误 012")
	}
	defer stmt.Close()
	defer rows.Close()
	for rows.Next() {
		var id string
		var mail string
		var token string
		err := rows.Scan(&id, &mail, &token)
		if err != nil {
			return nil, errors.New("dashboard 查询邮箱扫描出错 013")
		}
		u := email{
			Id:    id,
			Mail:  mail,
			Token: token,
		}
		emails = append(emails, u)
	}
	return emails, nil
}

func CreateEmail(s sessions.Session, address string, prefix string) (emails email, err error) {

	db, err := conf.SqlConn()
	if err != nil {
		return email{}, errors.New("数据库连接错误 037")
	}
	defer db.Close()
	match, _ := regexp.Match("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", []byte(address))
	if !match {
		return email{}, errors.New("邮箱错误！ 017")
	}
	if !exist_Email(address) {
		return email{}, errors.New("邮箱已存在！ 018")
	}

	stmt, err := db.Prepare("insert into mail(userid,username,mail,token) value(?,?,?,?)")
	if err != nil {
		return email{}, errors.New("邮箱创建错误！ 012")
	}
	var t string
	for {
		t = conf.RandomString(16)
		reg2, err2 := regexp2.Compile("^(?=.*[0-9].*)(?=.*[A-Z].*)(?=.*[a-z].*).{9,}$", 0)
		if err2 != nil {
			return email{}, errors.New("正则表达式创建失败 038")
		}
		matchString, err2 := reg2.MatchString(t)
		if err != nil {
			return email{}, errors.New("正则表达式判断失败 039")
		}
		if matchString {
			break
		}
	}

	if err = conf.AddEmail(address, t, prefix); err != nil {
		return email{}, errors.New("邮箱创建错误！ 016")
	}
	_, err = stmt.Exec(s.Get("id"), s.Get("username"), address, t)
	if err != nil {
		return email{}, errors.New("邮箱创建错误！ 013")
	}
	return email{
		Id:    s.Get("id").(string),
		Mail:  s.Get("username").(string),
		Token: t,
	}, nil
}

func DeleteEmail(s sessions.Session, email string) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	p := conf.GetKeyData()
	p.Add("username", email)
	post, err := http.Post(conf.GetUrl()+"/plugin?action=a&name=mail_sys&s=delete_mailbox", "application/x-www-form-urlencoded", strings.NewReader(p.Encode()))
	if err != nil {
		return false
	}
	defer post.Body.Close()
	all, err := ioutil.ReadAll(post.Body)
	if err != nil {
		return false
	}
	var d dEmail
	json.Unmarshal(all, &d)
	if d.status {
		return false
	}
	stmt, err := db.Prepare("delete from mail where userid = ? and mail = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	re, err := stmt.Exec(s.Get("id"), email)
	if err != nil {
		return false
	}
	num, _ := re.RowsAffected()
	if num == 1 {
		return true
	}
	return false
}
func exist_Email(mail string) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("select count(0) from mail where mail = ? limit 1")
	if err != nil {
		return false
	}
	defer stmt.Close()
	var re struct {
		num int
	}
	err = stmt.QueryRow(mail).Scan(&re.num)
	if err != nil {
		return false
	}
	if re.num > 0 {
		return false
	}

	return true

}

func Is_Self(id string, mail string) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("select count(0) from mail where userid = ? and mail = ?")
	if err != nil {
		return false
	}
	if err != nil {
		return false
	}
	var re struct {
		num int
	}
	defer stmt.Close()
	err = stmt.QueryRow(id, mail).Scan(&re.num)
	if err != nil {
		return false
	}
	if re.num > 0 {
		return true
	}
	return false
}

//判断数据库邮箱与密码匹配是否存在
func Is_Email(email string, token string) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("select count(0) from mail where mail = ? and token = ?")
	if err != nil {
		return false
	}
	if err != nil {
		return false
	}
	var re struct {
		num int
	}
	defer stmt.Close()
	err = stmt.QueryRow(email, token).Scan(&re.num)
	if err != nil {
		return false
	}
	if re.num > 0 {
		return true
	}
	return false
}

//判断用户是否存在
func Is_Pid(pid string) (bool, string) {
	db, err := conf.SqlConn()
	if err != nil {
		return false, ""
	}
	defer db.Close()
	stmt, err := db.Prepare("select username from user where id = ?")
	if err != nil {
		return false, ""
	}
	var re struct {
		username string
	}
	defer stmt.Close()
	err = stmt.QueryRow(pid).Scan(&re.username)
	if err != nil {
		return false, ""
	}
	//fmt.Println(re.username)
	if re.username != "" {
		return true, re.username
	}
	return false, ""
}

//
// Transfer
//  @Description: 修改邮箱所有人
//  @param id 要修改的所有人id
//  @param mail 邮箱
//  @return bool
//
func Transfer(id string, mail string, username string) bool {
	db, err := conf.SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()
	stmt, err := db.Prepare("update mail set userid = ?, username = ? where mail = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, username, mail)
	if err != nil {
		return false
	}
	return true
}
