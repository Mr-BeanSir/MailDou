package controller

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"yxd/Dao"
	"yxd/modle"
)

type EnvelopeInfo struct {
	Uid           string
	Title         string
	SenderName    string
	SenderAddress string
	Time          string
}

func IndexMail(c *gin.Context) {
	session := sessions.Default(c)
	usernameMail, b := c.GetPostForm("usernameMail")
	if !b {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱用户名为空",
			"tag":   "021",
		})
		return
	}
	passwordMail, b := c.GetPostForm("passwordMail")
	if !b {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱密码为空",
			"tag":   "022",
		})
		return
	}
	conn, err := modle.ConnectMail()
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "链接邮箱服务器失败，我他妈是本地怎么失败了我草",
			"tag":   "023",
		})
		return
	}
	defer conn.Close()
	if !Dao.Is_Self(session.Get("id").(string), usernameMail) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	err = modle.LoginMail(conn, usernameMail, passwordMail)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "登录邮箱失败，你tm别乱搞我操",
			"tag":   "024",
		})
		return
	}
	defer conn.Logout()
	session.Set("usernameMail", usernameMail)
	session.Set("passwordMail", passwordMail)
	session.Set("Mail", "1")
	session.Save()
	log.Println(session.Get("usernameMail"))
	log.Println(session.Get("passwordMail"))
	c.Header("location", "/mail/list/INBOX")
	c.AbortWithStatus(http.StatusFound)

}

func ListMail(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("Mail") != "1" {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "你还没tm登录邮箱呢，看个鸡皮",
			"tag":   "025",
		})
		return
	}
	usernameMail := session.Get("usernameMail").(string)
	passwordMail := session.Get("passwordMail").(string)
	conn, err := modle.ConnectMail()
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "链接邮箱服务器失败，我他妈是本地怎么失败了我草",
			"tag":   "023",
		})
		return
	}
	defer conn.Close()
	if !Dao.Is_Self(session.Get("id").(string), usernameMail) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	err = modle.LoginMail(conn, usernameMail, passwordMail)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "登录邮箱失败，你tm别乱搞我操",
			"tag":   "024",
		})
		return
	}
	defer conn.Logout()
	mbox, err := conn.Select(c.Param("type"), false)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "选择邮箱失败，你没有这个类型的邮箱",
			"tag":   "028",
		})
		return
	}
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 20 {
		from = mbox.Messages - 20
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 30)
	done := make(chan error, 1)
	imap.CharsetReader = charset.Reader
	go func() {
		done <- conn.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	//查询邮件报错直接寄
	bo := make(chan bool)
	go func(c *gin.Context) {
		if err := <-done; err != nil {
			if len(messages) != 0 {
				log.Println(err)
				c.HTML(http.StatusOK, "index/error.html", gin.H{
					"error": err,
					"tag":   "300",
				})
				go func() {
					for _ = range messages {
					}
				}()
				bo <- true
				return
			}
			fmt.Println(err)
		}
		bo <- false
	}(c)
	if <-bo {
		return
	}

	mailRe := make([]EnvelopeInfo, 0)
	for msg := range messages {
		mailRe = append(mailRe, EnvelopeInfo{
			Uid:           strconv.Itoa(int(msg.SeqNum)),
			Title:         msg.Envelope.Subject,
			SenderName:    msg.Envelope.From[0].PersonalName,
			SenderAddress: msg.Envelope.From[0].MailboxName + msg.Envelope.From[0].HostName,
			Time:          msg.Envelope.Date.String(),
		})
	}
	Reverse(&mailRe, len(mailRe))
	c.HTML(http.StatusOK, "mail/list.html", gin.H{
		"mailRe":   mailRe,
		"username": session.Get("username"),
		"id":       session.Get("id"),
		"type":     c.Param("type"),
	})
}

