package conf

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var token = "宝塔TOKEN"
var urls = "宝塔API地址"

type rData struct {
	Status bool
	Msg    string `json:"msg"`
}

func DeleteEmail() {

}

func ChangePassword(address string) bool {
	db, err := SqlConn()
	if err != nil {
		return false
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE mail set token=? WHERE mail = ?")
	if err != nil {
		return false
	}
	pwd := RandomString(16)
	p := GetKeyData()
	p.Add("quota", "200 M")
	p.Add("username", address)
	p.Add("password", pwd)
	p.Add("full_name", address)
	p.Add("is_admin", "0")
	p.Add("active", "1")
	req, err := http.NewRequest(http.MethodPost, urls+"/plugin?action=a&name=mail_sys&s=update_mailbox", strings.NewReader(p.Encode()))
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 15,
	}
	do, err := client.Do(req)
	defer do.Body.Close()
	all, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return false
	}
	var r rData
	err = json.Unmarshal(all, &r)
	if err != nil {
		return false
	}
	if !r.Status {
		log.Println(r.Msg)
		return false
	}
	rows, err := stmt.Exec(pwd, address)
	if err != nil {
		return false
	}
	num, err := rows.RowsAffected()
	if err != nil {
		return false
	}
	if num != 1 {
		return false
	}
	return true
}

func AddEmail(address string, password string, prefix string) error {
	p := GetKeyData()
	p.Add("quota", "200 M")
	p.Add("username", address)
	p.Add("password", password)
	p.Add("full_name", prefix)
	p.Add("is_admin", "0")
	//log.Println(p.Encode())
	req, err := http.NewRequest(http.MethodPost, urls+"/plugin?action=a&name=mail_sys&s=add_mailbox", strings.NewReader(p.Encode()))

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 15,
	}
	do, err := client.Do(req)
	defer do.Body.Close()
	all, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return err
	}
	var r rData
	err = json.Unmarshal(all, &r)
	if err != nil {
		return err
	}
	log.Println(string(all))
	log.Println(r.Msg)
	if !r.Status {
		log.Println(r.Msg)
		return errors.New("错误，服务器创建邮箱失败")
	}
	return nil
}

func GetKeyData() url.Values {
	timeNow := time.Now().Unix()
	u := url.Values{}
	u.Add("request_time", strconv.Itoa(int(timeNow)))
	u.Add("request_token", Md5V(strconv.Itoa(int(timeNow))+Md5V(token)))
	return u
}

func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUrl() string {
	return urls
}
