/**
 *
 * @author nghiatc
 * @since Dec 6, 2019
 */

package npub

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/nats-io/nats.go"
	"log"
	"math"
	"ntc-gnats/nutil"
	"strings"
	"sync"
)

var mNP sync.Mutex
var mapInstanceNP = map[string]*NPublisher{}

type NPublisher struct {
	name string
	url  string
	auth string
	opts []nats.Option
	conn *nats.Conn
}

func (np *NPublisher) GetName() string {
	return np.name
}

func (np *NPublisher) GetUrl() string {
	return np.url
}

func (np *NPublisher) GetAuth() string {
	return np.auth
}

func (np *NPublisher) GetOption() []nats.Option {
	return np.opts
}

func (np *NPublisher) GetConnect() *nats.Conn {
	return np.conn
}

func NewNPublisher(name string) *NPublisher {
	if len(name) == 0 {
		fmt.Errorf("name config is not empty.")
		return nil
	}
	c := nconf.GetConfig()
	id := name + "_NPublisher_" + nutil.GetGUUID()
	popts := []nats.Option{nats.Name(id), nats.MaxReconnects(math.MaxInt32)}
	purl := c.GetString(name + ".pub.url")
	fmt.Printf("purl=%s\n", purl)
	pauth := c.GetString(name + ".pub.auth")
	//log.Printf("pauth: %s\n", pauth)
	if len(pauth) > 0 {
		arrauth := strings.Split(pauth, ":")
		if len(arrauth) == 2 {
			username := arrauth[0]
			password := arrauth[1]
			popts = append(popts, nats.UserInfo(username, password))
		}
	}
	// Connect to NATS
	nc, err := nats.Connect(purl, popts...)
	if err != nil {
		fmt.Printf("Connect to NATS Fail: %v\n", err)
	}
	return &NPublisher{name: id, url: purl, auth: pauth, opts: popts, conn: nc}
}

func GetInstance(name string) *NPublisher {
	instance := mapInstanceNP[name]
	if instance == nil || instance.conn.IsClosed() {
		mNP.Lock()
		defer mNP.Unlock()
		instance = mapInstanceNP[name]
		if instance == nil || instance.conn.IsClosed() {
			instance = NewNPublisher(name)
			mapInstanceNP[name] = instance
		}
	}
	return instance
}

func (np *NPublisher) Publish(subject string, msg string) error {
	if len(subject) > 0 && len(msg) > 0 {
		np.conn.Publish(subject, []byte(msg))
		np.conn.Flush()
		err := np.conn.LastError()
		if err != nil {
			log.Fatalf("NPublisher publish message Error: %v\n", err)
		}
		return err
	}
	return nil
}

func (np *NPublisher) PublishByte(subject string, msg []byte) error {
	if len(subject) > 0 && len(msg) > 0 {
		np.conn.Publish(subject, msg)
		np.conn.Flush()
		err := np.conn.LastError()
		if err != nil {
			log.Fatalf("NPublisher publish message Error: %v\n", err)
		}
		return err
	}
	return nil
}

func (np *NPublisher) Close() {
	if np != nil {
		np.conn.Close()
	}
}

func Publish(name string, subject string, msg string) error {
	np := GetInstance(name)
	return np.Publish(subject, msg)
}

func PublishByte(name string, subject string, msg []byte) error {
	np := GetInstance(name)
	return np.PublishByte(subject, msg)
}