func ReadMail(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("Mail") != "1" {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "你还没tm登录邮箱呢，看个鸡皮",
			"tag":   "025",
		})
		return
	}
	usernameMail := session.Get("usernameMail").(string)
	passwordMail := session.Get("passwordMail").(string)
	conn, err := modle.ConnectMail()
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "链接邮箱服务器失败，我他妈是本地怎么失败了我草",
			"tag":   "026",
		})
		return
	}
	defer conn.Close()
	if !Dao.Is_Self(session.Get("id").(string), usernameMail) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	err = modle.LoginMail(conn, usernameMail, passwordMail)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "登录邮箱失败，你tm别乱搞我操",
			"tag":   "027",
		})
		return
	}
	defer conn.Logout()
	_, err = conn.Select(c.Param("type"), false)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "选择邮箱失败，你没有这个类型的邮箱",
			"tag":   "028",
		})
		return
	}
	seqset := new(imap.SeqSet)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "id传递错误，你tm别乱搞我操",
			"tag":   "030",
		})
		return
	}
	seqset.AddNum(uint32(id))

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	//imap.FormatFlagsOp()
	//seenItem := imap.FormatFlagsOp(imap.SetFlags, true)
	//flags := []interface{}{imap.SeenFlag}
	//err = conn.UidStore(seqset, seenItem, flags, nil)
	//if err != nil {
	//	log.Println(err)
	//}

	go func() {
		done <- conn.Fetch(seqset, items, messages)
	}()

	imap.CharsetReader = charset.Reader
	msg := <-messages

	log.Println("Message has been marked as seen")
	mailinfo := msg.GetBody(section)

	log.Println(msg.Flags)

	if mailinfo == nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "空的邮件，联系站长查询",
			"tag":   "032",
		})
		return
	}

	m, err := mail.CreateReader(mailinfo)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "读取邮件错误！！！",
			"tag":   "033",
		})
		return
	}

	header := m.Header
	//log.Println(header.Date())
	date, _ := header.Date()
	//date, err := header.Date()
	//if err != nil {
	//	c.HTML(http.StatusOK, "index/error.html", gin.H{
	//		"error": "邮件时间存在错误哈哈哈，小问题防止出错不加载",
	//		"tag":   "034",
	//	})
	//	return
	//}
	defer m.Close()
	//var context string
	//var text string
	//for {
	//	p, err := m.NextPart()
	//	if err == io.EOF {
	//		break
	//	} else if err != nil {
	//		c.HTML(http.StatusOK, "index/error.html", gin.H{
	//			"error": "获取邮件正文出错",
	//			"tag":   "036",
	//		})
	//		return
	//	}
	//
	//	//var context string
	//	//var text string
	//	switch p.Header.(type) {
	//	case *mail.InlineHeader:
	//		// 正文消息文本
	//		b, _ := ioutil.ReadAll(p.Body)
	//		contentType := p.Header.Get("Content-Type")
	//		fmt.Println(contentType[:strings.Index(contentType, ";")])
	//		if contentType[:strings.Index(contentType, ";")] == "text/html" {
	//			context = string(b)
	//			context = delHMLTag(context)
	//		} else {
	//			text = string(b)
	//
	//		}
	//	}
	//
	//}
	//log.Println(header.Get("flag"))
	specificDate := date.Format("2006年01月2日")
	specificTime := date.Format("15:04:05")

	c.HTML(http.StatusOK, "mail/read.html", gin.H{
		"subject":      msg.Envelope.Subject,
		"from":         msg.Envelope.From[0].MailboxName + "@" + msg.Envelope.From[0].HostName,
		"to":           msg.Envelope.To[0].MailboxName + "@" + msg.Envelope.To[0].HostName,
		"date":         date.String(),
		"flag":         header.Get("flag"),
		"specificDate": specificDate,
		"specificTime": specificTime,
		"senderName":   msg.Envelope.From[0].PersonalName,
		"username":     session.Get("username"),
		"id":           session.Get("id"),
		"type":         c.Param("type"),
		"mailId":       c.Param("id"),
	})
}

