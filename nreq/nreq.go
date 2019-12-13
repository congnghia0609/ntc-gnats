package nreq

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"ntc-gnats/nutil"
	"strings"
	"time"
)

var rqurl string
var rqauth string
var rqopts []nats.Option
var DefaultTimeout int64 = 30
var rqtimeout time.Duration

func InitReqConf(name string) error {
	if name == "" {
		return errors.New("name config is not empty.")
	}
	rqopts = []nats.Option{nats.Name("NReq_" + nutil.GetGUUID())}
	c := nconf.GetConfig()
	rqurl = c.GetString(name+".req.url")
	//log.Printf("rqurl: %s", rqurl)
	rqauth = c.GetString(name+".req.auth")
	//log.Printf("rqauth: %s", rqauth)
	if len(rqauth) > 0 {
		arrauth := strings.Split(rqauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rqopts = append(rqopts, nats.UserInfo(username, password))
		}
	}
	timeout := c.GetInt64(name+".req.timeout")
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	rqtimeout = time.Duration(timeout)*time.Second
	return nil
}

// url = nats://127.0.0.1:4222
// auth = username:password
// timeout = 10*time.Second
func InitReqParams(url string, auth string, timeout time.Duration) error {
	rqopts = []nats.Option{nats.Name("NReq_" + nutil.GetGUUID())}
	rqurl = url
	rqauth = auth
	if len(rqauth) > 0 {
		arrauth := strings.Split(rqauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rqopts = append(rqopts, nats.UserInfo(username, password))
		}
	}
	if timeout <= 0 {
		rqtimeout = time.Duration(DefaultTimeout)*time.Second
	} else {
		rqtimeout = timeout
	}
	return nil
}

func GetUrl() string {
	return rqurl
}

func GetAuth() string {
	return rqauth
}

func GetOption() []nats.Option {
	return rqopts
}

func GetTimeout() time.Duration {
	return rqtimeout
}

func GetConnect() (*nats.Conn, error) {
	// Connect to NATS
	nc, err := nats.Connect(rqurl, rqopts...)
	if err != nil {
		log.Println(err)
	}
	return nc, err
}

func Request(subject string, data string) error {
	// Connect to NATS
	nc, err := nats.Connect(rqurl, rqopts...)
	defer nc.Close()
	if err != nil {
		log.Println(err)
	}
	nc.Request(subject, []byte(data), rqtimeout)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	return err
}
