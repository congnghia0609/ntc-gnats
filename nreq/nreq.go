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
	return msg, err
}
