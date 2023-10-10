package conf

var regkey = "asdbvc"        // 密码加盐，可以随便设置
var serv = "xxx.abc.com:143" // 邮箱服务器地址

func GetRegKey() string {
	return regkey
}

func GetServ() string {
	return serv
}
