/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */
package nreq

import (
	"errors"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
	"sync"
	"time"
)

var DefaultTimeout int64 = 30
var mNR sync.Mutex
var mapInstanceNR = map[string]*NRequest{}

type NRequest struct {
	Name    string
	Url     string
	Auth    string
	Opts    []nats.Option
	Timeout time.Duration
	Conn    *nats.Conn
}

func NewNRequest(name string) *NRequest {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NRequest_" + nutil.GetGUUID()
	rqopts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	rqurl := c.GetString(name + ".req.url")
	log.Printf("rqurl=%s", rqurl)
	rqauth := c.GetString(name + ".req.auth")
	//log.Printf("rqauth: %s", rqauth)
	if len(rqauth) > 0 {
		arrauth := strings.Split(rqauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			rqopts = append(rqopts, nats.UserInfo(username, password))
		}
	}
	timeout := c.GetInt64(name + ".req.timeout")
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	rqtimeout := time.Duration(timeout) * time.Second
	nc, err := nats.Connect(rqurl, rqopts...)
	if err != nil {
		log.Println(err)
	}
	return &NRequest{Name: id, Url: rqurl, Auth: rqauth, Opts: rqopts, Timeout: rqtimeout, Conn: nc}
}

func GetInstance(name string) *NRequest {
	instance := mapInstanceNR[name]
	if instance == nil || instance.Conn.IsClosed() {
		mNR.Lock()
		defer mNR.Unlock()
		instance = mapInstanceNR[name]
		if instance == nil || instance.Conn.IsClosed() {
			instance = NewNRequest(name)
			mapInstanceNR[name] = instance
		}
	}
	return instance
}

func (nr *NRequest) Request(subject string, data string) (*nats.Msg, error) {
	if len(subject) == 0 || len(data) == 0 {
		return nil, errors.New("Input params invalid")
	}
	// Request to NATS
	msg, err := nr.Conn.Request(subject, []byte(data), nr.Timeout)
	nr.Conn.Flush()
	if err := nr.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NReq Published [%s] : '%s'", subject, data)
	//log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	return msg, err
}

func (nr *NRequest) RequestByte(subject string, data []byte) (*nats.Msg, error) {
	if len(subject) == 0 || len(data) == 0 {
		return nil, errors.New("Input params invalid")
	}
	// Request to NATS
	msg, err := nr.Conn.Request(subject, data, nr.Timeout)
	nr.Conn.Flush()
	if err := nr.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NReq Published [%s] : '%s'", subject, data)
	//log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	return msg, err
}

func Request(name string, subject string, data string) (*nats.Msg, error) {
	if len(name) == 0 || len(subject) == 0 || len(data) == 0 {
		return nil, errors.New("Input params invalid")
	}
	nr := GetInstance(name)
	// Request to NATS
	msg, err := nr.Conn.Request(subject, []byte(data), nr.Timeout)
	nr.Conn.Flush()
	if err := nr.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NReq Published [%s] : '%s'", subject, data)
	//log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	return msg, err
}

func RequestByte(name string, subject string, data []byte) (*nats.Msg, error) {
	if len(name) == 0 || len(subject) == 0 || len(data) == 0 {
		return nil, errors.New("Input params invalid")
	}
	nr := GetInstance(name)
	// Request to NATS
	msg, err := nr.Conn.Request(subject, data, nr.Timeout)
	nr.Conn.Flush()
	if err := nr.Conn.LastError(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("NReq Published [%s] : '%s'", subject, data)
	//log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
	return msg, err
}

//var rqurl string
//var rqauth string
//var rqopts []nats.Option
//var rqtimeout time.Duration
//
//func InitReqConf(name string) error {
//	if name == "" {
//		return errors.New("name config is not empty.")
//	}
//	rqopts = []nats.Option{nats.Name("NReq_" + nutil.GetGUUID())}
//	c := nconf.GetConfig()
//	rqurl = c.GetString(name+".req.url")
//	log.Printf("rqurl=%s", rqurl)
//	rqauth = c.GetString(name+".req.auth")
//	//log.Printf("rqauth: %s", rqauth)
//	if len(rqauth) > 0 {
//		arrauth := strings.Split(rqauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			rqopts = append(rqopts, nats.UserInfo(username, password))
//		}
//	}
//	timeout := c.GetInt64(name+".req.timeout")
//	if timeout == 0 {
//		timeout = DefaultTimeout
//	}
//	rqtimeout = time.Duration(timeout)*time.Second
//	return nil
//}
//
//// url = nats://127.0.0.1:4222
//// auth = username:password
//// timeout = 10*time.Second
//func InitReqParams(url string, auth string, timeout time.Duration) error {
//	rqopts = []nats.Option{nats.Name("NReq_" + nutil.GetGUUID())}
//	rqurl = url
//	rqauth = auth
//	if len(rqauth) > 0 {
//		arrauth := strings.Split(rqauth, ":")
//		if len(arrauth) == 2 {
//			username := arrauth[0]
//			password := arrauth[1]
//			rqopts = append(rqopts, nats.UserInfo(username, password))
//		}
//	}
//	if timeout <= 0 {
//		rqtimeout = time.Duration(DefaultTimeout)*time.Second
//	} else {
//		rqtimeout = timeout
//	}
//	return nil
//}
//
//func GetUrl() string {
//	return rqurl
//}
//
//func GetAuth() string {
//	return rqauth
//}
//
//func GetOption() []nats.Option {
//	return rqopts
//}
//
//func GetTimeout() time.Duration {
//	return rqtimeout
//}
//
//func GetConnect() (*nats.Conn, error) {
//	// Connect to NATS
//	nc, err := nats.Connect(rqurl, rqopts...)
//	if err != nil {
//		log.Println(err)
//	}
//	return nc, err
//}

//func Request(subject string, data string) (*nats.Msg, error) {
//	// Connect to NATS
//	nc, err := nats.Connect(rqurl, rqopts...)
//	defer nc.Close()
//	if err != nil {
//		log.Println(err)
//	}
//	msg, err := nc.Request(subject, []byte(data), rqtimeout)
//	nc.Flush()
//	if err := nc.LastError(); err != nil {
//		log.Fatal(err)
//	}
//	//log.Printf("NReq Published [%s] : '%s'", subject, data)
//	//log.Printf("NReq Received  [%v] : '%s'", msg.Subject, string(msg.Data))
//	return msg, err
//}