func Reverse(arr *[]EnvelopeInfo, length int) {
	var temp EnvelopeInfo
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

func delHMLTag(html string) string {
	regEx_script := "<script[^>]*?>[\\s\\S]*?<\\/script>"
	regEx_style := "<style[^>]*?>[\\s\\S]*?<\\/style>"
	//regEx_html := "<[^>]+>"
	//regEx_space := "\\s*|\t|\r|\n"

	compile, _ := regexp2.Compile(regEx_style, 0)
	html, _ = compile.Replace(html, "", 0, 1)
	compile, _ = regexp2.Compile(regEx_script, 0)
	html, _ = compile.Replace(html, "", 0, 1)

	return html
}

func DelMail(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("Mail") != "1" {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "你还没tm登录邮箱呢，看个鸡皮",
			"tag":   "025",
		})
		return
	}
	usernameMail := session.Get("usernameMail").(string)
	passwordMail := session.Get("passwordMail").(string)
	conn, err := modle.ConnectMail()
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "链接邮箱服务器失败，我他妈是本地怎么失败了我草",
			"tag":   "026",
		})
		return
	}
	defer conn.Close()
	if !Dao.Is_Self(session.Get("id").(string), usernameMail) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	err = modle.LoginMail(conn, usernameMail, passwordMail)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "登录邮箱失败，你tm别乱搞我操",
			"tag":   "027",
		})
		return
	}

	_, err = conn.Select(c.Param("type"), false)
	if err != nil {
		c.Header("location", "/mail/list/INBOX")
		c.AbortWithStatus(302)
		return
	}

	seqSet := new(imap.SeqSet)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	seqSet.AddNum(uint32(id))
	item := imap.FormatFlagsOp(imap.AddFlags, false)
	flags := []interface{}{imap.DeletedFlag}
	err = conn.Store(seqSet, item, flags, nil)
	conn.Expunge(nil)
	if err != nil {
		c.Header("location", "/mail/list/INBOX")
		c.AbortWithStatus(302)
		return
	}
}

func Html(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("Mail") != "1" {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "你还没tm登录邮箱呢，看个鸡皮",
			"tag":   "025",
		})
		return
	}
	usernameMail := session.Get("usernameMail").(string)
	passwordMail := session.Get("passwordMail").(string)
	conn, err := modle.ConnectMail()
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "链接邮箱服务器失败，我他妈是本地怎么失败了我草",
			"tag":   "026",
		})
		return
	}
	defer conn.Close()
	if !Dao.Is_Self(session.Get("id").(string), usernameMail) {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "邮箱不属于你，你tm想干坏事",
			"tag":   "110",
		})
		return
	}
	err = modle.LoginMail(conn, usernameMail, passwordMail)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "登录邮箱失败，你tm别乱搞我操",
			"tag":   "027",
		})
		return
	}
	defer conn.Logout()
	_, err = conn.Select(c.Param("type"), false)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "选择邮箱失败，你没有这个类型的邮箱",
			"tag":   "028",
		})
		return
	}
	seqset := new(imap.SeqSet)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "id传递错误，你tm别乱搞我操",
			"tag":   "030",
		})
		return
	}
	seqset.AddNum(uint32(id))

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	//imap.FormatFlagsOp()
	//seenItem := imap.FormatFlagsOp(imap.SetFlags, true)
	//flags := []interface{}{imap.SeenFlag}
	//err = conn.UidStore(seqset, seenItem, flags, nil)
	//if err != nil {
	//	log.Println(err)
	//}

	go func() {
		done <- conn.Fetch(seqset, items, messages)
	}()

	imap.CharsetReader = charset.Reader
	msg := <-messages

	log.Println("Message has been marked as seen")
	mailinfo := msg.GetBody(section)

	log.Println(msg.Flags)

	if mailinfo == nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "空的邮件，联系站长查询",
			"tag":   "032",
		})
		return
	}

	m, err := mail.CreateReader(mailinfo)
	if err != nil {
		c.HTML(http.StatusOK, "index/error.html", gin.H{
			"error": "读取邮件错误！！！",
			"tag":   "033",
		})
		return
	}

	//log.Println(header.Date())
	//date, err := header.Date()
	//if err != nil {
	//	c.HTML(http.StatusOK, "index/error.html", gin.H{
	//		"error": "邮件时间存在错误哈哈哈，小问题防止出错不加载",
	//		"tag":   "034",
	//	})
	//	return
	//}
	defer m.Close()
	var context string
	var text string
	for {
		p, err := m.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			c.HTML(http.StatusOK, "index/error.html", gin.H{
				"error": "获取邮件正文出错",
				"tag":   "036",
			})
			return
		}

		//var context string
		//var text string
		switch p.Header.(type) {
		case *mail.InlineHeader:
			// 正文消息文本
			b, _ := ioutil.ReadAll(p.Body)
			contentType := p.Header.Get("Content-Type")
			fmt.Println(contentType[:strings.Index(contentType, ";")])
			if contentType[:strings.Index(contentType, ";")] == "text/html" {
				context = string(b)
				context = delHMLTag(context)
			} else {
				text = string(b)

			}
		}

	}
	//log.Println(header.Get("flag"))
	c.Header("Content-Type", "text/html; charset=utf-8")
	if context != "" {
		c.String(http.StatusOK, context)
	} else {
		c.String(http.StatusOK, text)
	}
}
