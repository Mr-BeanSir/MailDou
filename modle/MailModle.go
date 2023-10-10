package modle

import (
	"github.com/emersion/go-imap/client"
	"yxd/conf"
)

func ConnectMail() (*client.Client, error) {
	dial, err := client.Dial(conf.GetServ())
	if err != nil {
		return nil, err
	}
	return dial, err

}

func LoginMail(c *client.Client, username string, password string) error {
	err := c.Login(username, password)
	if err != nil {
		return err
	}
	return nil
}
